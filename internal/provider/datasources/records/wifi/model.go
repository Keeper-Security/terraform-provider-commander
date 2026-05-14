// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	records "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WifiDataSourceModel maps a Keeper `wifiCredentials` vault record for read-only access.
type WifiDataSourceModel struct {
	RecordUID    types.String               `tfsdk:"record_uid"`
	Id           types.String               `tfsdk:"id"`
	Title        types.String               `tfsdk:"title"`
	Folder       types.String               `tfsdk:"folder"`
	Notes        types.String               `tfsdk:"notes"`
	SSID         types.String               `tfsdk:"ssid"`
	Password     types.String               `tfsdk:"password"`
	Encryption   types.String               `tfsdk:"encryption"`
	IsSSIDHidden types.Bool                 `tfsdk:"is_ssid_hidden"`
	Custom       []records.CustomFieldModel `tfsdk:"custom"`
	Share        types.Map                  `tfsdk:"share"`
}
