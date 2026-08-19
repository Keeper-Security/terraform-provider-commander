// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DatabaseModel maps a Keeper `databaseCredentials` vault record.
// Shared between the resource and data source.
type DatabaseModel struct {
	utils.BaseVaultRecordModel

	Login    types.String `tfsdk:"login"`
	Hostname types.String `tfsdk:"hostname"`
	Port     types.String `tfsdk:"port"`
	Type     types.String `tfsdk:"type"`
	Password types.String `tfsdk:"password"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
