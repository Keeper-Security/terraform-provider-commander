// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package birthcertificate

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &BirthCertificateDataSource{}
var _ datasource.DataSourceWithConfigure = &BirthCertificateDataSource{}

type BirthCertificateDataSource struct {
	utils.BaseDataSource
}

func (d *BirthCertificateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_birth_certificate"
}

func (d *BirthCertificateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.ConfigureDataSource(ctx, req, resp)
}

func NewBirthCertificateDataSource() datasource.DataSource {
	return &BirthCertificateDataSource{}
}
