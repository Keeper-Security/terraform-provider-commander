// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamdatabase

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam"
	commonpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_database"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *PamDatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages a Nested-Shared (new) PAM database record with PAM settings and share permissions in your Keeper vault.\n\n" +
			"For more information, see the [PAM Database documentation](https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-database).",
		MarkdownDescription: "Creates and manages a **Nested-Shared (new) PAM database record** with **PAM settings** and **share permissions** in your Keeper vault.\n\n" +
			"For more information, see the [PAM Database documentation](https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-database).",
		Attributes: utils.MergeResourceAttributes(
			commonpamdatabase.SharedAttributes(),
			new_share.ResourceShareAttribute(),
		),
		Blocks: map[string]schema.Block{
			"pam_settings": commonpamrecords.CommonPamSettingsBlock(commonpamrecords.DatabaseProtocols),
		},
	}
}
