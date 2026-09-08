// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package login

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SharedAttributes returns the Login resource attribute map shared between classic
// and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			// need to make login as optional
			"login":           commonrecordsutils.RequiredLoginField(),
			"password":        commonrecordsutils.OptionalPasswordField(),
			"website_address": commonrecordsutils.OptionalWebsiteAddressField(),
			"custom":          commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed Login data source attributes shared
// between classic and new data sources. Callers add the lookup key (e.g. login)
// and share-extension attributes separately.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"login":           commonrecordsutils.ComputedLoginField(),
			"password":        commonrecordsutils.ComputedPasswordField(),
			"website_address": commonrecordsutils.ComputedWebsiteAddressField(),
			"custom":          commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
