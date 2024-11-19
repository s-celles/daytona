// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package target

import (
	"context"

	"github.com/daytonaio/daytona/pkg/logs"
	"github.com/daytonaio/daytona/pkg/models"
	"github.com/daytonaio/daytona/pkg/provisioner"
	"github.com/daytonaio/daytona/pkg/runners"
	"github.com/daytonaio/daytona/pkg/telemetry"
)

type TargetJobRunnerConfig struct {
	FindTarget               func(ctx context.Context, targetId string) (*models.Target, error)
	HandleSuccessfulCreation func(ctx context.Context, targetId string) error
	HandleSuccessfulRemoval  func(ctx context.Context, targetId string) error

	TrackTelemetryEvent func(event telemetry.ServerEvent, clientId string, props map[string]interface{}) error

	LoggerFactory logs.LoggerFactory
	Provisioner   provisioner.IProvisioner
}

func NewTargetJobRunner(config TargetJobRunnerConfig) runners.ITargetJobRunner {
	return &TargetJobRunner{
		findTarget:               config.FindTarget,
		handleSuccessfulCreation: config.HandleSuccessfulCreation,
		handleSuccessfulRemoval:  config.HandleSuccessfulRemoval,

		trackTelemetryEvent: config.TrackTelemetryEvent,

		loggerFactory: config.LoggerFactory,
		provisioner:   config.Provisioner,
	}
}

type TargetJobRunner struct {
	findTarget               func(ctx context.Context, targetId string) (*models.Target, error)
	handleSuccessfulCreation func(ctx context.Context, targetId string) error
	handleSuccessfulRemoval  func(ctx context.Context, targetId string) error

	trackTelemetryEvent func(event telemetry.ServerEvent, clientId string, props map[string]interface{}) error

	loggerFactory logs.LoggerFactory
	provisioner   provisioner.IProvisioner
}
