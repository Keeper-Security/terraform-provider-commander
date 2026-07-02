// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &PamMachineDataSource{}
var _ datasource.DataSourceWithConfigure = &PamMachineDataSource{}

// PamMachineDataSource reads a pamMachine vault record by UID.
type PamMachineDataSource struct {
	utils.BaseDataSource
}

func (d *PamMachineDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_classic_pam_machine"
}

func (d *PamMachineDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewPamMachineDataSource() datasource.DataSource {
	return &PamMachineDataSource{}
}
