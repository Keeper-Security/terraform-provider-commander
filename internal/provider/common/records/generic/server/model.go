// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ServerModel maps a Keeper `serverCredentials` vault record.
// Shared between the resource and data source.
type ServerModel struct {
	utils.BaseVaultRecordModel

	Login    types.String `tfsdk:"login"`
	Password types.String `tfsdk:"password"`
	Hostname types.String `tfsdk:"hostname"`
	Port     types.String `tfsdk:"port"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
