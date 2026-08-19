// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package driverlicense

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecorddriverlicense "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/driver_license"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DriverLicenseDataSourceModel adds a lookup key (`driver_license`) to the
// shared driver's license model.
type DriverLicenseDataSourceModel struct {
	DriverLicense types.String `tfsdk:"driver_license"`
	commonrecorddriverlicense.DriverLicenseModel
	classic_share.ShareModel
}
