// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	. "github.com/daytonaio/daytona/pkg/db/dto"
	"github.com/daytonaio/daytona/pkg/workspace"
)

type WorkspaceStore struct {
	db *gorm.DB
}

func NewWorkspaceStore(db *gorm.DB) (*WorkspaceStore, error) {
	err := db.AutoMigrate(&WorkspaceDTO{})
	if err != nil {
		return nil, err
	}

	return &WorkspaceStore{db: db}, nil
}

func (s *WorkspaceStore) List(filter *workspace.Filter) ([]*workspace.WorkspaceViewDTO, error) {
	workspaceDTOs := []WorkspaceDTO{}
	tx := processWorkspaceFilters(s.db, filter).Preload(clause.Associations).Find(&workspaceDTOs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	workspaceViewDTOs := []*workspace.WorkspaceViewDTO{}
	for _, workspaceDTO := range workspaceDTOs {
		workspaceViewDTOs = append(workspaceViewDTOs, ToWorkspaceViewDTO(workspaceDTO))
	}

	return workspaceViewDTOs, nil
}

func (w *WorkspaceStore) Find(filter *workspace.Filter) (*workspace.WorkspaceViewDTO, error) {
	workspaceDTO := WorkspaceDTO{}
	tx := processWorkspaceFilters(w.db, filter).Preload(clause.Associations).First(&workspaceDTO)
	if tx.Error != nil {
		if IsRecordNotFound(tx.Error) {
			return nil, workspace.ErrWorkspaceNotFound
		}
		return nil, tx.Error
	}

	return ToWorkspaceViewDTO(workspaceDTO), nil
}

func (s *WorkspaceStore) Save(workspace *workspace.Workspace) error {
	tx := s.db.Save(ToWorkspaceDTO(workspace))
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *WorkspaceStore) Delete(t *workspace.Workspace) error {
	tx := s.db.Delete(ToWorkspaceDTO(t))
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return workspace.ErrWorkspaceNotFound
	}

	return nil
}

func processWorkspaceFilters(tx *gorm.DB, filter *workspace.Filter) *gorm.DB {
	if filter != nil {
		if filter.IdOrName != nil {
			tx = tx.Where("id = ? OR name = ?", *filter.IdOrName, *filter.IdOrName)
		}
		if filter.States != nil && len(*filter.States) > 0 {
			placeholders := strings.Repeat("?,", len(*filter.States))
			placeholders = placeholders[:len(placeholders)-1]

			tx = tx.Where(fmt.Sprintf("state IN (%s)", placeholders), filter.StatesToInterface()...)
		}
	}

	return tx
}
