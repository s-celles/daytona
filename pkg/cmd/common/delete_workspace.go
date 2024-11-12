// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"context"
	"fmt"

	"github.com/daytonaio/daytona/cmd/daytona/config"
	apiclient_util "github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/pkg/apiclient"
	views_util "github.com/daytonaio/daytona/pkg/views/util"
)

func RemoveWorkspace(ctx context.Context, apiClient *apiclient.APIClient, workspace *apiclient.WorkspaceDTO, force bool) error {
	message := fmt.Sprintf("Deleting workspace %s", workspace.Name)
	err := views_util.WithInlineSpinner(message, func() error {
		res, err := apiClient.WorkspaceAPI.RemoveWorkspace(ctx, workspace.Id).Force(force).Execute()
		if err != nil {
			return apiclient_util.HandleErrorResponse(res, err)
		}
		c, err := config.GetConfig()
		if err != nil {
			return err
		}

		activeProfile, err := c.GetActiveProfile()
		if err != nil {
			return err
		}

		err = config.RemoveWorkspaceSshEntries(activeProfile.Id, workspace.Id)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
