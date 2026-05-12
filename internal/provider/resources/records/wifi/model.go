// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	records "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WifiResourceModel maps a Keeper `wifiCredentials` vault record.
type WifiResourceModel struct {
	records.BaseVaultRecordModel
	SSID         types.String `tfsdk:"ssid"`
	Password     types.String `tfsdk:"password"`
	Encryption   types.String `tfsdk:"encryption"`
	IsSSIDHidden types.Bool   `tfsdk:"is_ssid_hidden"`
}
