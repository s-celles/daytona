// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspace

import (
	"context"
	"fmt"

	"github.com/daytonaio/daytona/pkg/logs"
	"github.com/daytonaio/daytona/pkg/models"
	"github.com/daytonaio/daytona/pkg/stores"
	"github.com/daytonaio/daytona/pkg/views"
)

func (s *WorkspaceJobRunner) Create(ctx context.Context, j *models.Job) error {
	w, err := s.findWorkspace(ctx, j.ResourceId)
	if err != nil {
		return err
	}

	workspaceLogger := s.loggerFactory.CreateWorkspaceLogger(w.Id, w.Name, logs.LogSourceServer)
	defer workspaceLogger.Close()

	workspaceLogger.Write([]byte(fmt.Sprintf("Creating workspace %s\n", w.Name)))

	cr, err := s.findContainerRegistry(ctx, w.Image)
	if err != nil && !stores.IsContainerRegistryNotFound(err) {
		return err
	}

	var gc *models.GitProviderConfig

	if w.GitProviderConfigId != nil {
		gc, err = s.findGitProviderConfig(ctx, *w.GitProviderConfigId)
		if err != nil && !stores.IsGitProviderNotFound(err) {
			return err
		}
	}

	err = s.provisioner.CreateWorkspace(w, cr, gc)
	if err != nil {
		return err
	}

	workspaceLogger.Write([]byte(views.GetPrettyLogLine(fmt.Sprintf("Workspace %s created", w.Name))))
	return s.Start(ctx, j)
}
