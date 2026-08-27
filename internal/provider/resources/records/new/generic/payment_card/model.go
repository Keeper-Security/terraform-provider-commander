// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package paymentcard

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordpaymentcard "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/payment_card"
)

// PaymentCardResourceModel is the classic payment card resource state model: shared
// payment card fields plus the `share` attribute reconciled via classic_share.
type PaymentCardResourceModel struct {
	commonrecordpaymentcard.PaymentCardModel
	new_share.ShareModel
}
