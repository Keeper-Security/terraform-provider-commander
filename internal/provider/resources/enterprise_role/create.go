// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EnterpriseRoleResourceModel

	// Get planned data from Terraform
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that only one of teams or managing_nodes is provided
	if err := validateTeamsAndManagingNodesMutualExclusivity(data.Teams, data.ManagingNodes); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			err.Error(),
		)
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

	if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, data.ManagedCompany, func() error {
		users, err := utils.FetchAndProcessUsers(ctx, r.ApiManager, types.SetNull(types.StringType), data.Users)
		if err != nil {
			return err
		}
		teams, err := utils.FetchAndProcessTeams(ctx, r.ApiManager, types.SetNull(types.StringType), data.Teams)
		if err != nil {
			return err
		}
		if err := addRoleBasicAttributes(ctx, r.ApiManager, &data); err != nil {
			return err
		}
		currentScopeNodes, err := r.ApiManager.ExecuteCommand(ctx, "enterprise-info -n --format json", "Unable to fetch enterprise nodes for the managed company")
		if err != nil {
			return err
		}
		if err := validateManagingNodes(ctx, data.ManagingNodes, currentScopeNodes.Data); err != nil {
			return fmt.Errorf("managing nodes validation failed: %w", err)
		}
		if err := processManagingNodes(ctx, r.ApiManager, data.Id.ValueString(), data.ManagingNodes); err != nil {
			return err
		}
		if err := processEnforcementPolicies(ctx, r.ApiManager, data.Id.ValueString(), data.EnforcementPolicies); err != nil {
			return err
		}
		if users != "" {
			command := fmt.Sprintf("enterprise-role '%s' -f %s", data.Id.ValueString(), users)
			if _, err = r.ApiManager.ExecuteCommand(ctx, command, "Unable to add users to the enterprise role"); err != nil {
				return err
			}
		}
		if teams != "" {
			command := fmt.Sprintf("enterprise-role '%s' -f %s", data.Id.ValueString(), teams)
			if _, err = r.ApiManager.ExecuteCommand(ctx, command, "Unable to add teams to the enterprise role"); err != nil {
				return err
			}
		}
		return nil
	}, "Create Enterprise Role Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the ID in the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func addRoleBasicAttributes(ctx context.Context, apiManager *api.ApiManager, data *EnterpriseRoleResourceModel) error {
	var parts []string

	parts = append(parts, "enterprise-role")

	// Required parameters
	parts = append(parts, fmt.Sprintf("--add '%s'", data.Name.ValueString()))

	// Optional parameters
	if !data.Node.IsNull() {
		parts = append(parts, fmt.Sprintf("--node '%s'", data.Node.ValueString()))
	}

	command := strings.Join(parts, " ")

	createdRoleResponse, err := apiManager.ExecuteCommand(ctx, command, "Unable to add basic role attributes to the enterprise role")
	if err != nil {
		return err
	}

	createdRoleId, isCreatedRoleIdExtracted := extractRoleIdFromCreateRoleResponse(string(createdRoleResponse.Message))

	if isCreatedRoleIdExtracted {
		data.Id = types.StringValue(createdRoleId)
	} else {
		return fmt.Errorf("failed to extract role id from create role response: %s", string(createdRoleResponse.Message))
	}

	return nil
}
