// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_records"
)

type PamDatabaseResourceModel struct {
	commonpamrecords.CommonPamRecordsResourceModel

	HostnameOrIP   *commonpamrecords.HostnameOrIPModel `tfsdk:"hostname_or_ip"`
	UseSSL         types.Bool                          `tfsdk:"use_ssl"`
	DatabaseId     types.String                        `tfsdk:"database_id"`
	DatabaseType   types.String                        `tfsdk:"database_type"`
	ProviderGroup  types.String                        `tfsdk:"provider_group"`
	ProviderRegion types.String                        `tfsdk:"provider_region"`

	PamSettings *commonpamrecords.CommonPamSettingsFieldResourceModel `tfsdk:"pam_settings"`
}
