// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package bankaccount

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordbankaccount "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/bank_account"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BankAccountDataSourceModel adds a lookup key (`account`) to the
// shared bank account model.
type BankAccountDataSourceModel struct {
	Account types.String `tfsdk:"account"`
	commonrecordbankaccount.BankAccountModel
	new_share.ShareModel
}
