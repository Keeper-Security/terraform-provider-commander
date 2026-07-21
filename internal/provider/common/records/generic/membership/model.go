// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package membership holds shared model and helpers for the membership record type.
package membership

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MembershipModel maps a Keeper `membership` vault record.
// Shared between the resource and data source.
type MembershipModel struct {
	utils.BaseVaultRecordModel

	AccountNumber types.String     `tfsdk:"account_number"`
	Name          *utils.NameValue `tfsdk:"name"`
	Password      types.String     `tfsdk:"password"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
