// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package ssncard

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordssncard "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/ssn_card"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *SsnCardResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeResourceAttributes(
			commonrecordssncard.SharedAttributes(),
			new_share.ResourceShareAttribute(),
		),
	}
}
