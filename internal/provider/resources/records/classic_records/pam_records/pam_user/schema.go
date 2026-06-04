// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_user"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *PamUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         commonpamuser.SchemaDescription,
		MarkdownDescription: commonpamuser.SchemaMarkdownDescription,
		Attributes: folderutils.MergeResourceAttributes(
			commonpamuser.SharedAttributes(),
			classic_share.ResourceShareAttribute(),
		),
	}
}
