// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	"context"
	"strings"

	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_remote_browser"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamRemoteBrowserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := strings.TrimSpace(req.ID)
	if importID == "" {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_INVALID_IMPORT_ID,
			"Import ID cannot be empty. Use the PAM remote browser record UID when defined.",
		)
		return
	}

	state := commonpamremotebrowser.PamRemoteBrowserResourceModel{
		Id:                       types.StringValue(importID),
		Title:                    types.StringNull(),
		Url:                      types.StringNull(),
		Notes:                    types.StringNull(),
		Folder:                   types.StringNull(),
		PamRemoteBrowserSettings: nil,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
