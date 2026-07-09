// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package softwarelicense

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SharedAttributes returns the softwareLicense resource attribute map shared between
// classic and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"software_license_key": commonrecordsutils.OptionalSoftwareLicenseKeyField(),
			"expiration_date":      commonrecordsutils.OptionalDateStringField("Expiration date"),
			"date_active":          commonrecordsutils.OptionalDateStringField("Date active"),
			"custom":               commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed softwareLicense data source attributes
// shared between classic and new data sources.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"software_license_key": commonrecordsutils.ComputedSoftwareLicenseKeyField(),
			"expiration_date": commonrecordsutils.ComputedDateStringField(
				"License expiration date (YYYY-MM-DD).",
				"License **expiration date** (`YYYY-MM-DD`).",
			),
			"date_active": commonrecordsutils.ComputedDateStringField(
				"Date the license became active (YYYY-MM-DD).",
				"Date the license became **active** (`YYYY-MM-DD`).",
			),
			"custom": commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
