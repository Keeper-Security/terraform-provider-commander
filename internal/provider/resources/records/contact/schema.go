// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	"context"

	records "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *ContactResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	attrs := records.BaseRecordAttributes()
	attrs["name"] = records.NameNestedSchema(false)
	attrs["company"] = records.OptionalStringField("Company", "Company name.", "Company name (`f.text.company`).")
	attrs["email"] = records.OptionalStringField("Email", "Email address.", "Email address.")
	attrs["phone"] = records.PhoneListSchema()
	attrs["address_ref"] = records.RefUIDField("Linked Address record UID.", "UID of an `address` record linked via `addressRef`.")

	resp.Schema = schema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes:          attrs,
		Blocks:              records.BaseRecordBlocks(),
	}
}
