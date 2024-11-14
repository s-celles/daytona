// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package config

import "errors"

type TargetConfigStore interface {
	List(filter *Filter) ([]*TargetConfig, error)
	Find(filter *Filter) (*TargetConfig, error)
	Save(targetConfig *TargetConfig) error
	Delete(targetConfig *TargetConfig) error
}

type Filter struct {
	Name *string
}

var (
	ErrTargetConfigNotFound = errors.New("target config not found")
)

func IsTargetConfigNotFound(err error) bool {
	return err.Error() == ErrTargetConfigNotFound.Error()
}
