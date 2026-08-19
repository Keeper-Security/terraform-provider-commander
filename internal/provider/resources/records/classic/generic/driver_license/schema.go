// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package driverlicense

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecorddriverlicense "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/driver_license"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *DriverLicenseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeResourceAttributes(
			commonrecorddriverlicense.SharedAttributes(),
			classic_share.ResourceShareAttribute(),
		),
	}
}
