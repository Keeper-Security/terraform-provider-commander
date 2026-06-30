// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamuser

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_user"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := strings.TrimSpace(req.ID)
	if importID == "" {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_INVALID_IMPORT_ID,
			"Import ID cannot be empty. Use the new PAM user record UID when defined.",
		)
		return
	}

	state := PamUserResourceModel{
		PamUserSharedModel: commonpamuser.PamUserSharedModel{
			BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
				Id:             types.StringValue(importID),
				Title:          types.StringNull(),
				Notes:          types.StringNull(),
				FolderLocation: types.StringNull(),
			},
			Login:                types.StringNull(),
			Password:             types.StringNull(),
			DistinguishedName:    types.StringNull(),
			PrivatePEMKey:        types.StringNull(),
			PublicKey:            types.StringNull(),
			PrivateKeyPassphrase: types.StringNull(),
			ConnectDatabase:      types.StringNull(),
			Managed:              types.BoolNull(),
			RotationSettings:     nil,
		},
		ShareModel: new_share.ShareModel{
			Share: types.MapNull(new_share.ShareEntryAttrType),
		},
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
