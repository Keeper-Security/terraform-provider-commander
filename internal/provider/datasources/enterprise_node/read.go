// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *EnterpriseNodesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EnterpriseNodesDataSourceModel

	// Get configuration data
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	if err := utils.RunWithManagedCompanyContext(ctx, d.ApiManager, data.ManagedCompany, func() error {
		nodeInfo, err := utils.FetchEnterpriseNodeByNameOrId(ctx, d.ApiManager, data.Node.ValueString())
		if err != nil {
			return err
		}
		if nodeInfo == nil {
			return fmt.Errorf("enterprise node: '%s' not found", data.Node.ValueString())
		}
		data.Id = types.StringValue(strconv.Itoa(nodeInfo.NodeId))
		data.Name = types.StringValue(nodeInfo.Name)
		data.Parent = types.StringValue(nodeInfo.ParentNodeName)
		data.ParentId = types.StringValue(strconv.Itoa(nodeInfo.ParentNodeId))
		data.ManagedCompany = types.StringNull()
		return nil
	}, "Read Enterprise Node Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
