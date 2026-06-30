// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpammachine

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_machine"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PamMachineDataSourceModel struct {
	PamMachine types.String `tfsdk:"pam_machine"`
	commonpammachine.PamMachineResourceModel
	new_share.ShareModel
}
