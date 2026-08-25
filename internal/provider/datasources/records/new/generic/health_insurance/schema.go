// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package healthinsurance

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordhealthinsurance "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/health_insurance"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *HealthInsuranceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"health_insurance": dschema.StringAttribute{
					Required:            true,
					Description:         "New (NSF) Health Insurance record UID to look up.",
					MarkdownDescription: "New (NSF) Health Insurance record **UID** to look up.",
				},
			},
			commonrecordhealthinsurance.SharedDataSourceAttributes(),
			new_share.DataSourceShareAttribute(),
		),
	}
}
