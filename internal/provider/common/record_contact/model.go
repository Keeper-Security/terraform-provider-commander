// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package recordcontact holds shared model and helpers for the contact record type.
package recordcontact

import (
	records "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ContactModel maps a Keeper `contact` vault record.
// Shared between the resource and data source.
type ContactModel struct {
	records.BaseVaultRecordModel
	Name       *records.NameValue   `tfsdk:"name"`
	Company    types.String         `tfsdk:"company"`
	Email      types.String         `tfsdk:"email"`
	Phone      []records.PhoneValue `tfsdk:"phone"`
	AddressRef types.String         `tfsdk:"address_ref"`
}
