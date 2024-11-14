// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package runner

import (
	"context"

	workspace_jobs "github.com/daytonaio/daytona/pkg/jobs/workspace"
	"github.com/daytonaio/daytona/pkg/server/workspaces"
)

type IRunner interface {
	Run(ctx context.Context, req interface{}) (response interface{}, err error)
}

type RunnerConfig struct {
	WorkspaceService workspaces.IWorkspaceService
	JobStore         workspace_jobs.Store
}

type Runner struct {
	workspaceService workspaces.IWorkspaceService
}

func NewRunner(config RunnerConfig) *Runner {
	return &Runner{
		workspaceService: config.WorkspaceService,
	}
}

func (r *Runner) Run(ctx context.Context, req interface{}) (response interface{}, err error) {
	return nil, nil
}
