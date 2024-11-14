// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package target

import "errors"

type Store interface {
	List(filter *Filter) ([]*TargetViewDTO, error)
	Find(filter *Filter) (*TargetViewDTO, error)
	Save(target *Target) error
	Delete(target *Target) error
}

type Filter struct {
	IdOrName *string
	Default  *bool
}

type TargetViewDTO struct {
	Target
	WorkspaceCount int `json:"workspaceCount" validate:"required"`
} // @name TargetViewDTO

var (
	ErrTargetNotFound = errors.New("target not found")
)

func IsTargetNotFound(err error) bool {
	return err.Error() == ErrTargetNotFound.Error()
}
