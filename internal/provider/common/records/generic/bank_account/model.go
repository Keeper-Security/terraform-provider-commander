// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package bankaccount holds shared model and helpers for the bankAccount record type.
package bankaccount

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BankAccountModel maps a Keeper `bankAccount` vault record.
// Shared between the resource and data source.
type BankAccountModel struct {
	utils.BaseVaultRecordModel

	BankAccount    *utils.BankAccountValue `tfsdk:"bank_account"`
	Name           *utils.NameValue        `tfsdk:"name"`
	Login          types.String            `tfsdk:"login"`
	Password       types.String            `tfsdk:"password"`
	WebsiteAddress types.String            `tfsdk:"website_address"`
	CardRef        types.String            `tfsdk:"card_ref"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
