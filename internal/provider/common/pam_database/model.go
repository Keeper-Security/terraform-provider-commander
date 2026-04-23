// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase

import (
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_records"
)

/*
	NOTE: CURRENTLY KeeperDbProxy ( --keeper-db-proxy ) is not supported in the ADMIN CONSOLE, So not implementing it for now.
*/
// type PamDatabaseTunnelResourceModel struct {
// 	commonpamrecords.CommonPamSettingsTunnelResourceModel
// 	KeeperDbProxy types.Bool `tfsdk:"keeper_db_proxy"`
// }

// type PamDatabaseSettingsResourceModel struct {
// 	AllowSupplyHost types.Bool                                                 `tfsdk:"allow_supply_host"`
// 	Connection      *commonpamrecords.CommonPamSettingsConnectionResourceModel `tfsdk:"connection"`
// 	Tunnel          *PamDatabaseTunnelResourceModel                            `tfsdk:"tunnel"`
// 	Configuration   types.String                                               `tfsdk:"configuration"`
// }

type PamDatabaseResourceModel struct {
	commonpamrecords.CommonPamRecordsResourceModel

	PamSettings *commonpamrecords.CommonPamSettingsFieldResourceModel `tfsdk:"pam_settings"`
}
