// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package dto

import "github.com/daytonaio/daytona/pkg/jobs/workspace"

type JobDTO struct {
	Id          string              `json:"id" validate:"required"`
	WorkspaceId string              `json:"workspaceId" validate:"required"`
	Action      workspace.JobAction `json:"action" validate:"required"`
	State       workspace.JobState  `json:"state" validate:"required"`
	Error       *string             `json:"error" validate:"optional"`
	CreatedAt   string              `json:"createdAt" validate:"required"`
} // @name JobDTO
