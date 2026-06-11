// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_machine"
)

// PamMachineResourceModel is the classic PAM Machine resource state model:
// the shared PAM Machine fields plus the `share` attribute reconciled via
// the classic_share package and the `share-record` Commander CLI.
type PamMachineResourceModel struct {
	commonpammachine.PamMachineResourceModel
	classic_share.ShareModel
}
