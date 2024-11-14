//go:build testing

// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"fmt"

	"github.com/daytonaio/daytona/pkg/jobs/workspace"
)

type InMemoryJobStore struct {
	jobs map[string]*workspace.Job
}

func NewInMemoryJobStore() workspace.Store {
	return &InMemoryJobStore{
		jobs: make(map[string]*workspace.Job),
	}
}

func (s *InMemoryJobStore) List(filter *workspace.Filter) ([]*workspace.Job, error) {
	jobs, err := s.processFilters(filter)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (s *InMemoryJobStore) Find(filter *workspace.Filter) (*workspace.Job, error) {
	jobs, err := s.processFilters(filter)
	if err != nil {
		return nil, err
	}
	if len(jobs) == 0 {
		return nil, workspace.ErrJobNotFound
	}

	return jobs[0], nil
}

func (s *InMemoryJobStore) Save(job *workspace.Job) error {
	s.jobs[job.Id] = job
	return nil
}

func (s *InMemoryJobStore) Delete(job *workspace.Job) error {
	delete(s.jobs, job.Id)
	return nil
}

func (s *InMemoryJobStore) processFilters(filter *workspace.Filter) ([]*workspace.Job, error) {
	var result []*workspace.Job
	filteredJobs := make(map[string]*workspace.Job)
	for k, v := range s.jobs {
		filteredJobs[k] = v
	}

	if filter != nil {
		if filter.Id != nil {
			job, ok := s.jobs[*filter.Id]
			if ok {
				return []*workspace.Job{job}, nil
			} else {
				return []*workspace.Job{}, fmt.Errorf("job with id %s not found", *filter.Id)
			}
		}
		if filter.States != nil {
			for _, job := range filteredJobs {
				check := false
				for _, state := range *filter.States {
					if job.State == state {
						check = true
						break
					}
				}
				if !check {
					delete(filteredJobs, job.Id)
				}
			}
		}
		if filter.Actions != nil {
			for _, job := range filteredJobs {
				check := false
				for _, action := range *filter.Actions {
					if job.Action == action {
						check = true
						break
					}
				}
				if !check {
					delete(filteredJobs, job.Id)
				}
			}
		}
	}

	for _, job := range filteredJobs {
		result = append(result, job)
	}

	return result, nil
}
