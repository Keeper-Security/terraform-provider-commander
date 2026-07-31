// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package bankaccount

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordbankaccount "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/bank_account"
)

// BankAccountResourceModel is the classic bank account resource state model: shared
// bank account fields plus the `share` attribute reconciled via classic_share.
type BankAccountResourceModel struct {
	commonrecordbankaccount.BankAccountModel
	classic_share.ShareModel
}
