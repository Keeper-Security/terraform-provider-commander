// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EnterpriseNodeResourceModel

	// Get planned data from Terraform
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Restrict use of toggle_isolated in create bec cli does not support it
	if !data.ToggleIsolated.IsNull() && data.ToggleIsolated.ValueBool() {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"toggle_isolated is not supported in create",
		)
		return
	}

	// Execute with managed company context if provided
	if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, data.ManagedCompany, func() error {
		return addNodeBasicAttributes(ctx, r.ApiManager, &data)
	}, "Create Enterprise Node Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the ID in the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func addNodeBasicAttributes(ctx context.Context, apiManager *api.ApiManager, data *EnterpriseNodeResourceModel) error {
	var parts []string

	parts = append(parts, "enterprise-node")

	// Required parameters
	parts = append(parts, fmt.Sprintf("--add '%s'", data.Name.ValueString()))

	// Optional parameters
	if !data.Parent.IsNull() {
		parts = append(parts, fmt.Sprintf("--parent '%s'", data.Parent.ValueString()))
	}

	command := strings.Join(parts, " ")

	createNodeResponse, err := apiManager.ExecuteCommand(ctx, command, "Unable to create enterprise node")
	if err != nil {
		return fmt.Errorf("create enterprise node failed: %w", err)
	}

	createdNodeId, isCreatedNodeIdExtracted := utils.ExtractNodeIDFromCreateNodeResponse(string(createNodeResponse.Message))

	if isCreatedNodeIdExtracted {
		data.Id = types.StringValue(createdNodeId)
	} else {
		return fmt.Errorf("failed to extract node id from create node response: %s", string(createNodeResponse.Message))
	}
	return nil
}
