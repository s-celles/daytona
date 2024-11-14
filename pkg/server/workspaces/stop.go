// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspaces

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/daytonaio/daytona/pkg/logs"
	"github.com/daytonaio/daytona/pkg/target"
	"github.com/daytonaio/daytona/pkg/telemetry"
	"github.com/daytonaio/daytona/pkg/views"
	"github.com/daytonaio/daytona/pkg/workspace"
	log "github.com/sirupsen/logrus"
)

func (s *WorkspaceService) StopWorkspace(ctx context.Context, workspaceId string) error {
	ws, err := s.workspaceStore.Find(&workspace.Filter{IdOrName: &workspaceId})
	if err != nil {
		return s.handleStopError(ctx, &ws.Workspace, ErrWorkspaceNotFound)
	}

	ws.State = workspace.WorkspaceStateStopping
	err = s.workspaceStore.Save(&ws.Workspace)
	return s.handleStopError(ctx, &ws.Workspace, err)
}

func (s *WorkspaceService) RunStopWorkspace(ctx context.Context, workspaceId string) error {
	ws, err := s.workspaceStore.Find(&workspace.Filter{IdOrName: &workspaceId})
	if err != nil {
		return s.handleStopError(ctx, &ws.Workspace, ErrWorkspaceNotFound)
	}

	target, err := s.targetStore.Find(&target.Filter{IdOrName: &ws.TargetId})
	if err != nil {
		return s.handleStopError(ctx, &ws.Workspace, err)
	}

	workspaceLogger := s.loggerFactory.CreateWorkspaceLogger(ws.Id, ws.Name, logs.LogSourceServer)
	defer workspaceLogger.Close()

	err = s.stopWorkspace(ctx, &ws.Workspace, &target.Target, workspaceLogger)
	if err != nil {
		return s.handleStartError(ctx, &ws.Workspace, err)
	}

	err = s.workspaceStore.Save(&ws.Workspace)
	return s.handleStopError(ctx, &ws.Workspace, err)
}

func (s *WorkspaceService) stopWorkspace(ctx context.Context, w *workspace.Workspace, target *target.Target, logger io.Writer) error {
	logger.Write([]byte(fmt.Sprintf("Stopping workspace %s\n", w.Name)))

	err := s.provisioner.StopWorkspace(w, target)
	if err != nil {
		return s.handleStopError(ctx, w, err)
	}

	if w.Metadata != nil {
		w.Metadata.Uptime = 0
		w.Metadata.UpdatedAt = time.Now().Format(time.RFC1123)
	}

	logger.Write([]byte(views.GetPrettyLogLine(fmt.Sprintf("Workspace %s stopped", w.Name))))

	return nil
}

func (s *WorkspaceService) handleStopError(ctx context.Context, w *workspace.Workspace, err error) error {
	if !telemetry.TelemetryEnabled(ctx) {
		return err
	}

	clientId := telemetry.ClientId(ctx)

	telemetryProps := telemetry.NewWorkspaceEventProps(ctx, w)
	event := telemetry.ServerEventWorkspaceStopped
	if err != nil {
		telemetryProps["error"] = err.Error()
		event = telemetry.ServerEventWorkspaceStopError
	}
	telemetryError := s.telemetryService.TrackServerEvent(event, clientId, telemetryProps)
	if telemetryError != nil {
		log.Trace(err)
	}

	return err
}
