// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"github.com/daytonaio/daytona/pkg/apikey"
	. "github.com/daytonaio/daytona/pkg/db/dto"
	"gorm.io/gorm"
)

type ApiKeyStore struct {
	db *gorm.DB
}

func NewApiKeyStore(db *gorm.DB) (*ApiKeyStore, error) {
	err := db.AutoMigrate(&ApiKeyDTO{})
	if err != nil {
		return nil, err
	}

	return &ApiKeyStore{db: db}, nil
}

func (s *ApiKeyStore) List() ([]*apikey.ApiKey, error) {
	apiKeyDTOs := []ApiKeyDTO{}
	tx := s.db.Find(&apiKeyDTOs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	apiKeys := []*apikey.ApiKey{}
	for _, apiKeyDTO := range apiKeyDTOs {
		apiKey := ToApiKey(apiKeyDTO)
		apiKeys = append(apiKeys, &apiKey)
	}
	return apiKeys, nil
}

func (s *ApiKeyStore) Find(key string) (*apikey.ApiKey, error) {
	apiKeyDTO := ApiKeyDTO{}
	tx := s.db.Where("key_hash = ?", key).First(&apiKeyDTO)
	if tx.Error != nil {
		if IsRecordNotFound(tx.Error) {
			return nil, apikey.ErrApiKeyNotFound
		}
		return nil, tx.Error
	}

	apiKey := ToApiKey(apiKeyDTO)

	return &apiKey, nil
}

func (s *ApiKeyStore) FindByName(name string) (*apikey.ApiKey, error) {
	apiKeyDTO := ApiKeyDTO{}
	tx := s.db.Where("name = ?", name).First(&apiKeyDTO)
	if tx.Error != nil {
		if IsRecordNotFound(tx.Error) {
			return nil, apikey.ErrApiKeyNotFound
		}
		return nil, tx.Error
	}

	apiKey := ToApiKey(apiKeyDTO)

	return &apiKey, nil
}

func (s *ApiKeyStore) Save(apiKey *apikey.ApiKey) error {
	apiKeyDTO := ToApiKeyDTO(*apiKey)
	tx := s.db.Save(&apiKeyDTO)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *ApiKeyStore) Delete(apiKey *apikey.ApiKey) error {
	tx := s.db.Where("key_hash = ?", apiKey.KeyHash).Delete(&ApiKeyDTO{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return apikey.ErrApiKeyNotFound
	}

	return nil
}
