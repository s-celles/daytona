// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspace

import (
	"context"

	"github.com/daytonaio/daytona/pkg/models"
)

func (s *WorkspaceJobRunner) Restart(ctx context.Context, j *models.Job) error {
	err := s.Stop(ctx, j)
	if err != nil {
		return err
	}

	return s.Start(ctx, j)
}
