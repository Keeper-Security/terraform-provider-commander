// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcompany

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *ManagedCompanyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ManagedCompanyResourceModel

	// Get current state (old values)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR,
			err.Error(),
		)
		return
	}

	err := utils.RunWithMspContext(ctx, r.ApiManager, func() error {
		if err := utils.MspDown(ctx, r.ApiManager); err != nil {
			return err
		}
		companyInfo, err := utils.FetchManagedCompanyByNameOrId(ctx, r.ApiManager, state.Id.ValueString())
		if err != nil {
			return err
		}
		if companyInfo == nil {
			resp.State.RemoveResource(ctx)
			return utils.ErrResourceRemoved
		}
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
			return fmt.Errorf("failed to build add_ons set: %v", diags.Errors())
		}
		state.AddOns = addOnsSet
		return nil
	}, "Read Managed Company Failed", &resp.Diagnostics)
	if err != nil && errors.Is(err, utils.ErrResourceRemoved) {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
