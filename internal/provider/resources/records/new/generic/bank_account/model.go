// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package bankaccount

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordbankaccount "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/bank_account"
)

// BankAccountResourceModel is the New (NSF) bank account resource state model.
type BankAccountResourceModel struct {
	commonrecordbankaccount.BankAccountModel
	new_share.ShareModel
}
