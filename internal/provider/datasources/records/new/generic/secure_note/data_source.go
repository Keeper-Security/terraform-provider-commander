// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package securenote

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &SecureNoteDataSource{}
var _ datasource.DataSourceWithConfigure = &SecureNoteDataSource{}

type SecureNoteDataSource struct {
	utils.BaseDataSource
}

func (d *SecureNoteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_secure_note"
}

func (d *SecureNoteDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewSecureNoteDataSource() datasource.DataSource {
	return &SecureNoteDataSource{}
}
