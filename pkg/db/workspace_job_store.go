// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"fmt"
	"strings"

	. "github.com/daytonaio/daytona/pkg/db/dto"
	"github.com/daytonaio/daytona/pkg/jobs/workspace"
	"gorm.io/gorm"
)

type WorkspaceJobStore struct {
	db *gorm.DB
}

func NewWorkspaceJobStore(db *gorm.DB) (*WorkspaceJobStore, error) {
	err := db.AutoMigrate(&WorkspaceJobDTO{})
	if err != nil {
		return nil, err
	}

	return &WorkspaceJobStore{db: db}, nil
}

func (s *WorkspaceJobStore) List(filter *workspace.Filter) ([]*workspace.Job, error) {
	workspaceJobDtos := []WorkspaceJobDTO{}
	tx := processWorkspaceJobFilters(s.db, filter).Find(&workspaceJobDtos)

	if tx.Error != nil {
		return nil, tx.Error
	}

	workspaceJobs := []*workspace.Job{}
	for _, workspaceJobDto := range workspaceJobDtos {
		workspaceJobs = append(workspaceJobs, ToWorkspaceJob(workspaceJobDto))
	}

	return workspaceJobs, nil
}

func (s *WorkspaceJobStore) Find(key string) (*workspace.Job, error) {
	workspaceJobDTO := WorkspaceJobDTO{}
	tx := s.db.Where("id = ?", key).First(&workspaceJobDTO)
	if tx.Error != nil {
		if IsRecordNotFound(tx.Error) {
			return nil, workspace.ErrJobNotFound
		}
		return nil, tx.Error
	}

	workspaceJob := ToWorkspaceJob(workspaceJobDTO)

	return workspaceJob, nil
}

func (s *WorkspaceJobStore) FindByName(name string) (*workspace.Job, error) {
	workspaceJobDTO := WorkspaceJobDTO{}
	tx := s.db.Where("name = ?", name).First(&workspaceJobDTO)
	if tx.Error != nil {
		if IsRecordNotFound(tx.Error) {
			return nil, workspace.ErrJobNotFound
		}
		return nil, tx.Error
	}

	workspaceJob := ToWorkspaceJob(workspaceJobDTO)

	return workspaceJob, nil
}

func (s *WorkspaceJobStore) Save(workspaceJob *workspace.Job) error {
	workspaceJobDTO := ToWorkspaceJobDTO(*workspaceJob)
	tx := s.db.Save(&workspaceJobDTO)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *WorkspaceJobStore) Delete(workspaceJob *workspace.Job) error {
	tx := s.db.Delete(ToWorkspaceJobDTO(*workspaceJob))
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return workspace.ErrJobNotFound
	}

	return nil
}

func processWorkspaceJobFilters(tx *gorm.DB, filter *workspace.Filter) *gorm.DB {
	if filter != nil {
		if filter.Id != nil {
			tx = tx.Where("id = ?", *filter.Id)
		}
		if filter.States != nil && len(*filter.States) > 0 {
			placeholders := strings.Repeat("?,", len(*filter.States))
			placeholders = placeholders[:len(placeholders)-1]

			tx = tx.Where(fmt.Sprintf("state IN (%s)", placeholders), filter.StatesToInterface()...)
		}
		if filter.Actions != nil && len(*filter.Actions) > 0 {
			placeholders := strings.Repeat("?,", len(*filter.Actions))
			placeholders = placeholders[:len(placeholders)-1]

			tx = tx.Where(fmt.Sprintf("action IN (%s)", placeholders), filter.ActionsToInterface()...)
		}
	}
	return tx
}
