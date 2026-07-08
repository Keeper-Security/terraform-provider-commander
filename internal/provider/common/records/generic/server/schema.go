// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package server

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SharedAttributes returns the serverCredentials resource attribute map shared between
// classic and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"login": schema.StringAttribute{
				Optional:            true,
				Description:         LoginDescription,
				MarkdownDescription: LoginMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Login", 1, true),
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
			"hostname": schema.StringAttribute{
				Optional:            true,
				Description:         HostnameDescription,
				MarkdownDescription: HostnameMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Hostname", 1, true),
				},
			},
			"port": schema.StringAttribute{
				Optional:            true,
				Description:         PortDescription,
				MarkdownDescription: PortMarkdownDescription,
				Validators: []validator.String{
					utils.NumericStringValidator("Port", true),
				},
			},
			"custom": commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed serverCredentials data source attributes
// shared between classic and new data sources.
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
			"hostname": dschema.StringAttribute{
				Computed:            true,
				Description:         HostnameDescription,
				MarkdownDescription: HostnameMarkdownDescription,
			},
			"port": dschema.StringAttribute{
				Computed:            true,
				Description:         PortDescription,
				MarkdownDescription: PortMarkdownDescription,
			},
			"custom": commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
