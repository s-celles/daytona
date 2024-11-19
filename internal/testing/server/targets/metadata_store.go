// go:build testing

// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package targets

import (
	"github.com/daytonaio/daytona/pkg/models"
	"github.com/daytonaio/daytona/pkg/stores"
)

type InMemoryTargetMetadataStore struct {
	targetMetadataEntries map[string]*models.TargetMetadata
}

func NewInMemoryTargetMetadataStore() stores.TargetMetadataStore {
	return &InMemoryTargetMetadataStore{
		targetMetadataEntries: make(map[string]*models.TargetMetadata),
	}
}

func (s *InMemoryTargetMetadataStore) Find(filter *stores.TargetMetadataFilter) (*models.TargetMetadata, error) {
	for _, m := range s.targetMetadataEntries {
		if m.TargetId == *filter.TargetId {
			return m, nil
		}
	}

	return nil, stores.ErrTargetMetadataNotFound
}

func (s *InMemoryTargetMetadataStore) Save(targetMetadata *models.TargetMetadata) error {
	s.targetMetadataEntries[targetMetadata.TargetId] = targetMetadata
	return nil
}

func (s *InMemoryTargetMetadataStore) Delete(targetMetadata *models.TargetMetadata) error {
	delete(s.targetMetadataEntries, targetMetadata.TargetId)
	return nil
}

func (s *InMemoryTargetMetadataStore) processFilters(filter *stores.TargetMetadataFilter) ([]*models.TargetMetadata, error) {
	var result []*models.TargetMetadata
	filteredTargetMetadata := make(map[string]*models.TargetMetadata)
	for k, v := range s.targetMetadataEntries {
		filteredTargetMetadata[k] = v
	}

	if filter != nil {
		if filter.Id != nil {
			m, ok := s.targetMetadataEntries[*filter.Id]
			if ok {
				return []*models.TargetMetadata{m}, nil
			} else {
				return []*models.TargetMetadata{}, nil
			}
		}
		if filter.TargetId != nil {
			for _, m := range filteredTargetMetadata {
				if m.TargetId != *filter.TargetId {
					delete(filteredTargetMetadata, m.TargetId)
				}
			}
		}
	}

	for _, m := range filteredTargetMetadata {
		result = append(result, m)
	}

	return result, nil
}
