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

func (s *TargetJobRunner) Stop(ctx context.Context, j *models.Job) error {
	t, err := s.findTarget(ctx, j.ResourceId)
	if err != nil {
		return err
	}

	targetLogger := s.loggerFactory.CreateTargetLogger(t.Id, t.Name, logs.LogSourceServer)
	defer targetLogger.Close()

	targetLogger.Write([]byte(fmt.Sprintf("Stopping target %s\n", t.Name)))

	//	todo: go routines
	err = s.provisioner.StopTarget(t)
	if err != nil {
		return err
	}

	targetLogger.Write([]byte(views.GetPrettyLogLine(fmt.Sprintf("Target %s stopped", t.Name))))

	return nil
}
