// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase

import (
	"context"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_database"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *PamDatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages PAM database record with pam settings in your Keeper vault.\n\n" +
			"A PAM Database record is a type of KeeperPAM resource that represents a database server, either on-prem or hosted in the cloud.\n\n" +
			"For more information, see the https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-database.",
		MarkdownDescription: "Creates and manages **PAM database record with pam settings** in your Keeper vault.\n\n" +
			"A PAM Database record is a type of KeeperPAM resource that represents a database server, either on-prem or hosted in the cloud.\n\n" +
			"For more information, see the [PAM Database documentation](https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-database).",
		Attributes: commonpamdatabase.SharedAttributes(),
		Blocks: map[string]schema.Block{
			"pam_settings": commonpamrecords.CommonPamSettingsBlock(),
		},
	}
}
