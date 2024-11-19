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

func (s *TargetJobRunner) Start(ctx context.Context, j *models.Job) error {
	tg, err := s.findTarget(ctx, j.ResourceId)
	if err != nil {
		return err
	}

	targetLogger := s.loggerFactory.CreateTargetLogger(j.ResourceId, "", logs.LogSourceServer)
	defer targetLogger.Close()

	targetLogger.Write([]byte("Starting target\n"))

	err = s.provisioner.StartTarget(tg)
	if err != nil {
		return err
	}

	targetLogger.Write([]byte(views.GetPrettyLogLine(fmt.Sprintf("Target %s started", tg.Name))))

	return nil
}
