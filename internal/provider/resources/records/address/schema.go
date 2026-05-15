// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package address

import (
	"context"

	records "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *AddressResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	attrs := records.BaseRecordAttributes()
	attrs["street1"] = records.OptionalStringField("Street Address Line 1", "Street address line 1.", "Street address line 1.")
	attrs["street2"] = records.OptionalStringField("Street Address Line 2", "Street address line 2.", "Street address line 2.")
	attrs["city"] = records.OptionalStringField("City", "City.", "City.")
	attrs["state"] = records.OptionalStringField("State", "State or province.", "State or province.")
	attrs["zip"] = records.OptionalStringField("Zip", "Zip or postal code.", "Zip or postal code.")
	attrs["country"] = records.OptionalStringField("Country", "ISO 3166-1 alpha-2 country code.", "ISO 3166-1 alpha-2 country code.")
	attrs["share"] = records.ShareAttribute()

	resp.Schema = schema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes:          attrs,
	}
}
