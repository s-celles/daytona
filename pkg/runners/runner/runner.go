// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package runner

import (
	"context"
	"errors"
	"time"

	"github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/pkg/build"
	"github.com/daytonaio/daytona/pkg/logs"
	"github.com/daytonaio/daytona/pkg/models"
	"github.com/daytonaio/daytona/pkg/provisioner"
	"github.com/daytonaio/daytona/pkg/runners"
	"github.com/daytonaio/daytona/pkg/scheduler"
	log "github.com/sirupsen/logrus"
)

type JobRunnerConfig struct {
	ListPendingJobs func(ctx context.Context) ([]*models.Job, error)
	UpdateJobState  func(ctx context.Context, job *models.Job, state models.JobState, err *error) error

	WorkspaceJobRunner runners.IWorkspaceJobRunner
	TargetJobRunner    runners.ITargetJobRunner
	LoggerFactory      logs.LoggerFactory
	Provisioner        provisioner.IProvisioner
}

func NewJobRunner(config JobRunnerConfig) runners.IJobRunner {
	return &JobRunner{
		listPendingJobs: config.ListPendingJobs,
		updateJobState:  config.UpdateJobState,

		workspaceJobRunner: config.WorkspaceJobRunner,
		targetJobRunner:    config.TargetJobRunner,
		loggerFactory:      config.LoggerFactory,
		provisioner:        config.Provisioner,
	}
}

type JobRunner struct {
	listPendingJobs func(ctx context.Context) ([]*models.Job, error)
	updateJobState  func(ctx context.Context, job *models.Job, state models.JobState, err *error) error

	workspaceJobRunner runners.IWorkspaceJobRunner
	targetJobRunner    runners.ITargetJobRunner
	loggerFactory      logs.LoggerFactory
	provisioner        provisioner.IProvisioner
}

func (s *JobRunner) StartRunner(ctx context.Context) error {
	scheduler := scheduler.NewCronScheduler()

	// Make sure the API is up
	for {
		_, err := apiclient.GetApiClient(nil)
		if err != nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	err := scheduler.AddFunc(build.DEFAULT_JOB_POLL_INTERVAL, func() {
		err := s.CheckAndRunJobs(ctx)
		if err != nil {
			log.Error(err)
		}
	})
	if err != nil {
		return err
	}

	scheduler.Start()
	return nil
}

func (s *JobRunner) CheckAndRunJobs(ctx context.Context) error {
	jobs, err := s.listPendingJobs(ctx)
	if err != nil {
		return err
	}

	// goroutines, sync group
	for _, job := range jobs {
		err = s.runJob(ctx, job)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *JobRunner) runJob(ctx context.Context, j *models.Job) error {
	err := s.updateJobState(ctx, j, models.JobStateRunning, nil)
	if err != nil {
		return err
	}

	switch j.ResourceType {
	case models.ResourceTypeTarget:
		err = s.runTargetJob(ctx, j)
	case models.ResourceTypeWorkspace:
		err = s.runWorkspaceJob(ctx, j)
	}

	if err != nil {
		return s.updateJobState(ctx, j, models.JobStateError, &err)
	}

	return s.updateJobState(ctx, j, models.JobStateSuccess, nil)
}

func (s *JobRunner) runWorkspaceJob(ctx context.Context, j *models.Job) error {
	switch j.Action {
	case models.JobActionCreate:
		return s.workspaceJobRunner.Create(ctx, j)
	case models.JobActionStart:
		return s.workspaceJobRunner.Start(ctx, j)
	case models.JobActionStop:
		return s.workspaceJobRunner.Stop(ctx, j)
	case models.JobActionRestart:
		return s.workspaceJobRunner.Restart(ctx, j)
	case models.JobActionDelete:
		return s.workspaceJobRunner.Delete(ctx, j, false)
	case models.JobActionForceDelete:
		return s.workspaceJobRunner.Delete(ctx, j, true)
	}
	return errors.New("invalid job action")
}

func (s *JobRunner) runTargetJob(ctx context.Context, j *models.Job) error {
	switch j.Action {
	case models.JobActionCreate:
		return s.targetJobRunner.Create(ctx, j)
	case models.JobActionStart:
		return s.targetJobRunner.Start(ctx, j)
	case models.JobActionStop:
		return s.targetJobRunner.Stop(ctx, j)
	case models.JobActionRestart:
		return s.targetJobRunner.Restart(ctx, j)
	case models.JobActionDelete:
		return s.targetJobRunner.Delete(ctx, j, false)
	case models.JobActionForceDelete:
		return s.targetJobRunner.Delete(ctx, j, true)
	}
	return errors.New("invalid job action")
}
