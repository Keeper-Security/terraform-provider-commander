// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package birthcertificate

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordbirthcertificate "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/birth_certificate"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *BirthCertificateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"birth_certificate": dschema.StringAttribute{
					Required:            true,
					Description:         "New (NSF) Birth Certificate record title or UID to look up.",
					MarkdownDescription: "New (NSF) Birth Certificate record **title** or **UID** to look up.",
				},
			},
			commonrecordbirthcertificate.SharedDataSourceAttributes(),
			new_share.DataSourceShareAttribute(),
		),
	}
}
