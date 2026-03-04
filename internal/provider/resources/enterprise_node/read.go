// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseNodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnterpriseNodeResourceModel

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

	// Execute with managed company context if provided
	err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, state.ManagedCompany, func() error {
		nodeInfo, err := utils.FetchEnterpriseNodeByNameOrId(ctx, r.ApiManager, state.Id.ValueString())
		if err != nil {
			return err
		}
		if nodeInfo == nil {
			resp.State.RemoveResource(ctx)
			return utils.ErrResourceRemoved
		}
		state.Id = types.StringValue(strconv.Itoa(nodeInfo.NodeId))
		state.Name = types.StringValue(nodeInfo.Name)
		state.ToggleIsolated = types.BoolValue(nodeInfo.Isolated)

		parentNodeVal, err := utils.RestoreUserInputFormatForNode(ctx, r.ApiManager, nodeInfo.ParentNodeName, state.Parent)
		if err != nil {
			return fmt.Errorf("failed to convert parent node to original format: %w", err)
		}
		state.Parent = parentNodeVal
		return nil
	}, "Read Enterprise Node Failed", &resp.Diagnostics)
	if err != nil && errors.Is(err, utils.ErrResourceRemoved) {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
