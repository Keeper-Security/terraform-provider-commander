// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package ssncard holds shared model and helpers for the ssnCard (Identity Card) record type.
package ssncard

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SsnCardModel maps a Keeper `ssnCard` (Identity Card) vault record.
// Shared between the resource and data source.
type SsnCardModel struct {
	utils.BaseVaultRecordModel

	AccountNumber types.String     `tfsdk:"account_number"`
	Name          *utils.NameValue `tfsdk:"name"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
