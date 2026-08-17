// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordsaasconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/saas_configuration"
)

// SaasConfigurationResourceModel is the classic SaaS configuration resource state model.
type SaasConfigurationResourceModel struct {
	commonrecordsaasconfiguration.SaasConfigurationModel
	new_share.ShareModel
}
