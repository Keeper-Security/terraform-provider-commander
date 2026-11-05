// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SaasConfigurationModel maps a Keeper `saasConfiguration` vault record.
// Shared between the resource and data source.
type SaasConfigurationModel struct {
	commonrecordsutils.BaseVaultRecordModel
	Configuration types.String                          `tfsdk:"configuration"`
	Gateway       types.String                          `tfsdk:"gateway"`
	Custom        []commonrecordsutils.CustomFieldModel `tfsdk:"custom"`
}
