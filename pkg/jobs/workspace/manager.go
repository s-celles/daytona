// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspace

type IJobManager interface {
	CreateJob(job *Job) error
	GetJob(filter *Filter) (*Job, error)
	ListJobs(filter *Filter) ([]*Job, error)
	DeleteJob(job *Job) error
}

type JobManager struct {
	store Store
}

func NewJobManager(store Store) *JobManager {
	return &JobManager{
		store: store,
	}
}

func (m *JobManager) CreateJob(job *Job) error {
	return m.store.Save(job)
}

func (m *JobManager) GetJob(filter *Filter) (*Job, error) {
	return m.store.Find(filter)
}

func (m *JobManager) ListJobs(filter *Filter) ([]*Job, error) {
	return m.store.List(filter)
}

func (m *JobManager) DeleteJob(job *Job) error {
	return m.store.Delete(job)
}
