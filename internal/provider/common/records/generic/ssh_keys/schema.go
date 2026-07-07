// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SharedAttributes returns the sshKeys resource attribute map shared between classic
// and new resources. Callers add any share-extension attribute separately.
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
			"passphrase": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         PassphraseDescription,
				MarkdownDescription: PassphraseMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Passphrase", 1, true),
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
					utils.StringMinLengthValidator("Port", 1, true),
				},
			},
			"public_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         PublicKeyDescription,
				MarkdownDescription: PublicKeyMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Public key", 1, true),
				},
			},
			"private_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         PrivateKeyDescription,
				MarkdownDescription: PrivateKeyMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Private key", 1, true),
				},
			},
			"custom": commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed sshKeys data source attributes shared
// between classic and new data sources.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"login": dschema.StringAttribute{
				Computed:            true,
				Description:         LoginDescription,
				MarkdownDescription: LoginMarkdownDescription,
			},
			"passphrase": dschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         PassphraseDescription,
				MarkdownDescription: PassphraseMarkdownDescription,
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
			"public_key": dschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         PublicKeyDescription,
				MarkdownDescription: PublicKeyMarkdownDescription,
			},
			"private_key": dschema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         PrivateKeyDescription,
				MarkdownDescription: PrivateKeyMarkdownDescription,
			},
			"custom": commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
}
