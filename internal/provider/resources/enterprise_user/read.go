// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser

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

func (r *EnterpriseUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnterpriseUserResourceModel

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
		userInfo, err := utils.FetchEnterpriseUserByEmailOrId(ctx, r.ApiManager, state.Id.ValueString())
		if err != nil {
			return err
		}
		if userInfo == nil {
			resp.State.RemoveResource(ctx)
			return utils.ErrResourceRemoved
		}
		if err := mapUserReadResponseToModel(ctx, r.ApiManager, *userInfo, &state); err != nil {
			return fmt.Errorf("failed to map user response to model: %w", err)
		}
		return nil
	}, "Read Enterprise User Failed", &resp.Diagnostics)
	if err != nil && errors.Is(err, utils.ErrResourceRemoved) {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func mapUserReadResponseToModel(ctx context.Context, apiManager *api.ApiManager, userInfo utils.EnterpriseUserResponse, state *EnterpriseUserResourceModel) error {
	// Map the response to the state
	state.Id = types.StringValue(strconv.Itoa(userInfo.UserId))
	state.Email = types.StringValue(userInfo.Email)
	if userInfo.Name == "" {
		state.Name = types.StringNull()
	} else {
		state.Name = types.StringValue(userInfo.Name)
	}
	if userInfo.JobTitle == "" {
		state.JobTitle = types.StringNull()
	} else {
		state.JobTitle = types.StringValue(userInfo.JobTitle)
	}
	nodeVal, err := utils.RestoreUserInputFormatForNode(ctx, apiManager, userInfo.Node, state.Node)
	if err != nil {
		return fmt.Errorf("failed to convert node to original format: %w", err)
	}
	state.Node = nodeVal

	switch userInfo.Status {
	case UserInvitedStatus:
		state.Status = types.StringValue(UserInvitedStatus)
	case UserActiveStatus:
		state.Status = types.StringValue(UserActiveStatus)
	case UserLockedStatus:
		state.Status = types.StringValue(UserLockedStatus)
	default:
		return fmt.Errorf("unexpected user status returned by commander cli: %s", userInfo.Status)
	}

	// Teams: preserve original format (name or team_uid) as user provided
	if len(userInfo.Teams) > 0 {
		teamsSet, err := utils.RestoreUserInputFormatForTeams(ctx, apiManager, userInfo.Teams, state.Teams)
		if err != nil {
			return fmt.Errorf("failed to convert teams to original format: %w", err)
		}
		state.Teams = teamsSet
	} else {
		state.Teams = types.SetNull(types.StringType)
	}

	// Roles: preserve original format (name or role_id) as user provided
	if len(userInfo.Roles) > 0 {
		rolesSet, err := utils.RestoreUserInputFormatForRoles(ctx, apiManager, userInfo.Roles, state.Roles)
		if err != nil {
			return fmt.Errorf("failed to convert roles to original format: %w", err)
		}
		state.Roles = rolesSet
	} else {
		state.Roles = types.SetNull(types.StringType)
	}

	return nil
}
