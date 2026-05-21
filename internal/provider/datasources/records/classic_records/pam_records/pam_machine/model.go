// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records/pam_machine"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PamMachineDataSourceModel struct {
	PamMachine types.String `tfsdk:"pam_machine"`
	commonpammachine.PamMachineResourceModel
}
