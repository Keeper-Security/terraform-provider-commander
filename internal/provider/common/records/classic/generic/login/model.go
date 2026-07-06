// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package login

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LoginModel maps a Keeper `login` vault record.
// Shared between the resource and data source.
type LoginModel struct {
	utils.BaseVaultRecordModel

	Login          types.String `tfsdk:"login"`
	Password       types.String `tfsdk:"password"`
	WebsiteAddress types.String `tfsdk:"website_address"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`

	classic_share.ShareModel
}
