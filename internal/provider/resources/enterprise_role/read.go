// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

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

func (r *EnterpriseRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnterpriseRoleResourceModel

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

	err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, state.ManagedCompany, func() error {
		roleInfo, err := utils.FetchEnterpriseRoleByNameOrId(ctx, r.ApiManager, state.Id.ValueString())
		if err != nil {
			return err
		}
		if roleInfo == nil {
			resp.State.RemoveResource(ctx)
			return utils.ErrResourceRemoved
		}
		if err := mapRoleReadResponseToModel(ctx, r.ApiManager, *roleInfo, &state); err != nil {
			return err
		}
		return nil
	}, "Read Enterprise Role Failed", &resp.Diagnostics)
	if err != nil && errors.Is(err, utils.ErrResourceRemoved) {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Note: We will remove stateId from the function parameters in the future, when we will have role_id in the response while creating the role.
func mapRoleReadResponseToModel(ctx context.Context, apiManager *api.ApiManager, roleInfo utils.EnterpriseRoleResponse, state *EnterpriseRoleResourceModel) error {
	// Map the response to the state
	state.Id = types.StringValue(strconv.Itoa(roleInfo.RoleId))
	state.Name = types.StringValue(roleInfo.Name)
	nodeVal, err := utils.RestoreUserInputFormatForNode(ctx, apiManager, roleInfo.Node, state.Node)
	if err != nil {
		return fmt.Errorf("failed to convert node to original format: %w", err)
	}
	state.Node = nodeVal

	// Convert API response identifiers back to original format from state
	// Users: preserve original format (email or ID) as user provided
	if len(roleInfo.Users) > 0 {
		usersSet, err := utils.RestoreUserInputFormatForUsers(ctx, apiManager, roleInfo.Users, state.Users)
		if err != nil {
			return fmt.Errorf("failed to convert users to original format: %w", err)
		}
		state.Users = usersSet
	} else {
		state.Users = types.SetNull(types.StringType)
	}

	// Teams: preserve original format (name or team_uid) as user provided
	if len(roleInfo.Teams) > 0 {
		teamsSet, err := utils.RestoreUserInputFormatForTeams(ctx, apiManager, roleInfo.Teams, state.Teams)
		if err != nil {
			return fmt.Errorf("failed to convert teams to original format: %w", err)
		}
		state.Teams = teamsSet
	} else {
		state.Teams = types.SetNull(types.StringType)
	}

	// Managing nodes: map API managed_nodes_permissions to state
	managingNodesMap, err := utils.MapManagedNodesPermissionsToState(ctx, roleInfo.ManagedNodesPermissions)
	if err != nil {
		return fmt.Errorf("failed to map managing nodes to state: %w", err)
	}
	state.ManagingNodes = managingNodesMap

	// Enforcement policies: map API enforcements (key -> string value) to state map; keys normalized to UPPER_SNAKE_CASE
	enforcementPoliciesMap, err := utils.MapEnforcementsToState(roleInfo.Enforcements, GeneratedPasswordComplexity)
	if err != nil {
		return fmt.Errorf("failed to map enforcement policies to state: %w", err)
	}
	state.EnforcementPolicies = enforcementPoliciesMap

	return nil
}
