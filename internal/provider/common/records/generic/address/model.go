// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package address holds shared model and helpers for the address record type.
package address

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
)

// AddressModel maps a Keeper `address` vault record.
// Shared between the resource and data source.
type AddressModel struct {
	utils.BaseVaultRecordModel

	Address *utils.AddressValue `tfsdk:"address"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
