// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package database

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SharedAttributes returns the databaseCredentials resource attribute map shared between
// classic and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"login":    commonrecordsutils.OptionalLoginField(),
			"hostname": commonrecordsutils.OptionalHostnameField(),
			"port":     commonrecordsutils.OptionalPortField(),
			"type":     commonrecordsutils.OptionalDatabaseTypeField(),
			"password": commonrecordsutils.OptionalPasswordField(),
			"custom":   commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed databaseCredentials data source attributes
// shared between classic and new data sources.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"login":    commonrecordsutils.ComputedLoginField(),
			"hostname": commonrecordsutils.ComputedHostnameField(),
			"port":     commonrecordsutils.ComputedPortField(),
			"type":     commonrecordsutils.ComputedDatabaseTypeField(),
			"password": commonrecordsutils.ComputedPasswordField(),
			"custom":   commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
