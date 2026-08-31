// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
)

// SaasConfigurationModel maps a Keeper `saasConfiguration` vault record.
// Shared between the resource and data source.
type SaasConfigurationModel struct {
	commonrecordsutils.BaseVaultRecordModel
	Custom []commonrecordsutils.CustomFieldModel `tfsdk:"custom"`
}
