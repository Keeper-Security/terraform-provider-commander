// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SharedAttributes returns the WiFi resource attribute map shared between classic
// and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"ssid": schema.StringAttribute{
				Required:            true,
				Description:         SSIDDescription,
				MarkdownDescription: SSIDMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("SSID", 1, false),
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
			"encryption": schema.StringAttribute{
				Optional:            true,
				Description:         EncryptionDescription,
				MarkdownDescription: EncryptionMarkdownDescription,
				Validators: []validator.String{
					utils.StringOneOfValidator("Encryption", AllowedEncryptions, true),
				},
			},
			"is_ssid_hidden": schema.BoolAttribute{
				Optional:            true,
				Description:         IsSSIDHiddenDescription,
				MarkdownDescription: IsSSIDHiddenMarkdownDescription,
			},
			"custom": commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed WiFi data source attributes shared
// between classic and new data sources. Callers add the lookup key (e.g. wifi)
// and share-extension attributes separately.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"ssid": dschema.StringAttribute{
				Computed:            true,
				Description:         SSIDDescription,
				MarkdownDescription: SSIDMarkdownDescription,
			},
			"password": dschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         PasswordDescription,
				MarkdownDescription: PasswordMarkdownDescription,
			},
			"encryption": dschema.StringAttribute{
				Computed:            true,
				Description:         EncryptionDescription,
				MarkdownDescription: EncryptionMarkdownDescription,
			},
			"is_ssid_hidden": dschema.BoolAttribute{
				Computed:            true,
				Description:         IsSSIDHiddenDescription,
				MarkdownDescription: IsSSIDHiddenMarkdownDescription,
			},
			"custom": commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
