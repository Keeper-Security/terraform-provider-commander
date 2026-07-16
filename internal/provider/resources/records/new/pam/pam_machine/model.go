// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpammachine

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_machine"
)

// PamMachineResourceModel is the new (nested-shared) PAM Machine resource
// state model: the shared PAM Machine fields plus the `share` attribute.
type PamMachineResourceModel struct {
	commonpammachine.PamMachineResourceModel
	new_share.ShareModel
}
