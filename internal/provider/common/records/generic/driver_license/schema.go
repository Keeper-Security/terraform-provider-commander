// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package driverlicense

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SharedAttributes returns the Driver's License resource attribute map shared between
// classic and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"account_number":  commonrecordsutils.OptionalAccountNumberField(),
			"name":            commonrecordsutils.NameNestedSchema(true),
			"birth_date":      commonrecordsutils.OptionalDateStringField("Birth date"),
			"address_ref":     commonrecordsutils.OptionalAddressRefField(),
			"expiration_date": commonrecordsutils.OptionalDateStringField("Expiration date"),
			"custom":          commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed Driver's License data source attributes
// shared between classic and new data sources. Callers add the lookup key (e.g.
// driver_license) and share-extension attributes separately.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"account_number": commonrecordsutils.ComputedAccountNumberField(),
			"name":           commonrecordsutils.NameDataSourceNestedSchema(),
			"birth_date": commonrecordsutils.ComputedDateStringField(
				"Date of birth (YYYY-MM-DD).",
				"Date of **birth** (`YYYY-MM-DD`).",
			),
			"address_ref": commonrecordsutils.ComputedAddressRefField(),
			"expiration_date": commonrecordsutils.ComputedDateStringField(
				"License expiration date (YYYY-MM-DD).",
				"License **expiration date** (`YYYY-MM-DD`).",
			),
			"custom": commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
