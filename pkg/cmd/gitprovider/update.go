// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package gitprovider

import (
	"context"

	"github.com/daytonaio/daytona/internal/util"
	apiclient_util "github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/cmd/common"
	"github.com/daytonaio/daytona/pkg/views"
	gitprovider_view "github.com/daytonaio/daytona/pkg/views/gitprovider"
	"github.com/daytonaio/daytona/pkg/views/selection"
	views_util "github.com/daytonaio/daytona/pkg/views/util"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a Git provider",
	Aliases: common.GetAliases("update"),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		apiClient, err := apiclient_util.GetApiClient(nil)
		if err != nil {
			return err
		}

		gitProviders, res, err := apiClient.GitProviderAPI.ListGitProviders(ctx).Execute()
		if err != nil {
			return apiclient_util.HandleErrorResponse(res, err)
		}

		if len(gitProviders) == 0 {
			views_util.NotifyEmptyGitProviderList(true)
			return nil
		}

		selectedGitProvider := selection.GetGitProviderConfigFromPrompt(selection.GetGitProviderConfigParams{
			GitProviderConfigs: gitProviders,
			ActionVerb:         "Update",
		})

		if selectedGitProvider == nil {
			return nil
		}

		existingAliases := util.ArrayMap(gitProviders, func(gp apiclient.GitProvider) string {
			return gp.Alias
		})

		setGitProviderConfig := apiclient.SetGitProviderConfig{
			Id:            &selectedGitProvider.Id,
			ProviderId:    selectedGitProvider.ProviderId,
			Token:         selectedGitProvider.Token,
			BaseApiUrl:    selectedGitProvider.BaseApiUrl,
			Username:      &selectedGitProvider.Username,
			Alias:         &selectedGitProvider.Alias,
			SigningMethod: selectedGitProvider.SigningMethod,
			SigningKey:    selectedGitProvider.SigningKey,
		}

		flags := map[string]string{}
		err = gitprovider_view.GitProviderCreationView(ctx, apiClient, &setGitProviderConfig, existingAliases, flags)
		if err != nil {
			return err
		}

		res, err = apiClient.GitProviderAPI.SaveGitProvider(ctx).GitProviderConfig(setGitProviderConfig).Execute()
		if err != nil {
			return apiclient_util.HandleErrorResponse(res, err)
		}

		views.RenderInfoMessage("Git provider has been updated")
		return nil
	},
}
