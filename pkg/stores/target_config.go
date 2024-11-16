// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package stores

import (
	"errors"

	"github.com/daytonaio/daytona/pkg/models"
)

type TargetConfigFilter struct {
	Name    *string
	Deleted *bool
}

type TargetConfigStore interface {
	List(filter *TargetConfigFilter) ([]*models.TargetConfig, error)
	Find(filter *TargetConfigFilter) (*models.TargetConfig, error)
	Save(targetConfig *models.TargetConfig) error
}

var (
	ErrTargetConfigNotFound = errors.New("target config not found")
	ErrTargetAlreadyExists  = errors.New("target already exists")
)

func IsTargetConfigNotFound(err error) bool {
	return err.Error() == ErrTargetConfigNotFound.Error()
}
