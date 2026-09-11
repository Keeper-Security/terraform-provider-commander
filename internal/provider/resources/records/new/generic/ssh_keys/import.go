// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordsshkeys "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/ssh_keys"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *SshKeysResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := strings.TrimSpace(req.ID)
	if id == "" {
		resp.Diagnostics.AddError(utils.ERR_MSG_INVALID_IMPORT_ID, "Import ID cannot be empty. Use the SSH keys record UID.")
		return
	}
	m := SshKeysResourceModel{
		SshKeysModel: commonrecordsshkeys.SshKeysModel{
			BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
				Id:             types.StringValue(id),
				Title:          types.StringNull(),
				Notes:          types.StringNull(),
				FolderLocation: types.StringNull(),
			},
			Login:      types.StringNull(),
			Passphrase: types.StringNull(),
			Hostname:   types.StringNull(),
			Port:       types.StringNull(),
			PublicKey:  types.StringNull(),
			PrivateKey: types.StringNull(),
			Custom:     nil,
		},
		ShareModel: new_share.ShareModel{
			Share: types.MapNull(new_share.ShareEntryAttrType),
		},
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &m)...)
}
