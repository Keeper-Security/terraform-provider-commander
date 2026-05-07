// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	records "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ContactResourceModel maps a Keeper `contact` vault record.
type ContactResourceModel struct {
	records.BaseVaultRecordModel
	Name       *records.NameValue   `tfsdk:"name"`
	Company    types.String         `tfsdk:"company"`
	Email      types.String         `tfsdk:"email"`
	Phone      []records.PhoneValue `tfsdk:"phone"`
	AddressRef types.String         `tfsdk:"address_ref"`
}
