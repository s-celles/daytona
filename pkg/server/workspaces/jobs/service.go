// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package jobs

import (
	workspace_job "github.com/daytonaio/daytona/pkg/jobs/workspace"
)

type IJobService interface {
	Delete(job *workspace_job.Job) error
	Find(filter *workspace_job.Filter) (*workspace_job.Job, error)
	List(filter *workspace_job.Filter) ([]*workspace_job.Job, error)
	Save(job *workspace_job.Job) error
}

type JobServiceConfig struct {
	JobStore workspace_job.Store
}

type JobService struct {
	jobStore workspace_job.Store
}

func NewJobService(config JobServiceConfig) IJobService {
	return &JobService{
		jobStore: config.JobStore,
	}
}

func (s *JobService) List(filter *workspace_job.Filter) ([]*workspace_job.Job, error) {
	return s.jobStore.List(filter)
}

func (s *JobService) Find(filter *workspace_job.Filter) (*workspace_job.Job, error) {
	return s.jobStore.Find(filter)
}

func (s *JobService) Save(job *workspace_job.Job) error {
	return s.jobStore.Save(job)
}

func (s *JobService) Delete(job *workspace_job.Job) error {
	return s.jobStore.Delete(job)
}
