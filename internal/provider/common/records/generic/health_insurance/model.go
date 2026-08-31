// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package healthinsurance holds shared model and helpers for the healthInsurance record type.
package healthinsurance

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// HealthInsuranceModel maps a Keeper `healthInsurance` vault record.
// Shared between the resource and data source.
type HealthInsuranceModel struct {
	utils.BaseVaultRecordModel

	AccountNumber  types.String     `tfsdk:"account_number"`
	Name           *utils.NameValue `tfsdk:"name"`
	Login          types.String     `tfsdk:"login"`
	Password       types.String     `tfsdk:"password"`
	WebsiteAddress types.String     `tfsdk:"website_address"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
