// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package paymentcard

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SharedAttributes returns the Payment Card resource attribute map shared between
// classic and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"payment_card":    commonrecordsutils.PaymentCardNestedSchema(),
			"cardholder_name": commonrecordsutils.OptionalCardholderNameField(),
			"pin_code":        commonrecordsutils.OptionalPinCodeField(),
			"address_ref":     commonrecordsutils.OptionalAddressRefField(),
			"custom":          commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed Payment Card data source attributes
// shared between classic and new data sources. Callers add the lookup key
// (e.g. bank_card) and share-extension attributes separately.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"payment_card":    commonrecordsutils.PaymentCardDataSourceNestedSchema(),
			"cardholder_name": commonrecordsutils.ComputedCardholderNameField(),
			"pin_code":        commonrecordsutils.ComputedPinCodeField(),
			"address_ref":     commonrecordsutils.ComputedAddressRefField(),
			"custom":          commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
