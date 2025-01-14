// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package targetconfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"

	"github.com/daytonaio/daytona/cmd/daytona/config"
	"github.com/daytonaio/daytona/internal/util"
	internal_util "github.com/daytonaio/daytona/internal/util"
	apiclient_util "github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/internal/util/apiclient/conversion"
	"github.com/daytonaio/daytona/pkg/apiclient"
	cmd_common "github.com/daytonaio/daytona/pkg/cmd/common"
	"github.com/daytonaio/daytona/pkg/cmd/provider"
	"github.com/daytonaio/daytona/pkg/common"
	"github.com/daytonaio/daytona/pkg/views"
	provider_view "github.com/daytonaio/daytona/pkg/views/provider"
	"github.com/daytonaio/daytona/pkg/views/targetconfig"
	"github.com/spf13/cobra"
)

var pipeFile string

var TargetConfigCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create target config",
	Args:    cobra.NoArgs,
	Aliases: cmd_common.GetAliases("create"),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		var input []byte
		var err error

		if pipeFile == "-" {
			input, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			return handleTargetConfigJSON(input)
		} else if pipeFile != "" {
			input, err = os.ReadFile(pipeFile)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", pipeFile, err)
			}
			return handleTargetConfigJSON(input)
		}

		apiClient, err := apiclient_util.GetApiClient(nil)
		if err != nil {
			return err
		}

		c, err := config.GetConfig()
		if err != nil {
			return err
		}

		activeProfile, err := c.GetActiveProfile()
		if err != nil {
			return err
		}

		targetConfig, err := TargetConfigCreationFlow(ctx, apiClient, activeProfile.Name)
		if err != nil {
			return err
		}

		if targetConfig == nil {
			return nil
		}

		views.RenderInfoMessage(fmt.Sprintf("Target config '%s' set successfully", targetConfig.Name))
		return nil
	},
}

func TargetConfigCreationFlow(ctx context.Context, apiClient *apiclient.APIClient, activeProfileName string) (*targetconfig.TargetConfigView, error) {
	serverConfig, res, err := apiClient.ServerAPI.GetConfigExecute(apiclient.ApiGetConfigRequest{})
	if err != nil {
		return nil, apiclient_util.HandleErrorResponse(res, err)
	}

	providersManifest, err := util.GetProvidersManifest(serverConfig.RegistryUrl)
	if err != nil {
		return nil, err
	}

	var latestProviders []apiclient.ProviderInfo
	if providersManifest != nil {
		providersManifestLatest := providersManifest.GetLatestVersions()
		if providersManifestLatest == nil {
			return nil, errors.New("could not get latest provider versions")
		}

		latestProviders = conversion.GetProviderListFromManifest(providersManifestLatest)
	} else {
		fmt.Println("Could not get provider manifest. Can't check for new providers to install")
	}

	providerViewList, err := provider.GetProviderViewOptions(ctx, apiClient, latestProviders)
	if err != nil {
		return nil, err
	}

	selectedProvider, err := provider_view.GetProviderFromPrompt(providerViewList, "Choose a Provider", false)
	if err != nil {
		if common.IsCtrlCAbort(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	if selectedProvider == nil {
		return nil, nil
	}

	selectedTargetConfig := &targetconfig.TargetConfigView{
		Name:    "",
		Options: "{}",
		ProviderInfo: targetconfig.ProviderInfo{
			Name:            selectedProvider.Name,
			AgentlessTarget: selectedProvider.AgentlessTarget,
			RunnerId:        selectedProvider.RunnerId,
			RunnerName:      selectedProvider.RunnerName,
			Version:         selectedProvider.Version,
			Label:           selectedProvider.Label,
		},
	}

	if selectedProvider.Installed != nil && !*selectedProvider.Installed {
		if providersManifest == nil {
			return nil, errors.New("could not get providers manifest")
		}

		selectedRunner, err := cmd_common.GetRunnerFlow(apiClient, "Manage Providers")
		if err != nil {
			if common.IsCtrlCAbort(err) {
				return nil, nil
			} else {
				return nil, err
			}
		}

		if selectedRunner == nil {
			return nil, nil
		}

		err = provider.InstallProvider(apiClient, selectedRunner.Id, *selectedProvider, providersManifest)
		if err != nil {
			return nil, err
		}

		selectedProvider.RunnerId = selectedRunner.Id
		selectedProvider.RunnerName = selectedRunner.Name
	}

	targetConfigs, res, err := apiClient.TargetConfigAPI.ListTargetConfigs(ctx).Execute()
	if err != nil {
		return nil, apiclient_util.HandleErrorResponse(res, err)
	}
	selectedTargetConfig.Name = ""
	err = targetconfig.NewTargetConfigNameInput(&selectedTargetConfig.Name, internal_util.ArrayMap(targetConfigs, func(t apiclient.TargetConfig) string {
		return t.Name
	}))
	if err != nil {
		return nil, err
	}

	err = targetconfig.CreateTargetConfigForm(selectedTargetConfig, selectedProvider.TargetConfigManifest)
	if err != nil {
		return nil, err
	}

	targetConfigData := apiclient.CreateTargetConfigDTO{
		Name:    selectedTargetConfig.Name,
		Options: selectedTargetConfig.Options,
		ProviderInfo: apiclient.ProviderInfo{
			AgentlessTarget:      selectedProvider.AgentlessTarget,
			Name:                 selectedProvider.Name,
			Version:              selectedProvider.Version,
			Label:                selectedProvider.Label,
			RunnerId:             selectedProvider.RunnerId,
			RunnerName:           selectedProvider.RunnerName,
			TargetConfigManifest: selectedProvider.TargetConfigManifest,
		},
	}

	targetConfig, res, err := apiClient.TargetConfigAPI.CreateTargetConfig(context.Background()).TargetConfig(targetConfigData).Execute()
	if err != nil {
		return nil, apiclient_util.HandleErrorResponse(res, err)
	}

	return &targetconfig.TargetConfigView{
		Id:      targetConfig.Id,
		Name:    targetConfig.Name,
		Options: targetConfig.Options,
		ProviderInfo: targetconfig.ProviderInfo{
			Name:            targetConfig.ProviderInfo.Name,
			AgentlessTarget: targetConfig.ProviderInfo.AgentlessTarget,
			RunnerId:        targetConfig.ProviderInfo.RunnerId,
			RunnerName:      targetConfig.ProviderInfo.RunnerName,
			Version:         targetConfig.ProviderInfo.Version,
			Label:           targetConfig.ProviderInfo.Label,
		},
	}, nil
}

func handleTargetConfigJSON(data []byte) error {
	ctx := context.Background()

	apiClient, err := apiclient_util.GetApiClient(nil)
	if err != nil {
		return err
	}
	var selectedTargetConfig *targetconfig.TargetConfigView
	err = parseJSON(data, &selectedTargetConfig)
	if err != nil {
		return fmt.Errorf("failed to parse input: %w", err)
	}

	if selectedTargetConfig.Name == "" {
		return errors.New("invalid input: 'name' field is required")
	}
	if selectedTargetConfig.Options == "" {
		return errors.New("option fields are required to setup your target config")
	}
	targetManifest := selectedTargetConfig.ProviderInfo.TargetConfigManifest

	err = validateProperty(targetManifest, selectedTargetConfig)
	if err != nil {
		return err
	}
	targetConfigData := apiclient.CreateTargetConfigDTO{
		Name:    selectedTargetConfig.Name,
		Options: selectedTargetConfig.Options,
		ProviderInfo: apiclient.ProviderInfo{
			Name:    selectedTargetConfig.ProviderInfo.Name,
			Version: selectedTargetConfig.ProviderInfo.Version,
		},
	}
	_, res, err := apiClient.TargetConfigAPI.CreateTargetConfig(ctx).TargetConfig(targetConfigData).Execute()
	if err != nil {
		return apiclient_util.HandleErrorResponse(res, err)
	}
	views.RenderInfoMessage("Target set successfully and will be used by default")
	return nil
}

func parseJSON(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err == nil {
		return nil
	}
	return errors.New("input is not a valid JSON")
}

func validateProperty(targetManifest map[string]apiclient.TargetConfigProperty, targetConfig *targetconfig.TargetConfigView) error {
	optionMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(targetConfig.Options), &optionMap); err != nil {
		return fmt.Errorf("failed to parse options JSON: %w", err)
	}
	for optionKey := range optionMap {
		if _, exists := targetManifest[optionKey]; !exists {
			return fmt.Errorf("invalid property '%s' for target manifest '%s'", optionKey, targetConfig.Name)
		}
	}

	sortedKeys := make([]string, 0, len(targetManifest))
	for k := range targetManifest {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, name := range sortedKeys {
		if _, present := optionMap[name]; !present {
			continue
		}

		property := targetManifest[name]
		if property.DisabledPredicate != nil && *property.DisabledPredicate != "" {
			if matched, err := regexp.Match(*property.DisabledPredicate, []byte(targetConfig.Name)); err == nil && matched {
				if !contains(property.Options, optionMap[name]) {
					return fmt.Errorf("unexpected property '%s' for target manifest '%s'", name, targetConfig.Name)
				}
				continue
			}
		}

		switch *property.Type {
		case apiclient.TargetConfigPropertyTypeFloat, apiclient.TargetConfigPropertyTypeInt:
			_, isNumber := optionMap[name].(float64)
			if !isNumber {
				return fmt.Errorf("invalid type for %s, expected number", name)
			}

		case apiclient.TargetConfigPropertyTypeString:
			_, isString := optionMap[name].(string)
			if !isString {
				return fmt.Errorf("invalid type for %s, expected string", name)
			}

		case apiclient.TargetConfigPropertyTypeBoolean:
			_, isBool := optionMap[name].(bool)
			if !isBool {
				return fmt.Errorf("invalid type for %s, expected boolean", name)
			}

		case apiclient.TargetConfigPropertyTypeOption:
			optionValue, ok := optionMap[name].(string)
			if !ok {
				return fmt.Errorf("invalid value for '%s': expected a string", name)
			}
			valid := false
			for _, allowedOption := range property.Options {
				if optionValue == allowedOption {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("unexpected property '%s' for target manifest '%s' : valid properties are %v", optionValue, targetConfig.Name, property.Options)
			}

		case apiclient.TargetConfigPropertyTypeFilePath:
			_, isString := optionMap[name].(string)
			if !isString {
				return fmt.Errorf("invalid type for %s, expected file path string", name)
			}

		default:
			return fmt.Errorf("unsupported provider type: %s", *property.Type)
		}
	}
	return nil
}

func contains(slice []string, item interface{}) bool {
	for _, val := range slice {
		if val == item {
			return true
		}
	}
	return false
}

func init() {
	TargetConfigCreateCmd.Flags().StringVarP(&pipeFile, "file", "f", "", "Path to JSON file for target configuration, use '-' to read from stdin")
}