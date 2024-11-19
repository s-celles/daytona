// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspace

import (
	"context"
	"fmt"

	"github.com/daytonaio/daytona/pkg/logs"
	"github.com/daytonaio/daytona/pkg/models"
	"github.com/daytonaio/daytona/pkg/views"
)

func (s *WorkspaceJobRunner) Start(ctx context.Context, j *models.Job) error {
	w, err := s.findWorkspace(ctx, j.ResourceId)
	if err != nil {
		return err
	}

	workspaceLogger := s.loggerFactory.CreateWorkspaceLogger(w.Id, w.Name, logs.LogSourceServer)
	defer workspaceLogger.Close()

	workspaceLogger.Write([]byte(fmt.Sprintf("Starting workspace %s\n", w.Name)))

	err = s.provisioner.StartWorkspace(w)
	if err != nil {
		return err
	}

	workspaceLogger.Write([]byte(views.GetPrettyLogLine(fmt.Sprintf("Workspace %s started", w.Name))))
	return nil
}
