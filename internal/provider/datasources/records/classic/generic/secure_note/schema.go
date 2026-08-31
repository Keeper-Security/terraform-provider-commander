// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package securenote

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordsecurenote "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/secure_note"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *SecureNoteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"secure_note": dschema.StringAttribute{
					Required:            true,
					Description:         "Secure note record UID to look up.",
					MarkdownDescription: "Secure note record **UID** to look up.",
				},
			},
			commonrecordsecurenote.SharedDataSourceAttributes(),
			classic_share.DataSourceShareAttribute(),
		),
	}
}
