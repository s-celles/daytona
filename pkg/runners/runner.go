// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package runners

import (
	"context"
)

type IJobRunner interface {
	StartRunner(ctx context.Context) error
	CheckAndRunJobs(ctx context.Context) error
}
