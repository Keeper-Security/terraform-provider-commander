// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	commonrecordwifi "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/record_wifi"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WifiDataSourceModel maps a Keeper `wifiCredentials` vault record for read-only access.
type WifiDataSourceModel struct {
	Wifi types.String `tfsdk:"wifi"`
	commonrecordwifi.WifiModel
}
