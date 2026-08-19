// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package softwarelicense

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SoftwareLicenseModel maps a Keeper `softwareLicense` vault record.
// Shared between the resource and data source.
type SoftwareLicenseModel struct {
	utils.BaseVaultRecordModel

	SoftwareLicenseKey types.String `tfsdk:"software_license_key"`
	ExpirationDate     types.String `tfsdk:"expiration_date"`
	DateActive         types.String `tfsdk:"date_active"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
