// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sshkeys

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordsshkeys "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/ssh_keys"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *SshKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: utils.MergeDataSourceAttributes(
			map[string]dschema.Attribute{
				"ssh_keys": dschema.StringAttribute{
					Required:            true,
					Description:         "SSH keys record title or UID to look up.",
					MarkdownDescription: "SSH keys record **title** or **UID** to look up.",
				},
			},
			commonrecordsshkeys.SharedDataSourceAttributes(),
			classic_share.DataSourceShareAttribute(),
		),
	}
}
