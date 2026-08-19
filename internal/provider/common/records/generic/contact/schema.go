// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SharedAttributes returns the Contact resource attribute map shared between
// classic and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"name":        commonrecordsutils.NameNestedSchema(false),
			"company":     commonrecordsutils.OptionalCompanyField(),
			"email":       commonrecordsutils.OptionalEmailField(),
			"phone":       commonrecordsutils.PhoneListSchema(),
			"address_ref": commonrecordsutils.OptionalAddressRefField(),
			"custom":      commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed Contact data source attributes
// shared between classic and new data sources. Callers add the lookup key
// (e.g. contact) and share-extension attributes separately.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"name":        commonrecordsutils.NameDataSourceNestedSchema(),
			"company":     commonrecordsutils.ComputedCompanyField(),
			"email":       commonrecordsutils.ComputedEmailField(),
			"phone":       commonrecordsutils.PhoneDataSourceListSchema(),
			"address_ref": commonrecordsutils.ComputedAddressRefField(),
			"custom":      commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
