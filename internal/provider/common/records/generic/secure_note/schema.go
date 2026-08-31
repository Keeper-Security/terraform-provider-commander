// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package securenote

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SharedAttributes returns the encryptedNotes resource attribute map shared between
// classic and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"secured_note": commonrecordsutils.OptionalSecuredNoteField(),
			"date":         commonrecordsutils.OptionalDateStringField("Date"),
			"custom":       commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed encryptedNotes data source attributes
// shared between classic and new data sources.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"secured_note": commonrecordsutils.ComputedSecuredNoteField(),
			"date": commonrecordsutils.ComputedDateStringField(
				"Record date (YYYY-MM-DD).",
				"Record **date** (`YYYY-MM-DD`).",
			),
			"custom": commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
