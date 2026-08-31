// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package paymentcard

// Commander CLI record field keys for payment card (Keeper record type `bankCard`).
const (
	FlagPaymentCard    = "f.paymentCard"
	FlagCardholderName = "f.text.cardholderName"
	FlagPinCode        = "f.pinCode"
	FlagAddressRef     = "f.addressRef"

	// CardholderNameFieldLabel is the Keeper `text` field label used for cardholderName.
	CardholderNameFieldLabel = "cardholderName"
)
