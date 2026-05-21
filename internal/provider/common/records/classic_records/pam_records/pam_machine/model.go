// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package pammachine holds shared Terraform schema fragments for PAM machine resources.
package pammachine

import (
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PamMachineResourceModel is the Terraform state for commander_pam_machine.
type PamMachineResourceModel struct {
	commonpamrecords.CommonPamRecordsResourceModel

	HostnameOrIP    *commonpamrecords.HostnameOrIPModel `tfsdk:"hostname_or_ip"`
	OperatingSystem types.String                        `tfsdk:"operating_system"`
	InstanceName    types.String                        `tfsdk:"instance_name"`
	InstanceId      types.String                        `tfsdk:"instance_id"`
	ProviderGroup   types.String                        `tfsdk:"provider_group"`
	ProviderRegion  types.String                        `tfsdk:"provider_region"`

	PamSettings *commonpamrecords.CommonPamSettingsFieldResourceModel `tfsdk:"pam_settings"`
}
