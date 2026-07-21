// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package membership

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &MembershipDataSource{}
var _ datasource.DataSourceWithConfigure = &MembershipDataSource{}

type MembershipDataSource struct {
	utils.BaseDataSource
}

func (d *MembershipDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_classic_membership"
}

func (d *MembershipDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewMembershipDataSource() datasource.DataSource {
	return &MembershipDataSource{}
}
