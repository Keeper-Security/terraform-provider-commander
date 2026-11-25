// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamcloud

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamcloud "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_cloud"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *PamCloudResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages a Nested-Shared (new) PAM cloud record with share permissions in your Keeper vault.\n\n" +
			"A PAM Cloud record represents a cloud resource, identified by its URL and account identifier.",
		MarkdownDescription: "Creates and manages a **Nested-Shared (new) PAM cloud record** with **share permissions** in your Keeper vault.\n\n" +
			"A PAM Cloud record represents a cloud resource, identified by its **URL** and **account identifier**.",
		Attributes: utils.MergeResourceAttributes(
			commonpamcloud.SharedAttributes(),
			new_share.ResourceShareAttribute(),
		),
	}
}
