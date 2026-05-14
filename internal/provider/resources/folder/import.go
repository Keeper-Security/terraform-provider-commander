// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package folder

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports terraform import by folder UID.
// Import ID: the folder UID (e.g. "NJiANrRnbuvVEOgnqjiYaw").
// After import, Terraform runs Read to refresh state from the API.
func (r *FolderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	importID := strings.TrimSpace(req.ID)
	if importID == "" {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_INVALID_IMPORT_ID,
			"Import ID cannot be empty. Use the folder UID.",
		)
		return
	}

	state := FolderResourceModel{
		Id:             types.StringValue(importID),
		Name:           types.StringNull(),
		FolderLocation: types.StringNull(),
		Color:          types.StringNull(),
		Records:        types.SetNull(types.StringType),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
