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
			"configuration": schema.StringAttribute{
				Required:            true,
				Description:         ConfigurationDescription,
				MarkdownDescription: ConfigurationMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("configuration", 1, true),
				},
			},
			"gateway": schema.StringAttribute{
				Required:            true,
				Description:         GatewayDescription,
				MarkdownDescription: GatewayMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("gateway", 1, true),
				},
			},
			"custom": schema.ListNestedAttribute{
				Required:            true,
				Description:         commonrecordsutils.CustomDescription,
				MarkdownDescription: commonrecordsutils.CustomMarkdownDescription,
				Validators: []validator.List{
					RequiredSaasTypeCustomFieldValidator(),
					RequiredActiveCustomFieldValidator(),
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
			"configuration": dschema.StringAttribute{
				Computed:            true,
				Description:         "PAM Configuration UID used for rotation.",
				MarkdownDescription: "The PAM configuration UID used for rotation.",
			},
			"gateway": dschema.StringAttribute{
				Computed:            true,
				Description:         "The configured gateway `UID` or `name`",
				MarkdownDescription: "The configured gateway `UID` or `name`",
			},
			"custom": commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
