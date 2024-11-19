// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"errors"
	"time"

	apiclient_util "github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/pkg/apiclient"
)

func AwaitWorkspaceState(workspaceId string, stateName apiclient.ModelsResourceStateName) error {
	for {
		ws, err := apiclient_util.GetWorkspace(workspaceId, false)
		if err != nil {
			return err
		}
		if ws.State.Name == stateName {
			return nil
		}
		if ws.State.Name == apiclient.ResourceStateNameError {
			var errorMessage string
			if ws.State.Error != nil {
				errorMessage = *ws.State.Error
			}
			return errors.New(errorMessage)
		}
		if ws.State.Name == apiclient.ResourceStateNameUnresponsive {
			return errors.New("workspace is unresponsive")
		}
		time.Sleep(time.Second)
	}
}

func AwaitTargetState(targetId string, stateName apiclient.ModelsResourceStateName) error {
	for {
		t, err := apiclient_util.GetTarget(targetId, false)
		if err != nil {
			return err
		}
		if t.State.Name == stateName || t.State.Name == apiclient.ResourceStateNameUndefined {
			return nil
		}
		if t.State.Name == apiclient.ResourceStateNameError {
			var errorMessage string
			if t.State.Error != nil {
				errorMessage = *t.State.Error
			}
			return errors.New(errorMessage)
		}
		if t.State.Name == apiclient.ResourceStateNameUnresponsive {
			return errors.New("target is unresponsive")
		}
		time.Sleep(time.Second)
	}
}
