// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package jobs_test

import (
	"testing"

	"github.com/daytonaio/daytona/internal/testing/job"
	"github.com/daytonaio/daytona/pkg/jobs/workspace"
	"github.com/daytonaio/daytona/pkg/server/workspaces/jobs"
	"github.com/stretchr/testify/suite"
)

var job1 = &workspace.Job{
	Id:          "123",
	WorkspaceId: "123",
	Action:      workspace.JobActionCreate,
	State:       workspace.JobStatePending,
}

var job2 = &workspace.Job{
	Id:          "456",
	WorkspaceId: "123",
	Action:      workspace.JobActionCreate,
	State:       workspace.JobStatePending,
}

var job3 = &workspace.Job{
	Id:          "789",
	WorkspaceId: "123",
	Action:      workspace.JobActionCreate,
	State:       workspace.JobStatePending,
}

var expectedJobs []*workspace.Job
var expectedJobsMap map[string]*workspace.Job

type JobServiceTestSuite struct {
	suite.Suite
	jobService jobs.IJobService
	jobStore   workspace.Store
}

func NewJobServiceTestSuite() *JobServiceTestSuite {
	return &JobServiceTestSuite{}
}

func (s *JobServiceTestSuite) SetupTest() {
	expectedJobs = []*workspace.Job{
		job1, job2, job3,
	}

	expectedJobsMap = map[string]*workspace.Job{
		job1.Id: job1,
		job2.Id: job2,
		job3.Id: job3,
	}

	s.jobStore = job.NewInMemoryJobStore()
	s.jobService = jobs.NewJobService(jobs.JobServiceConfig{
		JobStore: s.jobStore,
	})

	for _, job := range expectedJobs {
		_ = s.jobStore.Save(job)
	}
}

func TestJobService(t *testing.T) {
	suite.Run(t, NewJobServiceTestSuite())
}

func (s *JobServiceTestSuite) TestList() {
	require := s.Require()

	jobs, err := s.jobService.List(nil)
	require.Nil(err)
	require.ElementsMatch(expectedJobs, jobs)
}

func (s *JobServiceTestSuite) TestFind() {
	require := s.Require()

	job, err := s.jobService.Find(&workspace.Filter{
		Id: &job1.Id,
	})
	require.Nil(err)
	require.Equal(job1, job)
}

func (s *JobServiceTestSuite) TestSave() {
	require := s.Require()

	expectedJobs = append(expectedJobs, job3)

	err := s.jobService.Save(job3)
	require.Nil(err)

	jobs, err := s.jobService.List(nil)
	require.Nil(err)
	require.ElementsMatch(expectedJobs, jobs)
}

func (s *JobServiceTestSuite) TestDelete() {
	require := s.Require()

	expectedJobs = expectedJobs[:2]

	err := s.jobService.Delete(job3)
	require.Nil(err)

	jobs, err := s.jobService.List(nil)
	require.Nil(err)
	require.ElementsMatch(expectedJobs, jobs)
}
