// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package list

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/daytonaio/daytona/internal/util"
	"github.com/daytonaio/daytona/pkg/apiclient"
	"github.com/daytonaio/daytona/pkg/views"
	info_view "github.com/daytonaio/daytona/pkg/views/target/info"
	views_util "github.com/daytonaio/daytona/pkg/views/util"
)

type RowData struct {
	Name           string
	Provider       string
	WorkspaceCount string
	Default        bool
	Status         apiclient.ModelsResourceStateName
	Options        string
	Uptime         string
}

func ListTargets(targetList []apiclient.TargetDTO, verbose bool, activeProfileName string) {
	if len(targetList) == 0 {
		views_util.NotifyEmptyTargetList(true)
		return
	}

	SortTargets(&targetList)

	headers := []string{"Target", "Options", "# Workspaces", "Default", "Status"}

	data := util.ArrayMap(targetList, func(target apiclient.TargetDTO) []string {
		provider := target.ProviderInfo.Name
		if target.ProviderInfo.Label != nil {
			provider = *target.ProviderInfo.Label
		}

		rowData := RowData{
			Name:           target.Name,
			Provider:       provider,
			Options:        target.Options,
			WorkspaceCount: fmt.Sprint(len(target.Workspaces)),
			Default:        target.Default,
			Status:         target.State.Name,
		}

		if target.Metadata != nil && target.Metadata.Uptime > 0 {
			rowData.Uptime = util.FormatUptime(target.Metadata.Uptime)
		}

		return getRowFromRowData(rowData)
	})

	footer := lipgloss.NewStyle().Foreground(views.LightGray).Render(views.GetListFooter(activeProfileName, &views.Padding{}))

	table := views_util.GetTableView(data, headers, &footer, func() {
		renderUnstyledList(targetList)
	})

	fmt.Println(table)
}

func renderUnstyledList(targetList []apiclient.TargetDTO) {
	for _, target := range targetList {
		info_view.Render(&target, true)

		if target.Id != targetList[len(targetList)-1].Id {
			fmt.Printf("\n%s\n\n", views.SeparatorString)
		}
	}
}

func getRowFromRowData(rowData RowData) []string {
	stateLabel := views.GetStateLabel(rowData.Status)

	if rowData.Uptime != "" {
		stateLabel = fmt.Sprintf("%s (%s)", stateLabel, rowData.Uptime)
	}

	var isDefault string

	if rowData.Default {
		isDefault = "Yes"
	} else {
		isDefault = "/"
	}

	return []string{
		fmt.Sprintf("%s %s", views.NameStyle.Render(rowData.Name), views.DefaultRowDataStyle.Render(fmt.Sprintf("(%s)", rowData.Provider))),
		views.DefaultRowDataStyle.Render(rowData.Options),
		views.DefaultRowDataStyle.Render(rowData.WorkspaceCount),
		isDefault,
		stateLabel,
	}
}

func SortTargets(targetList *[]apiclient.TargetDTO) {
	sort.Slice(*targetList, func(i, j int) bool {
		return (*targetList)[i].Default && !(*targetList)[j].Default
	})
}
