// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcompany

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportState supports import ID format:
//   - Company name or company ID (e.g. "Test Company" or 1169425105420462)
//
// Managed companies are listed in MSP context. After import, Terraform runs Read to refresh state from the API.
func (r *ManagedCompanyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if err := r.EnsureApiManager(); err != nil {
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
			"Import ID cannot be empty. Use the managed company name or company ID, e.g. \"Test Company\" or 1169425105420462.",
		)
		return
	}

	state := ManagedCompanyResourceModel{
		Id:       types.StringValue(importID),
		Name:     types.StringNull(),
		Node:     types.StringNull(),
		Plan:     types.StringNull(),
		Seats:    types.Int64Null(),
		FilePlan: types.StringNull(),
		AddOns:   types.SetNull(types.StringType),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
