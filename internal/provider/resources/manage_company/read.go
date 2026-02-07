// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *ManageCompanyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ManageCompanyResourceModel

	// Get current state (old values)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Only switch to MSP if it's an MSP account
	if r.apiManager.IsMspAccount {
		if err := utils.SwitchToMsp(ctx, r.apiManager); err != nil {
			resp.Diagnostics.AddError(
				"Read Managed Company Failed",
				fmt.Sprintf("Failed to switch to MSP context: %s", err.Error()),
			)
			return
		}
	}

	// Execute msp-down command (setup/initialization)
	if err := utils.MspDown(ctx, r.apiManager); err != nil {
		resp.Diagnostics.AddError(
			"Read Managed Company Failed",
			err.Error(),
		)
		return
	}

	companyInfo, err := utils.FetchManageCompanyByNameOrId(ctx, r.apiManager, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Managed Company Failed",
			err.Error(),
		)
		return
	}

	if companyInfo == nil {
		// Resource not found - remove from state
		// Terraform will detect this and mark the resource for destruction
		resp.State.RemoveResource(ctx)
		return
	}

	// Map the response to the model
	state.Id = types.StringValue(strconv.Itoa(companyInfo.CompanyId))
	state.Name = types.StringValue(companyInfo.CompanyName)

	state.Node = types.StringValue(utils.ExtractNodeName(companyInfo.NodeName))

	state.Plan = types.StringValue(companyInfo.Plan)
	state.Seats = types.Int64Value(int64(companyInfo.Allocated))

	// Convert storage format: "100GB" -> "100gb", "1TB" -> "1tb", "10TB" -> "10tb"
	storageLower := strings.ToLower(companyInfo.Storage)
	state.FilePlan = types.StringValue(storageLower)

	// Convert addons array of strings to types.Set
	addOnsElements := make([]types.String, len(companyInfo.Addons))
	for i, addon := range companyInfo.Addons {
		addOnsElements[i] = types.StringValue(addon)
	}
	addOnsSet, diags := types.SetValueFrom(ctx, types.StringType, addOnsElements)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	state.AddOns = addOnsSet

	// Set the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
