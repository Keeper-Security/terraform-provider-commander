// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam"
	commonpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_database"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *PamDatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages a classic PAM database record with PAM settings and per-user share permissions in your Keeper vault.\n\n" +
			"A PAM Database record is a type of KeeperPAM resource that represents a database server, either on-prem or hosted in the cloud.\n\n" +
			"For more information, see the https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-database.",
		MarkdownDescription: "Creates and manages a **classic PAM database record** with **PAM settings** and **per-user share permissions** in your Keeper vault.\n\n" +
			"A PAM Database record is a type of KeeperPAM resource that represents a database server, either on-prem or hosted in the cloud.\n\n" +
			"For more information, see the [PAM Database documentation](https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-database).",
		Attributes: utils.MergeResourceAttributes(
			commonpamdatabase.SharedAttributes(),
			classic_share.ResourceShareAttribute(),
		),
		Blocks: map[string]schema.Block{
			"pam_settings": commonpamrecords.CommonPamSettingsBlock(commonpamrecords.DatabaseProtocols),
		},
	}
}
