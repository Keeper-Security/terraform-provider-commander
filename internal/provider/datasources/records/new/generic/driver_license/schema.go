// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package driverlicense

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecorddriverlicense "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/driver_license"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *DriverLicenseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"driver_license": dschema.StringAttribute{
					Required:            true,
					Description:         "Driver's License record title or UID to look up.",
					MarkdownDescription: "Driver's License record **title** or **UID** to look up.",
				},
			},
			commonrecorddriverlicense.SharedDataSourceAttributes(),
			new_share.DataSourceShareAttribute(),
		),
	}
}
