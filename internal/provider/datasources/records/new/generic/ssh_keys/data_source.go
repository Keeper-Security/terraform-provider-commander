// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &SshKeysDataSource{}
var _ datasource.DataSourceWithConfigure = &SshKeysDataSource{}

type SshKeysDataSource struct {
	utils.BaseDataSource
}

func (d *SshKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_ssh_keys"
}

func (d *SshKeysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewSshKeysDataSource() datasource.DataSource {
	return &SshKeysDataSource{}
}
