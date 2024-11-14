// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspace

import "errors"

type Store interface {
	List(filter *Filter) ([]*WorkspaceViewDTO, error)
	Find(filter *Filter) (*WorkspaceViewDTO, error)
	Save(workspace *Workspace) error
	Delete(workspace *Workspace) error
}

type Filter struct {
	IdOrName *string
	States   *[]WorkspaceState
}

func (f *Filter) StatesToInterface() []interface{} {
	args := make([]interface{}, len(*f.States))
	for i, v := range *f.States {
		args[i] = v
	}
	return args
}

type WorkspaceViewDTO struct {
	Workspace
	TargetName string `json:"targetName" validate:"required"`
} // @name WorkspaceViewDTO

var (
	ErrWorkspaceNotFound = errors.New("workspace not found")
)

func IsWorkspaceNotFound(err error) bool {
	return err.Error() == ErrWorkspaceNotFound.Error()
}
