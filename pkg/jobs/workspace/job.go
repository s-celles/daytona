// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspace

import (
	"time"

	"gvisor.dev/gvisor/pkg/errors"
)

type Job struct {
	Id          string        `json:"id" validate:"required"`
	WorkspaceId string        `json:"workspaceId" validate:"required"`
	Action      JobAction     `json:"action" validate:"required"`
	State       JobState      `json:"state" validate:"required"`
	Error       *errors.Error `json:"error" validate:"optional"`
	CreatedAt   time.Time     `json:"createdAt" validate:"required"`
	UpdatedAt   time.Time     `json:"updatedAt" validate:"required"`
} // @name WorkspaceJob

type JobState string // @name JobState

const (
	JobStatePending JobState = "pending"
	JobStateRunning JobState = "running"
	JobStateError   JobState = "error"
	JobStateSuccess JobState = "success"
)

type JobAction string

const (
	JobActionCreate      JobAction = "create"
	JobActionStart       JobAction = "start"
	JobActionStop        JobAction = "stop"
	JobActionDelete      JobAction = "delete"
	JobActionForceDelete JobAction = "force-delete"
	JobActionRestart     JobAction = "restart"
)
