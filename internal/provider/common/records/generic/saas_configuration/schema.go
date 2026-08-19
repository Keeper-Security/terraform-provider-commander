// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SharedAttributes returns the saasConfiguration resource attribute map shared between
// classic and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"custom": schema.ListNestedAttribute{
				Required:            true,
				Description:         commonrecordsutils.CustomDescription,
				MarkdownDescription: commonrecordsutils.CustomMarkdownDescription,
				Validators: []validator.List{
					RequiredSaasTypeCustomFieldValidator(),
				},
				NestedObject: commonrecordsutils.CustomFieldNestedAttributeObject(),
			},
		},
	)
}

// SharedDataSourceAttributes returns computed saasConfiguration data source attributes
// shared between classic and new data sources.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"custom": commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
