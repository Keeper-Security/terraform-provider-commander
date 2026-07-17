// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package paymentcard holds shared model and helpers for the payment card (bankCard) record type.
package paymentcard

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PaymentCardModel maps a Keeper `bankCard` vault record.
// Shared between the resource and data source.
type PaymentCardModel struct {
	utils.BaseVaultRecordModel

	PaymentCard    *utils.PaymentCardValue `tfsdk:"payment_card"`
	CardholderName types.String            `tfsdk:"cardholder_name"`
	PinCode        types.String            `tfsdk:"pin_code"`
	AddressRef     types.String            `tfsdk:"address_ref"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
