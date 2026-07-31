// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package paymentcard

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordpaymentcard "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/payment_card"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PaymentCardDataSourceModel adds a lookup key (`bank_card`) to the shared payment
// card model. A distinct lookup key name is required because `payment_card` is
// already used by the nested payment_card attribute in PaymentCardModel.
type PaymentCardDataSourceModel struct {
	BankCard types.String `tfsdk:"bank_card"`
	commonrecordpaymentcard.PaymentCardModel
	classic_share.ShareModel
}
