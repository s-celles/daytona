// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"context"
	"errors"

	"github.com/daytonaio/daytona/internal/util"
	"github.com/daytonaio/daytona/pkg/build"
	"github.com/daytonaio/daytona/pkg/db"
	"github.com/daytonaio/daytona/pkg/gitprovider"
	"github.com/daytonaio/daytona/pkg/logs"
	"github.com/daytonaio/daytona/pkg/models"
	"github.com/daytonaio/daytona/pkg/provider/manager"
	"github.com/daytonaio/daytona/pkg/provisioner"
	"github.com/daytonaio/daytona/pkg/runners"
	"github.com/daytonaio/daytona/pkg/server"
	"github.com/daytonaio/daytona/pkg/server/apikeys"
	"github.com/daytonaio/daytona/pkg/server/builds"
	"github.com/daytonaio/daytona/pkg/server/containerregistries"
	"github.com/daytonaio/daytona/pkg/server/gitproviders"
	"github.com/daytonaio/daytona/pkg/server/jobs"
	"github.com/daytonaio/daytona/pkg/server/targetconfigs"
	"github.com/daytonaio/daytona/pkg/server/targets"
	"github.com/daytonaio/daytona/pkg/server/workspaces"
	"github.com/daytonaio/daytona/pkg/stores"
	"github.com/daytonaio/daytona/pkg/telemetry"

	"github.com/daytonaio/daytona/pkg/runners/runner"
	target_runner "github.com/daytonaio/daytona/pkg/runners/target"
	workspace_runner "github.com/daytonaio/daytona/pkg/runners/workspace"
)

func GetJobRunner(c *server.Config, configDir string, version string, telemetryService telemetry.TelemetryService) (runners.IJobRunner, error) {
	dbPath, err := getDbPath()
	if err != nil {
		return nil, err
	}

	dbConnection := db.GetSQLiteConnection(dbPath)

	jobStore, err := db.NewJobStore(dbConnection)
	if err != nil {
		return nil, err
	}

	jobService := jobs.NewJobService(jobs.JobServiceConfig{
		JobStore: jobStore,
	})

	workspaceJobRunner, err := GetWorkspaceJobRunner(c, configDir, version, telemetryService)
	if err != nil {
		return nil, err
	}

	targetJobRunner, err := GetTargetJobRunner(c, configDir, version, telemetryService)
	if err != nil {
		return nil, err
	}

	return runner.NewJobRunner(runner.JobRunnerConfig{
		ListPendingJobs: func(ctx context.Context) ([]*models.Job, error) {
			return jobService.List(&stores.JobFilter{
				States: &[]models.JobState{models.JobStatePending},
			})
		},
		UpdateJobState: func(ctx context.Context, job *models.Job, state models.JobState, err *error) error {
			job.State = state
			if err != nil {
				job.Error = util.Pointer((*err).Error())
			}
			return jobService.Save(job)
		},
		WorkspaceJobRunner: workspaceJobRunner,
		TargetJobRunner:    targetJobRunner,
		LoggerFactory:      logs.NewLoggerFactory(nil, nil),
		Provisioner:        provisioner.NewProvisioner(provisioner.ProvisionerConfig{}),
	}), nil
}

func GetWorkspaceJobRunner(c *server.Config, configDir string, version string, telemetryService telemetry.TelemetryService) (runners.IWorkspaceJobRunner, error) {
	dbPath, err := getDbPath()
	if err != nil {
		return nil, err
	}

	dbConnection := db.GetSQLiteConnection(dbPath)

	jobStore, err := db.NewJobStore(dbConnection)
	if err != nil {
		return nil, err
	}

	buildStore, err := db.NewBuildStore(dbConnection)
	if err != nil {
		return nil, err
	}
	workspaceConfigStore, err := db.NewWorkspaceConfigStore(dbConnection)
	if err != nil {
		return nil, err
	}

	jobService := jobs.NewJobService(jobs.JobServiceConfig{
		JobStore: jobStore,
	})

	targetStore, err := db.NewTargetStore(dbConnection)
	if err != nil {
		return nil, err
	}

	targetMetadataStore, err := db.NewTargetMetadataStore(dbConnection)
	if err != nil {
		return nil, err
	}

	containerRegistryStore, err := db.NewContainerRegistryStore(dbConnection)
	if err != nil {
		return nil, err
	}

	containerRegistryService := containerregistries.NewContainerRegistryService(containerregistries.ContainerRegistryServiceConfig{
		Store: containerRegistryStore,
	})

	gitProviderConfigStore, err := db.NewGitProviderConfigStore(dbConnection)
	if err != nil {
		return nil, err
	}
	gitProviderService := gitproviders.NewGitProviderService(gitproviders.GitProviderServiceConfig{
		ConfigStore: gitProviderConfigStore,
	})

	targetConfigStore, err := db.NewTargetConfigStore(dbConnection)
	if err != nil {
		return nil, err
	}

	targetConfigService := targetconfigs.NewTargetConfigService(targetconfigs.TargetConfigServiceConfig{
		TargetConfigStore: targetConfigStore,
	})

	apiKeyStore, err := db.NewApiKeyStore(dbConnection)
	if err != nil {
		return nil, err
	}

	apiKeyService := apikeys.NewApiKeyService(apikeys.ApiKeyServiceConfig{
		ApiKeyStore: apiKeyStore,
	})

	targetLogsDir, err := server.GetTargetLogsDir(configDir)
	if err != nil {
		return nil, err
	}
	buildLogsDir, err := build.GetBuildLogsDir()
	if err != nil {
		return nil, err
	}
	loggerFactory := logs.NewLoggerFactory(&targetLogsDir, &buildLogsDir)

	headscaleUrl := util.GetFrpcHeadscaleUrl(c.Frps.Protocol, c.Id, c.Frps.Domain)

	workspaceStore, err := db.NewWorkspaceStore(dbConnection)
	if err != nil {
		return nil, err
	}

	workspaceMetadataStore, err := db.NewWorkspaceMetadataStore(dbConnection)
	if err != nil {
		return nil, err
	}

	providerManager := manager.GetProviderManager(nil)

	provisioner := provisioner.NewProvisioner(provisioner.ProvisionerConfig{
		ProviderManager: providerManager,
	})

	buildService := builds.NewBuildService(builds.BuildServiceConfig{
		BuildStore: buildStore,
		FindWorkspaceConfig: func(ctx context.Context, name string) (*models.WorkspaceConfig, error) {
			return workspaceConfigStore.Find(&stores.WorkspaceConfigFilter{
				Name: &name,
			})
		},
		GetRepositoryContext: func(ctx context.Context, url, branch string) (*gitprovider.GitRepository, error) {
			gitProvider, _, err := gitProviderService.GetGitProviderForUrl(url)
			if err != nil {
				return nil, err
			}

			repo, err := gitProvider.GetRepositoryContext(gitprovider.GetRepositoryContext{
				Url: url,
			})

			return repo, err
		},
		LoggerFactory: loggerFactory,
	})

	targetService := targets.NewTargetService(targets.TargetServiceConfig{
		TargetStore:         targetStore,
		TargetMetadataStore: targetMetadataStore,
		FindTargetConfig: func(ctx context.Context, name string) (*models.TargetConfig, error) {
			return targetConfigService.Find(&stores.TargetConfigFilter{Name: &name})
		},
		GenerateApiKey: func(ctx context.Context, name string) (string, error) {
			return apiKeyService.Generate(models.ApiKeyTypeTarget, name)
		},
		RevokeApiKey: func(ctx context.Context, name string) error {
			return apiKeyService.Revoke(name)
		},
		CreateJob: func(ctx context.Context, targetId string, action models.JobAction) error {
			return jobService.Save(&models.Job{
				ResourceId:   targetId,
				ResourceType: models.ResourceTypeTarget,
				Action:       action,
				State:        models.JobStatePending,
			})
		},
		ServerApiUrl:     util.GetFrpcApiUrl(c.Frps.Protocol, c.Id, c.Frps.Domain),
		ServerVersion:    version,
		ServerUrl:        headscaleUrl,
		Provisioner:      provisioner,
		LoggerFactory:    loggerFactory,
		TelemetryService: telemetryService,
	})

	workspaceService := workspaces.NewWorkspaceService(workspaces.WorkspaceServiceConfig{
		WorkspaceStore:         workspaceStore,
		WorkspaceMetadataStore: workspaceMetadataStore,
		FindTarget: func(ctx context.Context, targetId string) (*models.Target, error) {
			t, err := targetService.GetTarget(ctx, &stores.TargetFilter{IdOrName: &targetId}, false)
			if err != nil {
				return nil, err
			}
			return &t.Target, nil
		},
		FindContainerRegistry: func(ctx context.Context, image string) (*models.ContainerRegistry, error) {
			return containerRegistryService.FindByImageName(image)
		},
		FindCachedBuild: func(ctx context.Context, w *models.Workspace) (*models.CachedBuild, error) {
			validStates := &[]models.BuildState{
				models.BuildStatePublished,
			}

			build, err := buildService.Find(&stores.BuildFilter{
				States:        validStates,
				RepositoryUrl: &w.Repository.Url,
				Branch:        &w.Repository.Branch,
				EnvVars:       &w.EnvVars,
				BuildConfig:   w.BuildConfig,
				GetNewest:     util.Pointer(true),
			})
			if err != nil {
				return nil, err
			}

			if build.Image == nil || build.User == nil {
				return nil, errors.New("cached build is missing image or user")
			}

			return &models.CachedBuild{
				User:  *build.User,
				Image: *build.Image,
			}, nil
		},
		GenerateApiKey: func(ctx context.Context, name string) (string, error) {
			return apiKeyService.Generate(models.ApiKeyTypeWorkspace, name)
		},
		RevokeApiKey: func(ctx context.Context, name string) error {
			return apiKeyService.Revoke(name)
		},
		ListGitProviderConfigs: func(ctx context.Context, repoUrl string) ([]*models.GitProviderConfig, error) {
			return gitProviderService.ListConfigsForUrl(repoUrl)
		},
		FindGitProviderConfig: func(ctx context.Context, id string) (*models.GitProviderConfig, error) {
			return gitProviderService.GetConfig(id)
		},
		GetLastCommitSha: func(ctx context.Context, repo *gitprovider.GitRepository) (string, error) {
			return gitProviderService.GetLastCommitSha(repo)
		},
		CreateJob: func(ctx context.Context, workspaceId string, action models.JobAction) error {
			return jobService.Save(&models.Job{
				ResourceId:   workspaceId,
				ResourceType: models.ResourceTypeWorkspace,
				Action:       action,
				State:        models.JobStatePending,
			})
		},
		TrackTelemetryEvent:   telemetryService.TrackServerEvent,
		ServerApiUrl:          util.GetFrpcApiUrl(c.Frps.Protocol, c.Id, c.Frps.Domain),
		ServerVersion:         version,
		ServerUrl:             headscaleUrl,
		DefaultWorkspaceImage: c.DefaultWorkspaceImage,
		DefaultWorkspaceUser:  c.DefaultWorkspaceUser,
		Provisioner:           provisioner,
		LoggerFactory:         loggerFactory,
	})

	return workspace_runner.NewWorkspaceJobRunner(workspace_runner.WorkspaceJobRunnerConfig{
		FindWorkspace: func(ctx context.Context, workspaceId string) (*models.Workspace, error) {
			workspaceDto, err := workspaceService.GetWorkspace(ctx, workspaceId, false)
			if err != nil {
				return nil, err
			}
			return &workspaceDto.Workspace, nil
		},
		HandleSuccessfulRemoval: func(ctx context.Context, workspaceId string) error {
			return workspaceService.HandleSuccessfulRemoval(ctx, workspaceId)
		},
		FindTarget: func(ctx context.Context, targetId string) (*models.Target, error) {
			targetDto, err := targetService.GetTarget(ctx, &stores.TargetFilter{IdOrName: &targetId}, false)
			if err != nil {
				return nil, err
			}
			return &targetDto.Target, nil
		},
		FindContainerRegistry: func(ctx context.Context, image string) (*models.ContainerRegistry, error) {
			return containerRegistryService.Find(image)
		},
		FindGitProviderConfig: func(ctx context.Context, id string) (*models.GitProviderConfig, error) {
			return gitProviderService.GetConfig(id)
		},
		TrackTelemetryEvent: func(event telemetry.ServerEvent, clientId string, props map[string]interface{}) error {
			return telemetryService.TrackServerEvent(event, clientId, props)
		},
		LoggerFactory: loggerFactory,
		Provisioner:   provisioner,
	}), nil
}

func GetTargetJobRunner(c *server.Config, configDir string, version string, telemetryService telemetry.TelemetryService) (runners.ITargetJobRunner, error) {
	dbPath, err := getDbPath()
	if err != nil {
		return nil, err
	}

	dbConnection := db.GetSQLiteConnection(dbPath)

	targetStore, err := db.NewTargetStore(dbConnection)
	if err != nil {
		return nil, err
	}

	targetMetadataStore, err := db.NewTargetMetadataStore(dbConnection)
	if err != nil {
		return nil, err
	}

	targetConfigStore, err := db.NewTargetConfigStore(dbConnection)
	if err != nil {
		return nil, err
	}
	targetConfigService := targetconfigs.NewTargetConfigService(targetconfigs.TargetConfigServiceConfig{
		TargetConfigStore: targetConfigStore,
	})

	jobStore, err := db.NewJobStore(dbConnection)
	if err != nil {
		return nil, err
	}

	jobService := jobs.NewJobService(jobs.JobServiceConfig{
		JobStore: jobStore,
	})

	apiKeyStore, err := db.NewApiKeyStore(dbConnection)
	if err != nil {
		return nil, err
	}

	apiKeyService := apikeys.NewApiKeyService(apikeys.ApiKeyServiceConfig{
		ApiKeyStore: apiKeyStore,
	})

	targetLogsDir, err := server.GetTargetLogsDir(configDir)
	if err != nil {
		return nil, err
	}
	buildLogsDir, err := build.GetBuildLogsDir()
	if err != nil {
		return nil, err
	}
	loggerFactory := logs.NewLoggerFactory(&targetLogsDir, &buildLogsDir)

	headscaleUrl := util.GetFrpcHeadscaleUrl(c.Frps.Protocol, c.Id, c.Frps.Domain)

	providerManager := manager.GetProviderManager(nil)

	provisioner := provisioner.NewProvisioner(provisioner.ProvisionerConfig{
		ProviderManager: providerManager,
	})

	targetService := targets.NewTargetService(targets.TargetServiceConfig{
		TargetStore:         targetStore,
		TargetMetadataStore: targetMetadataStore,
		FindTargetConfig: func(ctx context.Context, name string) (*models.TargetConfig, error) {
			return targetConfigService.Find(&stores.TargetConfigFilter{Name: &name})
		},
		GenerateApiKey: func(ctx context.Context, name string) (string, error) {
			return apiKeyService.Generate(models.ApiKeyTypeTarget, name)
		},
		RevokeApiKey: func(ctx context.Context, name string) error {
			return apiKeyService.Revoke(name)
		},
		CreateJob: func(ctx context.Context, targetId string, action models.JobAction) error {
			return jobService.Save(&models.Job{
				ResourceId:   targetId,
				ResourceType: models.ResourceTypeTarget,
				Action:       action,
				State:        models.JobStatePending,
			})
		},
		ServerApiUrl:     util.GetFrpcApiUrl(c.Frps.Protocol, c.Id, c.Frps.Domain),
		ServerVersion:    version,
		ServerUrl:        headscaleUrl,
		Provisioner:      provisioner,
		LoggerFactory:    loggerFactory,
		TelemetryService: telemetryService,
	})

	return target_runner.NewTargetJobRunner(target_runner.TargetJobRunnerConfig{
		FindTarget: func(ctx context.Context, targetId string) (*models.Target, error) {
			targetDto, err := targetService.GetTarget(ctx, &stores.TargetFilter{IdOrName: &targetId}, false)
			if err != nil {
				return nil, err
			}
			return &targetDto.Target, nil
		},
		HandleSuccessfulCreation: func(ctx context.Context, targetId string) error {
			return targetService.HandleSuccessfulCreation(ctx, targetId)
		},
		HandleSuccessfulRemoval: func(ctx context.Context, targetId string) error {
			return targetService.HandleSuccessfulRemoval(ctx, targetId)
		},
		TrackTelemetryEvent: func(event telemetry.ServerEvent, clientId string, props map[string]interface{}) error {
			return telemetryService.TrackServerEvent(event, clientId, props)
		},
		LoggerFactory: loggerFactory,
		Provisioner:   provisioner,
	}), nil
}
