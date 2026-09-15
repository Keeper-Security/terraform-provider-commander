// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package softwarelicense

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordsoftwarelicense "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/software_license"
)

// SoftwareLicenseResourceModel is the classic softwareLicense resource state model.
type SoftwareLicenseResourceModel struct {
	commonrecordsoftwarelicense.SoftwareLicenseModel
	new_share.ShareModel
}
