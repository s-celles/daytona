//go:build testing

// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspaces

import (
	"github.com/daytonaio/daytona/pkg/workspace"
)

type InMemoryWorkspaceStore struct {
	workspaces map[string]*workspace.Workspace
}

func NewInMemoryWorkspaceStore() workspace.Store {
	return &InMemoryWorkspaceStore{
		workspaces: make(map[string]*workspace.Workspace),
	}
}

func (s *InMemoryWorkspaceStore) List(filter *workspace.Filter) ([]*workspace.WorkspaceViewDTO, error) {
	workspaces, err := s.processFilters(filter)
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (s *InMemoryWorkspaceStore) Find(filter *workspace.Filter) (*workspace.WorkspaceViewDTO, error) {
	workspaces, err := s.processFilters(filter)
	if err != nil {
		return nil, err
	}
	if len(workspaces) == 0 {
		return nil, workspace.ErrWorkspaceNotFound
	}

	return workspaces[0], nil
}

func (s *InMemoryWorkspaceStore) Save(workspace *workspace.Workspace) error {
	s.workspaces[workspace.Id] = workspace
	return nil
}

func (s *InMemoryWorkspaceStore) Delete(workspace *workspace.Workspace) error {
	delete(s.workspaces, workspace.Id)
	return nil
}

func (s *InMemoryWorkspaceStore) processFilters(filter *workspace.Filter) ([]*workspace.WorkspaceViewDTO, error) {
	var result []*workspace.WorkspaceViewDTO
	filteredWorkspaces := make(map[string]*workspace.Workspace)
	for k, v := range s.workspaces {
		filteredWorkspaces[k] = v
	}

	if filter != nil {
		if filter.IdOrName != nil {
			for _, w := range filteredWorkspaces {
				if w.Id != *filter.IdOrName && w.Name != *filter.IdOrName {
					delete(filteredWorkspaces, w.Id)
				}
			}
		}
		if filter.States != nil {
			for _, w := range filteredWorkspaces {
				check := false
				for _, state := range *filter.States {
					if w.State == state {
						check = true
						break
					}
				}
				if !check {
					delete(filteredWorkspaces, w.Id)
				}
			}
		}
	}

	for _, w := range filteredWorkspaces {
		result = append(result, &workspace.WorkspaceViewDTO{Workspace: *w})
	}

	return result, nil
}
