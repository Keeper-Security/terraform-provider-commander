// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_remote_browser"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *PamRemoteBrowserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a classic PAM remote browser record with per-user share permissions in your Keeper vault.\n\n" +
			"A PAM Remote Browser is a type of KeeperPAM resource that represents a remote browser isolation target, such as a protected internal application or cloud-based web app.\n\n" +
			"For more information, see https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-remote-browser.",
		MarkdownDescription: "Manages a **classic PAM remote browser** record with **per-user share permissions** in your Keeper vault.\n\n" +
			"A PAM Remote Browser is a type of KeeperPAM resource that represents a remote browser isolation target, such as a protected internal application or cloud-based web app.\n\n" +
			"For more information, see [Keeper PAM Remote Browser documentation](https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-remote-browser).",
		Attributes: folderutils.MergeResourceAttributes(
			commonpamremotebrowser.SharedAttributes(),
			classic_share.ResourceShareAttribute(),
		),
	}
}
