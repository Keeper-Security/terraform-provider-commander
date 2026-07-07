// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package login

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
			"custom": commonrecordsutils.CustomFieldAttributeSchema(),
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
			"custom": commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
