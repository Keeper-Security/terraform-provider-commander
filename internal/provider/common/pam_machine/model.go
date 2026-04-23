// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package pammachine holds shared Terraform schema fragments for PAM machine resources.
package pammachine

import (
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_records"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// HostnameOrIPModel maps the nested hostname_or_ip object.
type HostnameOrIPModel struct {
	HostName types.String `tfsdk:"hostname"`
	AdministrativePort types.Int32 `tfsdk:"administrative_port"`
}

// PamSettingsModel maps the nested pam_settings object.
// TODO: define fields when implementing PAM settings.
// type PamSettingsModel struct {
// 	commonpamrecords.CommonPamSettingsFieldResourceModel
// }

// PamMachineResourceModel is the Terraform state for commander_pam_machine.
type PamMachineResourceModel struct {
	commonpamrecords.CommonPamRecordsResourceModel

	HostnameOrIP    *HostnameOrIPModel `tfsdk:"hostname_or_ip"`
	OperatingSystem types.String       `tfsdk:"operating_system"`
	InstanceName    types.String       `tfsdk:"instance_name"`
	InstanceId      types.String       `tfsdk:"instance_id"`
	ProviderGroup   types.String       `tfsdk:"provider_group"`
	ProviderRegion  types.String       `tfsdk:"provider_region"`

	PamSettings *commonpamrecords.CommonPamSettingsFieldResourceModel `tfsdk:"pam_settings"`
}
