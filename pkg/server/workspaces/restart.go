// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspaces

import (
	"context"

	"github.com/daytonaio/daytona/pkg/logs"
	"github.com/daytonaio/daytona/pkg/target"
	"github.com/daytonaio/daytona/pkg/telemetry"
	"github.com/daytonaio/daytona/pkg/workspace"
	log "github.com/sirupsen/logrus"
)

func (s *WorkspaceService) RestartWorkspace(ctx context.Context, workspaceId string) error {
	ws, err := s.workspaceStore.Find(&workspace.Filter{IdOrName: &workspaceId})
	if err != nil {
		return s.handleRestartError(ctx, &ws.Workspace, ErrWorkspaceNotFound)
	}

	ws.State = workspace.WorkspaceStatePendingRestart
	err = s.workspaceStore.Save(&ws.Workspace)
	return s.handleRestartError(ctx, &ws.Workspace, err)
}

func (s *WorkspaceService) RunRestartWorkspace(ctx context.Context, workspaceId string) error {
	ws, err := s.workspaceStore.Find(&workspace.Filter{IdOrName: &workspaceId})
	if err != nil {
		return s.handleRestartError(ctx, &ws.Workspace, ErrWorkspaceNotFound)
	}

	target, err := s.targetStore.Find(&target.Filter{IdOrName: &ws.TargetId})
	if err != nil {
		return s.handleRestartError(ctx, &ws.Workspace, err)
	}

	workspaceLogger := s.loggerFactory.CreateWorkspaceLogger(ws.Id, ws.Name, logs.LogSourceServer)
	defer workspaceLogger.Close()

	err = s.stopWorkspace(ctx, &ws.Workspace, &target.Target, workspaceLogger)
	if err != nil {
		return s.handleRestartError(ctx, &ws.Workspace, err)
	}

	err = s.startWorkspace(&ws.Workspace, &target.Target, workspaceLogger)
	if err != nil {
		return s.handleRestartError(ctx, &ws.Workspace, err)
	}

	err = s.workspaceStore.Save(&ws.Workspace)
	return s.handleRestartError(ctx, &ws.Workspace, err)
}

func (s *WorkspaceService) handleRestartError(ctx context.Context, w *workspace.Workspace, err error) error {
	if !telemetry.TelemetryEnabled(ctx) {
		return err
	}

	clientId := telemetry.ClientId(ctx)

	telemetryProps := telemetry.NewWorkspaceEventProps(ctx, w)
	event := telemetry.ServerEventWorkspaceStarted
	if err != nil {
		telemetryProps["error"] = err.Error()
		event = telemetry.ServerEventWorkspaceStartError
	}
	telemetryError := s.telemetryService.TrackServerEvent(event, clientId, telemetryProps)
	if telemetryError != nil {
		log.Trace(err)
	}

	return err
}
