// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordsshkeys "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/ssh_keys"
)

// SshKeysResourceModel is the classic sshKeys resource state model: shared sshKeys
// fields plus the `share` attribute reconciled via classic_share.
type SshKeysResourceModel struct {
	commonrecordsshkeys.SshKeysModel
	classic_share.ShareModel
}
