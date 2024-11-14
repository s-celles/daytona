// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"gorm.io/gorm"

	. "github.com/daytonaio/daytona/pkg/db/dto"
	"github.com/daytonaio/daytona/pkg/gitprovider"
)

type GitProviderConfigStore struct {
	db *gorm.DB
}

func NewGitProviderConfigStore(db *gorm.DB) (*GitProviderConfigStore, error) {
	err := db.AutoMigrate(&GitProviderConfigDTO{})
	if err != nil {
		return nil, err
	}

	return &GitProviderConfigStore{db: db}, nil
}

func (s *GitProviderConfigStore) List() ([]*gitprovider.GitProviderConfig, error) {
	gitProviderDTOs := []GitProviderConfigDTO{}
	tx := s.db.Find(&gitProviderDTOs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	gitProviders := []*gitprovider.GitProviderConfig{}
	for _, gitProviderDTO := range gitProviderDTOs {
		gitProvider := ToGitProviderConfig(gitProviderDTO)
		gitProviders = append(gitProviders, &gitProvider)
	}

	return gitProviders, nil
}

func (s *GitProviderConfigStore) Find(id string) (*gitprovider.GitProviderConfig, error) {
	gitProviderDTO := GitProviderConfigDTO{}
	tx := s.db.Where("id = ?", id).First(&gitProviderDTO)
	if tx.Error != nil {
		if IsRecordNotFound(tx.Error) {
			return nil, gitprovider.ErrGitProviderConfigNotFound
		}
		return nil, tx.Error
	}

	gitProvider := ToGitProviderConfig(gitProviderDTO)

	return &gitProvider, nil
}

func (s *GitProviderConfigStore) Save(gitProvider *gitprovider.GitProviderConfig) error {
	gitProviderDTO := ToGitProviderConfigDTO(*gitProvider)
	tx := s.db.Save(&gitProviderDTO)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *GitProviderConfigStore) Delete(gitProvider *gitprovider.GitProviderConfig) error {
	gitProviderDTO := ToGitProviderConfigDTO(*gitProvider)
	tx := s.db.Delete(&gitProviderDTO)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gitprovider.ErrGitProviderConfigNotFound
	}

	return nil
}
