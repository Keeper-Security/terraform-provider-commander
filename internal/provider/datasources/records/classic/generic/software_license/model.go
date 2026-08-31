// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package softwarelicense

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordsoftwarelicense "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/software_license"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SoftwareLicenseDataSourceModel maps a Keeper `softwareLicense` vault record for read-only access.
type SoftwareLicenseDataSourceModel struct {
	SoftwareLicense types.String `tfsdk:"software_license"`
	commonrecordsoftwarelicense.SoftwareLicenseModel
	classic_share.ShareModel
}
