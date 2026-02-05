// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
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
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Execute with managed company context if provided
	// ExecuteWithManagedCompanyContext handles context switching and enterprise-down internally
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, state.ManagedCompany, func() error {

		nodeInfo, err := fetchEnterpriseNodeById(ctx, r.apiManager, state.Id.ValueString())
		if err != nil {
			return err
		}

		if nodeInfo == nil {
			// Resource not found - remove from state
			resp.State.RemoveResource(ctx)
			return utils.ErrResourceRemoved
		}

		// Map the response to the state
		state.Id = types.StringValue(strconv.Itoa(nodeInfo.NodeId))
		state.Name = types.StringValue(nodeInfo.Name)

		parentNodeVal, err := utils.RestoreUserInputFormatForNode(ctx, r.apiManager, nodeInfo.ParentNodeName, state.Parent)
		if err != nil {
			return fmt.Errorf("failed to convert parent node to original format: %w", err)
		}
		state.Parent = parentNodeVal

		return nil
	})

	if err != nil {
		if errors.Is(err, utils.ErrResourceRemoved) {
			return
		}
		resp.Diagnostics.AddError(
			"Read Enterprise Node Failed",
			err.Error(),
		)
		return
	}

	// Set the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func fetchEnterpriseNodeById(ctx context.Context, apiManager *api.ApiManager, id string) (*utils.EnterpriseNodeResponse, error) {

	// Build the Commander command string
	command := fmt.Sprintf("enterprise-info -n -v --format json --node '%s'", id)

	apiResp, err := apiManager.ExecuteCommand(ctx, command, "Unable to read enterprise node")
	if err != nil {
		return nil, err
	}

	// Parse the JSON response - it's an array of node objects
	var nodes []utils.EnterpriseNodeResponse

	if err := utils.UnmarshalApiResponse(apiResp.Data, &nodes); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise nodes list from API response: %w", err)
	}

	// Find the node matching state.Id (which is the node name)
	var nodeInfo *utils.EnterpriseNodeResponse

	for i := range nodes {
		// convert node id to string to compare with id
		if strconv.Itoa(nodes[i].NodeId) == id || nodes[i].Name == id {
			nodeInfo = &nodes[i]
			break
		}
	}

	// Node not in list - resource was likely deleted outside Terraform.
	// Return (nil, nil) so Read can remove it from state instead of erroring.
	if nodeInfo == nil {
		return nil, nil
	}

	return nodeInfo, nil
}
