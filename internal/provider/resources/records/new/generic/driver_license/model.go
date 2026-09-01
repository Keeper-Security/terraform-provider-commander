// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package driverlicense

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecorddriverlicense "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/driver_license"
)

// DriverLicenseResourceModel is the new driver's license resource state model:
// shared driver's license fields plus the `share` attribute reconciled via new_share.
type DriverLicenseResourceModel struct {
	commonrecorddriverlicense.DriverLicenseModel
	new_share.ShareModel
}
