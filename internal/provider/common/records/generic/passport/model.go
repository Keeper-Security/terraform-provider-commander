// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package passport holds shared model and helpers for the passport record type.
package passport

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PassportModel maps a Keeper `passport` vault record.
// Shared between the resource and data source.
type PassportModel struct {
	utils.BaseVaultRecordModel

	AccountNumber  types.String     `tfsdk:"account_number"`
	Name           *utils.NameValue `tfsdk:"name"`
	BirthDate      types.String     `tfsdk:"birth_date"`
	AddressRef     types.String     `tfsdk:"address_ref"`
	ExpirationDate types.String     `tfsdk:"expiration_date"`
	DateIssued     types.String     `tfsdk:"date_issued"`
	Password       types.String     `tfsdk:"password"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
