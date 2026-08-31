// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package birthcertificate

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SharedAttributes returns the Birth Certificate resource attribute map shared between
// classic and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"name":       commonrecordsutils.NameNestedSchema(true),
			"birth_date": commonrecordsutils.OptionalDateStringField("Birth date"),
			"custom":     commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed Birth Certificate data source attributes
// shared between classic and new data sources. Callers add the lookup key (e.g.
// birth_certificate) and share-extension attributes separately.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"name": commonrecordsutils.NameDataSourceNestedSchema(),
			"birth_date": commonrecordsutils.ComputedDateStringField(
				"Date of birth (YYYY-MM-DD).",
				"Date of **birth** (`YYYY-MM-DD`).",
			),
			"custom": commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
