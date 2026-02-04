// Copyright (c) HashiCorp, Inc.
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
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Execute with managed company context if provided
	// ExecuteWithManagedCompanyContext handles context switching and enterprise-down internally
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, data.ManagedCompany, func() error {
		if err := addNodeBasicAttributes(ctx, r.apiManager, &data); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Create Enterprise Node Failed",
			err.Error(),
		)
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
		value := data.Parent.ValueString()

		// We need to make the parent as "root" as if the parent is the same as the managed company bec. this like functionality implemente in commander cli.
		if !data.ManagedCompany.IsNull() && data.Parent.ValueString() == data.ManagedCompany.ValueString() {
			value = "root"
		}

		parts = append(parts, fmt.Sprintf("--parent '%s'", value))
	}

	// TODO: Currently its not working in
	// if !data.WipeOut.IsNull() {
	// 	parts = append(parts, "--wipe-out")
	// }

	// TODO: NEED TO CHECK IF THIS FLAG IS REQUIRED / usecase CHECK?
	if !data.ToggleIsolated.IsNull() {
		parts = append(parts, "--toggle-isolated")
	}

	// TODO: NEED TO CHECK HOW WE CAN ADD LOGO FILE - NOT WORKING IN COMMANDER CLI
	// if !data.LogoFile.IsNull() {
	// 	parts = append(parts, fmt.Sprintf("--logo-file '%s'", data.LogoFile.ValueString()))
	// }

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
