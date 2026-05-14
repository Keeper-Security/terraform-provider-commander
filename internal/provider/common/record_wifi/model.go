// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package recordwifi

import (
	records "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WifiModel maps a Keeper `wifiCredentials` vault record.
// Shared between the resource and data source.
type WifiModel struct {
	records.BaseVaultRecordModel
	SSID         types.String `tfsdk:"ssid"`
	Password     types.String `tfsdk:"password"`
	Encryption   types.String `tfsdk:"encryption"`
	IsSSIDHidden types.Bool   `tfsdk:"is_ssid_hidden"`
}
