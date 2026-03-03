// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharedfolder

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports import by shared folder UID.
// Import ID: the shared folder UID (e.g. "BTbjhOmqw9iYal3OQJ9UAQ").
// After import, Terraform runs Read to refresh state from the API.
func (r *SharedFolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(ErrSummaryProviderConfig, err.Error())
		return
	}

	importID := strings.TrimSpace(req.ID)
	if importID == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID cannot be empty. Use the shared folder UID (e.g. \"BTbjhOmqw9iYal3OQJ9UAQ\").",
		)
		return
	}

	state := SharedFolderResourceModel{
		Id:                types.StringValue(importID),
		Name:              types.StringNull(),
		FolderLocation:    types.StringNull(),
		UserPermissions:   nil,
		RecordPermissions: nil,
		Records:           types.MapNull(RecordEntryMapElemType),
		Users:             types.MapNull(UserEntryMapElemType),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
