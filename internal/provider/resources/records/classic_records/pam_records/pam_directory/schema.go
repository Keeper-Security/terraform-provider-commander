// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	"context"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_directory"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *PamDirectoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages PAM directory record with pam settings in your Keeper vault.\n\n" +
			"A PAM Directory record is a type of KeeperPAM resource that represents an Active Directory or OpenLDAP service, either on-prem or hosted in the cloud..\n\n" +
			"For more information, see the https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-directory.",
		MarkdownDescription: "Creates and manages **PAM directory record with pam settings** in your Keeper vault.\n\n" +
			"A PAM Directory record is a type of KeeperPAM resource that represents an Active Directory or OpenLDAP service, either on-prem or hosted in the cloud.\n\n" +
			"For more information, see the [PAM Directory documentation](https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-directory).",
		Attributes: commonpamdirectory.SharedAttributes(),
		Blocks: map[string]schema.Block{
			"pam_settings": commonpamrecords.CommonPamSettingsBlock(),
		},
	}
}
