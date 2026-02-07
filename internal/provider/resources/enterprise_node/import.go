// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports import ID formats:
//   - Node only: node name or node ID (e.g. "Root" or 1169425105420462)
//   - With managed company: "managed_company_name_or_id,node_name_or_id" (comma-separated, e.g. "Test Company,Root" or "1169425105420462,1169425105420462")
//
// After import, Terraform runs Read to refresh state from the API.
func (r *EnterpriseNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
			"Import ID cannot be empty. Use: (1) node name or node ID alone, e.g. Root or 1169425105420462; or (2) for a node in a managed company, use \"managed_company_name_or_id,node_name_or_id\" (comma-separated), e.g. \"Test Company,Root\" or \"1169425105420462,1169425105420462\".",
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
			"When using managed company format \"managed_company_name_or_id,node\", the node part cannot be empty. Examples: \"Test Company,Root\" or \"1169425105420462,1169425105420462\".",
		)
		return
	}

	state := EnterpriseNodeResourceModel{
		Id:     types.StringValue(resourceIdentifier),
		Name:   types.StringNull(),
		Parent: types.StringNull(),
		// WipeOut:        types.BoolNull(),
		ToggleIsolated: types.BoolNull(),
		ManagedCompany: types.StringNull(),
	}
	if managedCompany != "" {
		state.ManagedCompany = types.StringValue(managedCompany)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
