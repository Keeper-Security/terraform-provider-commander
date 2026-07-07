// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordsshkeys "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/ssh_keys"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SshKeysDataSourceModel adds a lookup key (`ssh_keys`) to the shared sshKeys model.
// The lookup attribute is named ssh_keys because login is a field on the record.
type SshKeysDataSourceModel struct {
	SshKeysRecord types.String `tfsdk:"ssh_keys"`
	commonrecordsshkeys.SshKeysModel
	classic_share.ShareModel
}
