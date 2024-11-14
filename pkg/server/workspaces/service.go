// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspaces

import (
	"context"
	"io"

	workspace_jobs "github.com/daytonaio/daytona/pkg/jobs/workspace"
	"github.com/daytonaio/daytona/pkg/logs"
	"github.com/daytonaio/daytona/pkg/provisioner"
	"github.com/daytonaio/daytona/pkg/server/apikeys"
	"github.com/daytonaio/daytona/pkg/server/builds"
	"github.com/daytonaio/daytona/pkg/server/containerregistries"
	"github.com/daytonaio/daytona/pkg/server/gitproviders"
	"github.com/daytonaio/daytona/pkg/server/workspaces/dto"
	"github.com/daytonaio/daytona/pkg/target"
	"github.com/daytonaio/daytona/pkg/telemetry"
	"github.com/daytonaio/daytona/pkg/workspace"
)

type IWorkspaceService interface {
	CreateWorkspace(ctx context.Context, req dto.CreateWorkspaceDTO) (*workspace.WorkspaceViewDTO, error)
	GetWorkspace(ctx context.Context, filter *workspace.Filter, verbose bool) (*dto.WorkspaceDTO, error)
	ListWorkspaces(ctx context.Context, filter *workspace.Filter, verbose bool) ([]dto.WorkspaceDTO, error)
	StartWorkspace(ctx context.Context, workspaceId string) error
	StopWorkspace(ctx context.Context, workspaceId string) error
	RestartWorkspace(ctx context.Context, workspaceId string) error
	RemoveWorkspace(ctx context.Context, workspaceId string) error
	ForceRemoveWorkspace(ctx context.Context, workspaceId string) error

	RunCreateWorkspace(ctx context.Context, ws workspace.Workspace) error
	RunStartWorkspace(ctx context.Context, workspaceId string) error
	RunStopWorkspace(ctx context.Context, workspaceId string) error
	RunRestartWorkspace(ctx context.Context, workspaceId string) error
	RunRemoveWorkspace(ctx context.Context, workspaceId string) error
	RunForceRemoveWorkspace(ctx context.Context, workspaceId string) error

	GetWorkspaceLogReader(workspaceId string) (io.Reader, error)
	SetWorkspaceMetadata(workspaceId string, state *workspace.WorkspaceMetadata) (*workspace.WorkspaceViewDTO, error)
}

type targetStore interface {
	Find(filter *target.Filter) (*target.TargetViewDTO, error)
}

type WorkspaceServiceConfig struct {
	WorkspaceStore           workspace.Store
	TargetStore              targetStore
	ContainerRegistryService containerregistries.IContainerRegistryService
	BuildService             builds.IBuildService
	ServerApiUrl             string
	ServerUrl                string
	ServerVersion            string
	Provisioner              provisioner.IProvisioner
	DefaultWorkspaceImage    string
	DefaultWorkspaceUser     string
	ApiKeyService            apikeys.IApiKeyService
	LoggerFactory            logs.LoggerFactory
	GitProviderService       gitproviders.IGitProviderService
	TelemetryService         telemetry.TelemetryService
	JobStore                 workspace_jobs.Store
}

func NewWorkspaceService(config WorkspaceServiceConfig) IWorkspaceService {
	return &WorkspaceService{
		workspaceStore:           config.WorkspaceStore,
		targetStore:              config.TargetStore,
		containerRegistryService: config.ContainerRegistryService,
		buildService:             config.BuildService,
		serverApiUrl:             config.ServerApiUrl,
		serverUrl:                config.ServerUrl,
		serverVersion:            config.ServerVersion,
		defaultWorkspaceImage:    config.DefaultWorkspaceImage,
		defaultWorkspaceUser:     config.DefaultWorkspaceUser,
		provisioner:              config.Provisioner,
		loggerFactory:            config.LoggerFactory,
		apiKeyService:            config.ApiKeyService,
		gitProviderService:       config.GitProviderService,
		telemetryService:         config.TelemetryService,
		jobStore:                 config.JobStore,
	}
}

type WorkspaceService struct {
	workspaceStore           workspace.Store
	targetStore              targetStore
	containerRegistryService containerregistries.IContainerRegistryService
	buildService             builds.IBuildService
	provisioner              provisioner.IProvisioner
	apiKeyService            apikeys.IApiKeyService
	serverApiUrl             string
	serverUrl                string
	serverVersion            string
	defaultWorkspaceImage    string
	defaultWorkspaceUser     string
	loggerFactory            logs.LoggerFactory
	gitProviderService       gitproviders.IGitProviderService
	telemetryService         telemetry.TelemetryService
	jobStore                 workspace_jobs.Store
}

func (s *WorkspaceService) SetWorkspaceMetadata(workspaceId string, state *workspace.WorkspaceMetadata) (*workspace.WorkspaceViewDTO, error) {
	ws, err := s.workspaceStore.Find(&workspace.Filter{IdOrName: &workspaceId})
	if err != nil {
		return nil, ErrWorkspaceNotFound
	}

	ws.Metadata = state
	return ws, s.workspaceStore.Save(&ws.Workspace)
}

func (s *WorkspaceService) GetWorkspaceLogReader(workspaceId string) (io.Reader, error) {
	return s.loggerFactory.CreateWorkspaceLogReader(workspaceId)
}
