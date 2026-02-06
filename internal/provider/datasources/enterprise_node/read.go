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
	if err := d.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	err := utils.ExecuteWithManagedCompanyContext(ctx, d.apiManager, data.ManagedCompany, func() error {
		nodeInfo, err := utils.FetchEnterpriseNodeByNameOrId(ctx, d.apiManager, data.Node.ValueString())
		if err != nil {
			return err
		}
		if nodeInfo == nil {
			return fmt.Errorf("Enterprise node not found")
		}
		data.Id = types.StringValue(strconv.Itoa(nodeInfo.NodeId))
		data.Name = types.StringValue(nodeInfo.Name)
		data.Parent = types.StringValue(nodeInfo.ParentNodeName)
		data.ParentId = types.StringValue(strconv.Itoa(nodeInfo.ParentNodeId))
		data.ManagedCompany = types.StringNull() // Bec we dotn want to return the managed company in the result
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Enterprise Node Failed",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
