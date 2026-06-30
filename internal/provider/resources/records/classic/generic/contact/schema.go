// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *ContactResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeResourceAttributes(
			commonrecordsutils.BaseRecordAttributes(),
			map[string]schema.Attribute{
				"name":        commonrecordsutils.NameNestedSchema(false),
				"company":     commonrecordsutils.OptionalStringField("Company", "Company name.", "Company name."),
				"email":       commonrecordsutils.OptionalStringField("Email", "Email address.", "Email address."),
				"phone":       commonrecordsutils.PhoneListSchema(),
				"address_ref": commonrecordsutils.RefUIDField("Linked Address record UID.", "UID of an `address` record linked via `addressRef`."),
			},
			classic_share.ResourceShareAttribute()),
	}
}
