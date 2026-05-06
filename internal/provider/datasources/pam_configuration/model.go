// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	commonpamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_configuration"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PamConfigurationDataSourceModel struct {
	PamConfiguration types.String `tfsdk:"pam_configuration"`
	commonpamconfiguration.PamConfigurationResourceModel
}
