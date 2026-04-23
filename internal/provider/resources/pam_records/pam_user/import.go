// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := strings.TrimSpace(req.ID)
	if importID == "" {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_INVALID_IMPORT_ID,
			"Import ID cannot be empty. Use the PAM User record UID.",
		)
		return
	}

	state := PamUserResourceModel{
		Id:                types.StringValue(importID),
		Title:             types.StringNull(),
		Login:             types.StringNull(),
		Password:          types.StringNull(),
		Folder:            types.StringNull(),
		Notes:             types.StringNull(),
		DistinguishedName: types.StringNull(),
		PrivatePEMKey:     types.StringNull(),
		ConnectDatabase:   types.StringNull(),
		Managed:           types.BoolNull(),
		RotationSettings:  nil,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
