// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package recordaddress holds shared model and helpers for the address record type.
package recordaddress

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	records "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records"
)

// AddressModel maps a Keeper `address` vault record.
// Shared between the resource and data source.
type AddressModel struct {
	records.BaseVaultRecordModel
	Street1 types.String `tfsdk:"street1"`
	Street2 types.String `tfsdk:"street2"`
	City    types.String `tfsdk:"city"`
	State   types.String `tfsdk:"state"`
	Zip     types.String `tfsdk:"zip"`
	Country types.String `tfsdk:"country"`
}
