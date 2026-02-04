// Copyright (c) HashiCorp, Inc.
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
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, state.ManagedCompany, func() error {
		// Build the Commander command string
		command := fmt.Sprintf("enterprise-info '%s' -r --format json --columns='visible_below,default_role,admin,node,user_count,users,team_count,teams' -q", state.Id.ValueString())

		apiResp, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to read enterprise role")
		if err != nil {
			return fmt.Errorf("Read Enterprise Role Failed: %w", err)
		}

		// Parse the JSON response - it's an array of role objects
		var roles []utils.EnterpriseRoleResponse

		// Unmarshal API response into roles struct
		if err := utils.UnmarshalApiResponse(apiResp.Data, &roles); err != nil {
			return fmt.Errorf("unable to parse enterprise roles list from API response: %w", err)
		}

		// Find the role matching state.Id
		var roleInfo *utils.EnterpriseRoleResponse
		stateId := state.Id.ValueString()
		for i := range roles {
			if strconv.Itoa(roles[i].RoleId) == stateId || roles[i].Name == stateId {
				roleInfo = &roles[i]
				break
			}
		}

		if roleInfo == nil {
			// Resource not found - remove from state
			resp.State.RemoveResource(ctx)
			return utils.ErrResourceRemoved
		}

		if err := mapRoleReadResponseToModel(ctx, r.apiManager, *roleInfo, &state); err != nil {
			return fmt.Errorf("failed to map role response to model: %w", err)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, utils.ErrResourceRemoved) {
			return
		}
		resp.Diagnostics.AddError(
			"Read Enterprise Role Failed",
			err.Error(),
		)
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
	state.Node = types.StringValue(utils.ExtractNodeName(roleInfo.Node))

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

	// TODO: Later we will add logic for enforcement policies, managing nodes to state when it is implemented in commander cli

	return nil
}
