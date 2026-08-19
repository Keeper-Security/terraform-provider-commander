// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package contact holds shared model and helpers for the contact record type.
package contact

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ContactModel maps a Keeper `contact` vault record.
// Shared between the resource and data source.
type ContactModel struct {
	utils.BaseVaultRecordModel

	Name       *utils.NameValue   `tfsdk:"name"`
	Company    types.String       `tfsdk:"company"`
	Email      types.String       `tfsdk:"email"`
	Phone      []utils.PhoneValue `tfsdk:"phone"`
	AddressRef types.String       `tfsdk:"address_ref"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
