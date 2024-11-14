// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspace

import (
	"fmt"
	"net/http"
	"time"

	"github.com/daytonaio/daytona/pkg/api/controllers/target/dto"
	"github.com/daytonaio/daytona/pkg/server"
	"github.com/daytonaio/daytona/pkg/workspace"
	"github.com/gin-gonic/gin"
)

// SetWorkspaceMetadata 			godoc
//
//	@Tags			workspace
//	@Summary		Set workspace state
//	@Description	Set workspace state
//	@Param			workspaceId	path	string					true	"Workspace ID"
//	@Param			setState	body	SetWorkspaceMetadata	true	"Set State"
//	@Success		200
//	@Router			/workspace/{workspaceId}/state [post]
//
//	@id				SetWorkspaceMetadata
func SetWorkspaceMetadata(ctx *gin.Context) {
	workspaceId := ctx.Param("workspaceId")

	var setWorkspaceMetadataDTO dto.SetWorkspaceMetadata
	err := ctx.BindJSON(&setWorkspaceMetadataDTO)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	server := server.GetInstance(nil)

	_, err = server.WorkspaceService.SetWorkspaceMetadata(workspaceId, &workspace.WorkspaceMetadata{
		Uptime:    setWorkspaceMetadataDTO.Uptime,
		UpdatedAt: time.Now().Format(time.RFC1123),
		GitStatus: setWorkspaceMetadataDTO.GitStatus,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to set workspace state for %s: %w", workspaceId, err))
		return
	}

	ctx.Status(200)
}
