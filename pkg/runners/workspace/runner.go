// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspace

import (
	"context"

	"github.com/daytonaio/daytona/pkg/logs"
	"github.com/daytonaio/daytona/pkg/models"
	"github.com/daytonaio/daytona/pkg/provisioner"
	"github.com/daytonaio/daytona/pkg/runners"
	"github.com/daytonaio/daytona/pkg/telemetry"
)

type WorkspaceJobRunnerConfig struct {
	FindWorkspace           func(ctx context.Context, workspaceId string) (*models.Workspace, error)
	HandleSuccessfulRemoval func(ctx context.Context, workspaceId string) error
	FindTarget              func(ctx context.Context, targetId string) (*models.Target, error)
	FindContainerRegistry   func(ctx context.Context, image string) (*models.ContainerRegistry, error)
	FindGitProviderConfig   func(ctx context.Context, id string) (*models.GitProviderConfig, error)
	TrackTelemetryEvent     func(event telemetry.ServerEvent, clientId string, props map[string]interface{}) error

	LoggerFactory logs.LoggerFactory
	Provisioner   provisioner.IProvisioner
}

func NewWorkspaceJobRunner(config WorkspaceJobRunnerConfig) runners.IWorkspaceJobRunner {
	return &WorkspaceJobRunner{
		findWorkspace:           config.FindWorkspace,
		handleSuccessfulRemoval: config.HandleSuccessfulRemoval,
		findTarget:              config.FindTarget,
		findContainerRegistry:   config.FindContainerRegistry,
		findGitProviderConfig:   config.FindGitProviderConfig,
		trackTelemetryEvent:     config.TrackTelemetryEvent,

		loggerFactory: config.LoggerFactory,
		provisioner:   config.Provisioner,
	}
}

type WorkspaceJobRunner struct {
	findWorkspace           func(ctx context.Context, workspaceId string) (*models.Workspace, error)
	handleSuccessfulRemoval func(ctx context.Context, workspaceId string) error
	findTarget              func(ctx context.Context, targetId string) (*models.Target, error)
	findContainerRegistry   func(ctx context.Context, image string) (*models.ContainerRegistry, error)
	findGitProviderConfig   func(ctx context.Context, id string) (*models.GitProviderConfig, error)
	trackTelemetryEvent     func(event telemetry.ServerEvent, clientId string, props map[string]interface{}) error

	loggerFactory logs.LoggerFactory
	provisioner   provisioner.IProvisioner
}
