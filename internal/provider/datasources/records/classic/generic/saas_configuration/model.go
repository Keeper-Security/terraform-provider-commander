// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordsaasconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/saas_configuration"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SaasConfigurationDataSourceModel maps a Keeper `saasConfiguration` vault record for read-only access.
type SaasConfigurationDataSourceModel struct {
	SaasConfiguration types.String `tfsdk:"saas_configuration"`
	commonrecordsaasconfiguration.SaasConfigurationModel
	classic_share.ShareModel
}
