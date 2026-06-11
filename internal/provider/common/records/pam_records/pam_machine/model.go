// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package pammachine holds shared Terraform schema fragments for PAM machine resources.
package pammachine

import (
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PamMachineResourceModel is the Terraform state shared by
// commander_classic_pam_machine and commander_new_pam_machine.
type PamMachineResourceModel struct {
	commonrecordsutils.BaseVaultRecordModel

	HostnameOrIP    *commonpamrecords.HostnameOrIPModel `tfsdk:"hostname_or_ip"`
	OperatingSystem types.String                        `tfsdk:"operating_system"`
	InstanceName    types.String                        `tfsdk:"instance_name"`
	InstanceId      types.String                        `tfsdk:"instance_id"`
	ProviderGroup   types.String                        `tfsdk:"provider_group"`
	ProviderRegion  types.String                        `tfsdk:"provider_region"`

	PamSettings *commonpamrecords.CommonPamSettingsFieldResourceModel `tfsdk:"pam_settings"`
}
