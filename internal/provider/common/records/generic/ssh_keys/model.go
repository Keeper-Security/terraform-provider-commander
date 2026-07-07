// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SshKeysModel maps a Keeper `sshKeys` vault record.
// Shared between classic and new resources and data sources.
type SshKeysModel struct {
	utils.BaseVaultRecordModel

	Login      types.String `tfsdk:"login"`
	Passphrase types.String `tfsdk:"passphrase"`
	Hostname   types.String `tfsdk:"hostname"`
	Port       types.String `tfsdk:"port"`
	PublicKey  types.String `tfsdk:"public_key"`
	PrivateKey types.String `tfsdk:"private_key"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
