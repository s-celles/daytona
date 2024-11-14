// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package target

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	apiclient_util "github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/cmd/common"
	"github.com/daytonaio/daytona/pkg/views"
	"github.com/daytonaio/daytona/pkg/views/target/selection"
	views_util "github.com/daytonaio/daytona/pkg/views/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var yesFlag bool
var forceFlag bool

var deleteCmd = &cobra.Command{
	Use:     "delete [TARGET]",
	Short:   "Delete a target",
	Aliases: []string{"remove", "rm"},
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.Background()

		var targetDeleteList = []*apiclient.TargetDTO{}
		var targetDeleteListNames = []string{}
		apiClient, err := apiclient_util.GetApiClient(nil)
		if err != nil {
			return err
		}

		workspaceList, res, err := apiClient.WorkspaceAPI.ListWorkspaces(ctx).Execute()
		if err != nil {
			return apiclient_util.HandleErrorResponse(res, err)
		}

		if allFlag {
			return deleteAllTargetsView(ctx, apiClient, workspaceList)
		}

		if len(args) == 0 {
			targetList, res, err := apiClient.TargetAPI.ListTargets(ctx).Execute()
			if err != nil {
				return apiclient_util.HandleErrorResponse(res, err)
			}

			if len(targetList) == 0 {
				views_util.NotifyEmptyTargetList(false)
				return nil
			}

			targetDeleteList = selection.GetTargetsFromPrompt(targetList, "Delete")
			for _, target := range targetDeleteList {
				targetDeleteListNames = append(targetDeleteListNames, target.Name)
			}
		} else {
			for _, arg := range args {
				target, err := apiclient_util.GetTarget(arg, false)
				if err != nil {
					log.Error(fmt.Sprintf("[ %s ] : %v", arg, err))
					continue
				}
				targetDeleteList = append(targetDeleteList, target)
				targetDeleteListNames = append(targetDeleteListNames, target.Name)
			}
		}

		if len(targetDeleteList) == 0 {
			return nil
		}

		var deleteTargetsFlag bool

		if !yesFlag {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title(fmt.Sprintf("Delete target(s): [%s]?", strings.Join(targetDeleteListNames, ", "))).
						Description(fmt.Sprintf("Are you sure you want to delete the target(s): [%s]?", strings.Join(targetDeleteListNames, ", "))).
						Value(&deleteTargetsFlag),
				),
			).WithTheme(views.GetCustomTheme())

			err := form.Run()
			if err != nil {
				return err
			}
		}

		if !yesFlag && !deleteTargetsFlag {
			fmt.Println("Operation canceled.")
		} else {
			for _, target := range targetDeleteList {
				err := removeTarget(ctx, apiClient, target, workspaceList)
				if err != nil {
					log.Error(fmt.Sprintf("[ %s ] : %v", target.Name, err))
				} else {
					views.RenderInfoMessage(fmt.Sprintf("Target '%s' successfully deleted", target.Name))
				}
			}
		}
		return nil
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getTargetNameCompletions()
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Delete all targets")
	deleteCmd.Flags().BoolVarP(&yesFlag, "yes", "y", false, "Confirm deletion without prompt")
	deleteCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Delete a target by force")
}

func deleteAllTargetsView(ctx context.Context, apiClient *apiclient.APIClient, workspaceList []apiclient.WorkspaceDTO) error {
	var deleteAllTargetsFlag bool

	if yesFlag {
		fmt.Println("Deleting all targets.")
		err := deleteAllTargets(ctx, apiClient, workspaceList)
		if err != nil {
			return err
		}
	} else {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Delete all targets?").
					Description("Are you sure you want to delete all targets?").
					Value(&deleteAllTargetsFlag),
			),
		).WithTheme(views.GetCustomTheme())

		err := form.Run()
		if err != nil {
			return err
		}

		if deleteAllTargetsFlag {
			err := deleteAllTargets(ctx, apiClient, workspaceList)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("Operation canceled.")
		}
	}

	return nil
}

func deleteAllTargets(ctx context.Context, apiClient *apiclient.APIClient, workspaceList []apiclient.WorkspaceDTO) error {
	targetList, res, err := apiClient.TargetAPI.ListTargets(ctx).Execute()
	if err != nil {
		return apiclient_util.HandleErrorResponse(res, err)
	}

	for _, target := range targetList {
		err := removeTarget(ctx, apiClient, &target, workspaceList)
		if err != nil {
			log.Errorf("Failed to delete target %s: %v", target.Name, err)
			continue
		}
		views.RenderInfoMessage(fmt.Sprintf("- Target '%s' successfully deleted", target.Name))
	}
	return nil
}

func removeTarget(ctx context.Context, apiClient *apiclient.APIClient, target *apiclient.TargetDTO, workspaceList []apiclient.WorkspaceDTO) error {
	targetWorkspaces := getTargetWorkspacesFromWorkspaceList(target.Id, workspaceList)

	if len(targetWorkspaces) > 0 {
		err := removeWorkspacesForTarget(ctx, apiClient, target.Name, targetWorkspaces)
		if err != nil {
			return err
		}
	}

	message := fmt.Sprintf("Deleting target %s", target.Name)
	err := views_util.WithInlineSpinner(message, func() error {
		res, err := apiClient.TargetAPI.RemoveTarget(ctx, target.Id).Force(forceFlag).Execute()
		if err != nil {
			return apiclient_util.HandleErrorResponse(res, err)
		}

		return nil
	})

	return err
}

func removeWorkspacesForTarget(ctx context.Context, apiClient *apiclient.APIClient, targetName string, targetWorkspaces []apiclient.WorkspaceDTO) error {
	var deleteWorkspacesFlag bool

	var targetWorkspacesNames []string
	for _, workspace := range targetWorkspaces {
		targetWorkspacesNames = append(targetWorkspacesNames, workspace.Name)
	}

	if !yesFlag {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("Target '%s' is in use by %d workspace(s). Delete workspaces: [%s]?", targetName, len(targetWorkspaces), strings.Join(targetWorkspacesNames, ", "))).
					Description(fmt.Sprintf("Do you want to delete workspace(s): [%s]?", strings.Join(targetWorkspacesNames, ", "))).
					Value(&deleteWorkspacesFlag),
			),
		).WithTheme(views.GetCustomTheme())

		err := form.Run()
		if err != nil {
			return err
		}
	}

	if yesFlag || deleteWorkspacesFlag {
		for _, workspace := range targetWorkspaces {
			err := common.RemoveWorkspace(ctx, apiClient, &workspace, forceFlag)
			if err != nil {
				log.Errorf("Failed to delete workspace %s: %v", workspace.Name, err)
				continue
			}
		}
	}

	return nil
}

func getTargetWorkspacesFromWorkspaceList(targetId string, workspaceList []apiclient.WorkspaceDTO) []apiclient.WorkspaceDTO {
	var workspaces []apiclient.WorkspaceDTO
	for _, workspace := range workspaceList {
		if workspace.TargetId == targetId {
			workspaces = append(workspaces, workspace)
		}
	}

	return workspaces
}
