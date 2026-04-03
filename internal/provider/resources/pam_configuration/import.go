// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	importID := strings.TrimSpace(req.ID)
	if importID == "" {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_INVALID_IMPORT_ID,
			"Import ID cannot be empty. Use the PAM configuration UID when defined.",
		)
		return
	}

	state := PamConfigurationResourceModel{
		Id:                            types.StringValue(importID),
		Environment:                   types.StringNull(),
		Title:                         types.StringNull(),
		Gateway:                       types.StringNull(),
		ApplicationFolder:             types.StringNull(),
		Schedule:                      types.StringNull(),
		PortMapping:                   types.SetNull(types.StringType),
		Connections:                   types.BoolNull(),
		Tunneling:                     types.BoolNull(),
		Rotation:                      types.BoolNull(),
		RemoteBrowserIsolation:        types.BoolNull(),
		ConnectionsRecording:          types.BoolNull(),
		TypescriptRecording:           types.BoolNull(),
		AIThreatDetection:             types.BoolNull(),
		AITerminateSessionOnDetection: types.BoolNull(),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
