// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpiseuser

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports import ID formats:
//   - User only: user email or user ID (e.g. user@example.com or 1326075447607317)
//   - With managed company: "managed_company_name_or_id,user_email_or_id" (comma-separated, e.g. "Test Company,1326075447607317" or "1169425105420462,user@example.com")
//
// After import, Terraform runs Read to refresh state from the API.
func (r *EnterpriseUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
			"Import ID cannot be empty. Use: (1) user email or user ID alone, e.g. user@example.com or 1326075447607317; or (2) for a user in a managed company, use \"managed_company_name_or_id,user_email_or_id\" (comma-separated), e.g. \"Test Company,1326075447607317\" or \"1169425105420462,user@example.com\".",
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
			"When using managed company format \"managed_company_name_or_id,user\", the user part cannot be empty. Examples: \"Test Company,1326075447607317\" or \"1169425105420462,user@example.com\".",
		)
		return
	}

	state := EnterpriseUserResourceModel{
		Id:             types.StringValue(resourceIdentifier), // id from import ID (email or user_id string)
		Email:          types.StringNull(),
		Name:           types.StringNull(),
		JobTitle:       types.StringNull(),
		Roles:          types.SetNull(types.StringType),
		Teams:          types.SetNull(types.StringType),
		Node:           types.StringNull(),
		ManagedCompany: types.StringNull(),
		Status:         types.StringNull(),
	}
	if managedCompany != "" {
		state.ManagedCompany = types.StringValue(managedCompany)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
