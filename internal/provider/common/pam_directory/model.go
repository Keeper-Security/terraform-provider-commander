// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_records"
)

type PamDirectoryResourceModel struct {
	commonpamrecords.CommonPamRecordsResourceModel

	PamSettings *commonpamrecords.CommonPamSettingsFieldResourceModel `tfsdk:"pam_settings"`
}
