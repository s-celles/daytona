// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package env

import (
	"context"
	"fmt"
	"strings"

	apiclient_util "github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/views"
	"github.com/daytonaio/daytona/pkg/views/env"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add [KEY=VALUE]...",
	Short:   "Add environment variables",
	Aliases: []string{"a", "new"},
	RunE: func(cmd *cobra.Command, args []string) error {
		envVarsMap := make(map[string]string)

		if len(args) > 0 {
			for _, arg := range args {
				kv := strings.Split(arg, "=")
				if len(kv) != 2 {
					return fmt.Errorf("invalid key-value pair: %s", arg)
				}
				envVarsMap[kv[0]] = kv[1]
			}
		} else {
			err := env.AddEnvVarsView(&envVarsMap)
			if err != nil {
				return err
			}
		}

		ctx := context.Background()

		apiClient, err := apiclient_util.GetApiClient(nil)
		if err != nil {
			return err
		}

		for key, value := range envVarsMap {
			res, err := apiClient.EnvVarAPI.AddEnvironmentVariable(ctx).EnvironmentVariable(apiclient.EnvironmentVariable{
				Key:   key,
				Value: value,
			}).Execute()
			if err != nil {
				log.Error(apiclient_util.HandleErrorResponse(res, err))
			}
		}

		views.RenderInfoMessageBold("Profile environment variables have been successfully added")

		return nil
	},
}
