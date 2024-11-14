// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package dto

import (
	"time"

	"github.com/daytonaio/daytona/pkg/jobs/workspace"
	"gvisor.dev/gvisor/pkg/errors"
)

type WorkspaceJobDTO struct {
	Id          string              `json:"id" gorm:"primaryKey"`
	WorkspaceId string              `json:"workspaceId" validate:"required"`
	Action      workspace.JobAction `json:"action" validate:"required"`
	State       workspace.JobState  `json:"state" validate:"required"`
	Error       *errors.Error       `json:"error" validate:"optional"`
	CreatedAt   time.Time           `json:"createdAt" validate:"required"`
	UpdatedAt   time.Time           `json:"updatedAt" validate:"required"`
}

func ToWorkspaceJobDTO(workspaceJob workspace.Job) WorkspaceJobDTO {
	return WorkspaceJobDTO{
		Id:          workspaceJob.Id,
		WorkspaceId: workspaceJob.WorkspaceId,
		Action:      workspaceJob.Action,
		State:       workspaceJob.State,
		Error:       workspaceJob.Error,
		CreatedAt:   workspaceJob.CreatedAt,
		UpdatedAt:   workspaceJob.UpdatedAt,
	}
}

func ToWorkspaceJob(workspaceJobDTO WorkspaceJobDTO) *workspace.Job {
	return &workspace.Job{
		Id:          workspaceJobDTO.Id,
		WorkspaceId: workspaceJobDTO.WorkspaceId,
		Action:      workspaceJobDTO.Action,
		State:       workspaceJobDTO.State,
		Error:       workspaceJobDTO.Error,
		CreatedAt:   workspaceJobDTO.CreatedAt,
		UpdatedAt:   workspaceJobDTO.UpdatedAt,
	}
}
