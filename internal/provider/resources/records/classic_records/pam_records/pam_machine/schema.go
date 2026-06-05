// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_machine"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *PamMachineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages a classic PAM machine record with PAM settings and per-user share permissions in your Keeper vault.\n\n" +
			"A PAM Machine record is a type of KeeperPAM resource that represents a workload, such as a Windows or Linux server.\n\n" +
			"For more information, see the [PAM Machine documentation](https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-machine).",
		MarkdownDescription: "Creates and manages a **classic PAM machine record** with **PAM settings** and **per-user share permissions** in your Keeper vault.\n\n" +
			"A PAM Machine record is a type of KeeperPAM resource that represents a workload, such as a Windows or Linux server.\n\n" +
			"For more information, see the [PAM Machine documentation](https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-machine).",
		Attributes: utils.MergeResourceAttributes(
			commonpammachine.SharedAttributes(),
			classic_share.ResourceShareAttribute(),
		),
		Blocks: map[string]schema.Block{
			"pam_settings": commonpamrecords.CommonPamSettingsBlock(),
		},
	}
}
