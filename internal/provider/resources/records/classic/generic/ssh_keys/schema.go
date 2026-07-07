// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordsshkeys "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/ssh_keys"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *SshKeysResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeResourceAttributes(
			commonrecordsshkeys.SharedAttributes(),
			classic_share.ResourceShareAttribute(),
		),
	}
}
