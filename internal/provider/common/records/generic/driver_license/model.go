// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package driverlicense holds shared model and helpers for the driverLicense record type.
package driverlicense

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DriverLicenseModel maps a Keeper `driverLicense` vault record.
// Shared between the resource and data source.
type DriverLicenseModel struct {
	utils.BaseVaultRecordModel

	AccountNumber  types.String     `tfsdk:"account_number"`
	Name           *utils.NameValue `tfsdk:"name"`
	BirthDate      types.String     `tfsdk:"birth_date"`
	AddressRef     types.String     `tfsdk:"address_ref"`
	ExpirationDate types.String     `tfsdk:"expiration_date"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
