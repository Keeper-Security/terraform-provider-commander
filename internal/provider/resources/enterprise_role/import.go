package enterpriserole

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports import ID formats:
//   - Role only: role name or role ID (e.g. "Admin" or 1234567890)
//   - With managed company: "managed_company_name_or_id,role_name_or_id" (comma-separated, e.g. "Test Company,Admin" or "1169425105420462,1234567890")
//
// After import, Terraform runs Read to refresh state from the API.
func (r *EnterpriseRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	importID := strings.TrimSpace(req.ID)
	if importID == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID cannot be empty. Use: (1) role name or role ID alone, e.g. Admin or 1234567890; or (2) for a role in a managed company, use \"managed_company_name_or_id,role_name_or_id\" (comma-separated), e.g. \"Test Company,Admin\" or \"1169425105420462,1234567890\".",
		)
		return
	}

	var resourceIdentifier, managedCompany string
	if parts := strings.SplitN(importID, ",", 2); len(parts) == 2 {
		managedCompany = strings.TrimSpace(parts[0])
		resourceIdentifier = strings.TrimSpace(parts[1])
	} else {
		resourceIdentifier = importID
	}

	if resourceIdentifier == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"When using managed company format \"managed_company_name_or_id,role\", the role part cannot be empty. Examples: \"Test Company,Admin\" or \"1169425105420462,1234567890\".",
		)
		return
	}

	managingNodesElemType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"privileges": types.SetType{ElemType: types.StringType},
			"cascade":    types.BoolType,
		},
	}

	state := EnterpriseRoleResourceModel{
		Id:                  types.StringValue(resourceIdentifier),
		Name:                types.StringNull(),
		Node:                types.StringNull(),
		Users:               types.SetNull(types.StringType),
		Teams:               types.SetNull(types.StringType),
		ManagingNodes:       types.MapNull(managingNodesElemType),
		EnforcementPolicies: types.MapNull(types.StringType),
		ManagedCompany:      types.StringNull(),
	}
	if managedCompany != "" {
		state.ManagedCompany = types.StringValue(managedCompany)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
