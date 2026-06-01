// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
)

// PamDirectoryResourceModel is the Terraform state shared by
// commander_classic_pam_directory and commander_new_pam_directory.
type PamDirectoryResourceModel struct {
	commonpamrecords.CommonPamRecordsResourceModel

	HostnameOrIP   *commonpamrecords.HostnameOrIPModel `tfsdk:"hostname_or_ip"`
	UseSSL         types.Bool                          `tfsdk:"use_ssl"`
	DomainName     types.String                        `tfsdk:"domain_name"`
	AlternativeIPs types.Set                           `tfsdk:"alternative_ips"`
	DirectoryId    types.String                        `tfsdk:"directory_id"`
	DirectoryType  types.String                        `tfsdk:"directory_type"`
	UserMatch      types.String                        `tfsdk:"user_match"`
	ProviderGroup  types.String                        `tfsdk:"provider_group"`
	ProviderRegion types.String                        `tfsdk:"provider_region"`

	PamSettings *commonpamrecords.CommonPamSettingsFieldResourceModel `tfsdk:"pam_settings"`
}
