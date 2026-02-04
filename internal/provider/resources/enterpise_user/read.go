package enterpiseuser

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
		command := fmt.Sprintf("enterprise-info '%s' -u --format json --columns='name,status,node,teams,roles,alias' -q", state.Id.ValueString())

		apiResp, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to read enterprise user")
		if err != nil {
			return fmt.Errorf("Read Enterprise User Failed: %w", err)
		}

		// Parse the JSON response - it's an array of user objects
		var users []utils.EnterpriseUserResponse

		// Unmarshal API response into users struct
		if err := utils.UnmarshalApiResponse(apiResp.Data, &users); err != nil {
			return fmt.Errorf("unable to parse enterprise users list from API response: %w", err)
		}

		// Find the user matching state.Id
		var userInfo *utils.EnterpriseUserResponse
		stateId := state.Id.ValueString()
		for i := range users {
			// TODO: Later will check user_id instead of email
			if users[i].Email == stateId || strconv.Itoa(users[i].UserId) == stateId {
				userInfo = &users[i]
				break
			}
		}

		if userInfo == nil {
			// Resource not found - remove from state
			resp.State.RemoveResource(ctx)
			return utils.ErrResourceRemoved
		}

		if err := mapUserReadResponseToModel(ctx, r.apiManager, *userInfo, &state); err != nil {
			return fmt.Errorf("failed to map user response to model: %w", err)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, utils.ErrResourceRemoved) {
			return
		}
		resp.Diagnostics.AddError(
			"Read Enterprise User Failed",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func mapUserReadResponseToModel(ctx context.Context, apiManager *api.ApiManager, userInfo utils.EnterpriseUserResponse, state *EnterpriseUserResourceModel) error {
	// Map the response to the state
	state.Id = types.StringValue(strconv.Itoa(userInfo.UserId)) // NOTE: For now we are using email as id, once we get user_id in commander cli response while creating user we will change to Int64 type
	state.Email = types.StringValue(userInfo.Email)
	state.Name = types.StringValue(userInfo.Name)
	state.JobTitle = types.StringValue(userInfo.JobTitle)
	state.Node = types.StringValue(utils.ExtractNodeName(userInfo.Node))

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
