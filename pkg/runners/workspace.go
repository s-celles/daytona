// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package runners

import (
	"context"

	"github.com/daytonaio/daytona/pkg/models"
)

type IWorkspaceJobRunner interface {
	Create(ctx context.Context, job *models.Job) error
	Start(ctx context.Context, job *models.Job) error
	Stop(ctx context.Context, job *models.Job) error
	Restart(ctx context.Context, job *models.Job) error
	Delete(ctx context.Context, job *models.Job, force bool) error
}
