// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package login

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SharedAttributes returns the Login resource attribute map shared between classic
// and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"login": schema.StringAttribute{
				Required:            true,
				Description:         LoginDescription,
				MarkdownDescription: LoginMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Login", 1, false),
				},
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         PasswordDescription,
				MarkdownDescription: PasswordMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Password", 1, true),
				},
			},
			"website_address": schema.StringAttribute{
				Optional:            true,
				Description:         WebsiteAddressDescription,
				MarkdownDescription: WebsiteAddressMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Website address", 1, true),
				},
			},
			"custom": customFieldResourceAttribute(),
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
			"login": dschema.StringAttribute{
				Computed:            true,
				Description:         LoginDescription,
				MarkdownDescription: LoginMarkdownDescription,
			},
			"password": dschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         PasswordDescription,
				MarkdownDescription: PasswordMarkdownDescription,
			},
			"website_address": dschema.StringAttribute{
				Computed:            true,
				Description:         WebsiteAddressDescription,
				MarkdownDescription: WebsiteAddressMarkdownDescription,
			},
			"custom": customFieldDataSourceAttribute(),
		},
	)
}

func customFieldResourceAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:            true,
		Description:         CustomDescription,
		MarkdownDescription: CustomMarkdownDescription,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:            true,
					Description:         CustomTypeDescription,
					MarkdownDescription: CustomTypeMarkdownDescription,
				},
				"label": schema.StringAttribute{
					Required:            true,
					Description:         CustomLabelDescription,
					MarkdownDescription: CustomLabelMarkdownDescription,
					Validators: []validator.String{
						utils.StringMinLengthValidator("Custom field label", 1, false),
					},
				},
				"value": schema.StringAttribute{
					Required:            true,
					Description:         CustomValueDescription,
					MarkdownDescription: CustomValueMarkdownDescription,
				},
				"sensitive": schema.BoolAttribute{
					Optional:            true,
					Computed:            true,
					Description:         CustomSensitiveDescription,
					MarkdownDescription: CustomSensitiveMarkdownDescription,
					Default:             booldefault.StaticBool(false),
				},
			},
		},
	}
}

func customFieldDataSourceAttribute() dschema.ListNestedAttribute {
	return dschema.ListNestedAttribute{
		Computed:            true,
		Description:         DSCustomDescription,
		MarkdownDescription: DSCustomMarkdownDescription,
		NestedObject: dschema.NestedAttributeObject{
			Attributes: map[string]dschema.Attribute{
				"type": dschema.StringAttribute{
					Computed:            true,
					Description:         DSCustomTypeDescription,
					MarkdownDescription: DSCustomTypeMarkdownDescription,
				},
				"label": dschema.StringAttribute{
					Computed:            true,
					Description:         CustomLabelDescription,
					MarkdownDescription: CustomLabelMarkdownDescription,
				},
				"value": dschema.StringAttribute{
					Computed:            true,
					Sensitive:           true,
					Description:         DSCustomValueDescription,
					MarkdownDescription: DSCustomValueMarkdownDescription,
				},
				"sensitive": dschema.BoolAttribute{
					Computed:            true,
					Description:         DSCustomSensitiveDescription,
					MarkdownDescription: DSCustomSensitiveMarkdownDescription,
				},
			},
		},
	}
}
