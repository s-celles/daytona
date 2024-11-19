// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package target

import (
	"context"
	"fmt"

	"github.com/daytonaio/daytona/pkg/logs"
	"github.com/daytonaio/daytona/pkg/models"
	"github.com/daytonaio/daytona/pkg/views"
)

func (s *TargetJobRunner) Create(ctx context.Context, j *models.Job) error {
	tg, err := s.findTarget(ctx, j.ResourceId)
	if err != nil {
		return err
	}

	targetLogger := s.loggerFactory.CreateTargetLogger(tg.Id, tg.Name, logs.LogSourceServer)
	defer targetLogger.Close()

	targetLogger.Write([]byte(fmt.Sprintf("Creating target %s (%s)\n", tg.Name, tg.Id)))

	err = s.provisioner.CreateTarget(tg)
	if err != nil {
		return err
	}

	err = s.handleSuccessfulCreation(ctx, tg.Id)
	if err != nil {
		return err
	}

	targetLogger.Write([]byte(views.GetPrettyLogLine("Target creation complete")))
	return s.Start(ctx, j)
}
