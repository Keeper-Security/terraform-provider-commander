package enterprisenode

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnterpriseNodesDataSourceModel struct {
	Node           types.String `tfsdk:"node"` // input: enterprise node name or ID to find the node
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Parent         types.String `tfsdk:"parent"`
	ParentId       types.String `tfsdk:"parent_id"`
	ManagedCompany types.String `tfsdk:"managed_company"`
}
