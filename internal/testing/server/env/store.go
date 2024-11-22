//go:build testing

// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package env

import (
	"github.com/daytonaio/daytona/pkg/models"
	"github.com/daytonaio/daytona/pkg/stores"
)

type InMemoryEnvironmentVariableStore struct {
	environmentVariables []*models.EnvironmentVariable
}

func NewInMemoryEnvironmentVariableStore() stores.EnvironmentVariableStore {
	return &InMemoryEnvironmentVariableStore{
		environmentVariables: make([]*models.EnvironmentVariable, 0),
	}
}

func (s *InMemoryEnvironmentVariableStore) List() ([]*models.EnvironmentVariable, error) {
	return s.environmentVariables, nil
}

func (s *InMemoryEnvironmentVariableStore) Save(environmentVariable *models.EnvironmentVariable) error {
	envVars := make([]*models.EnvironmentVariable, 0)
	for _, envVar := range s.environmentVariables {
		if envVar.Key != environmentVariable.Key {
			envVars = append(envVars, envVar)
		}
	}
	envVars = append(envVars, environmentVariable)
	s.environmentVariables = envVars
	return nil
}

func (s *InMemoryEnvironmentVariableStore) Delete(key string) error {
	envVars := make([]*models.EnvironmentVariable, 0)
	for _, envVar := range s.environmentVariables {
		if envVar.Key != key {
			envVars = append(envVars, envVar)
		}
	}
	s.environmentVariables = envVars
	return nil
}
