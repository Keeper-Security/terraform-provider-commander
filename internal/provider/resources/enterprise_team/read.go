package enterpriseteam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseTeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnterpriseTeamResourceModel

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

	// Use managed company from state if provided
	managedCompany := state.ManagedCompany

	// Execute with managed company context if provided
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, managedCompany, func() error {
		// Build command to get enterprise team info
		command := fmt.Sprintf("enterprise-info '%s' -t --format json --columns='users,roles,restricts,node' -q", state.Id.ValueString())

		// Execute the command
		apiResp, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to retrieve enterprise team information")
		if err != nil {
			return fmt.Errorf("Read Enterprise Team Failed: %w", err)
		}

		// Parse the response
		teams, err := parseEnterpriseTeamReadResponse(apiResp.Data)
		if err != nil {
			return fmt.Errorf("Failed to parse team response: %w", err)
		}

		// Find the team matching the ID
		var teamInfo *utils.EnterpriseTeamResponse
		for i := range teams {
			if teams[i].TeamUid == state.Id.ValueString() {
				teamInfo = &teams[i]
				break
			}
		}

		if teamInfo == nil {
			// Resource not found - remove from state
			resp.State.RemoveResource(ctx)
			return utils.ErrResourceRemoved
		}

		// Map the response to the model
		if err := mapTeamReadResponseToModel(ctx, r.apiManager, *teamInfo, &state); err != nil {
			return fmt.Errorf("Failed to map team response to model: %w", err)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, utils.ErrResourceRemoved) {
			return
		}
		resp.Diagnostics.AddError(
			"Read Enterprise Team Failed",
			err.Error(),
		)
		return
	}

	// Set the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// parseEnterpriseTeamReadResponse parses the JSON response from enterprise-info -t command
func parseEnterpriseTeamReadResponse(data interface{}) ([]utils.EnterpriseTeamResponse, error) {
	var teams []utils.EnterpriseTeamResponse

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to process the response from Keeper Commander Service Mode API: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &teams); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise team information from Service Mode API response: %w", err)
	}

	return teams, nil
}

// parseRestrictsString parses the restricts string (e.g., "R W S") and returns boolean values
// R -> enable_privacy_screen (RestrictView)
// S -> restrict_record_re_share (RestrictShare)
// W -> restrict_record_edit (RestrictEdit)
func parseRestrictsString(restricts string) (restrictEdit, restrictShare, restrictView bool) {
	restricts = strings.TrimSpace(restricts)
	if restricts == "" {
		return false, false, false
	}

	// Split by space and check for each flag
	parts := strings.Fields(restricts)
	for _, part := range parts {
		switch part {
		case "W":
			restrictEdit = true
		case "S":
			restrictShare = true
		case "R":
			restrictView = true
		}
	}

	return restrictEdit, restrictShare, restrictView
}

// mapTeamReadResponseToModel maps the API response to the resource model
func mapTeamReadResponseToModel(ctx context.Context, apiManager *api.ApiManager, teamResp utils.EnterpriseTeamResponse, state *EnterpriseTeamResourceModel) error {
	// Map ID
	state.Id = types.StringValue(teamResp.TeamUid)

	// Map Name
	state.Name = types.StringValue(teamResp.Name)

	// Parse and map restricts flags
	restrictEdit, restrictShare, restrictView := parseRestrictsString(teamResp.Restricts)
	state.RestrictEdit = types.BoolValue(restrictEdit)
	state.RestrictShare = types.BoolValue(restrictShare)
	state.RestrictView = types.BoolValue(restrictView)

	// Convert API response identifiers back to original format from state
	// Roles: preserve original format (name or ID) as user provided
	if len(teamResp.Roles) > 0 {
		rolesSet, err := utils.RestoreUserInputFormatForRoles(ctx, apiManager, teamResp.Roles, state.Roles)
		if err != nil {
			return fmt.Errorf("failed to convert roles to original format: %w", err)
		}
		state.Roles = rolesSet
	} else {
		state.Roles = types.SetNull(types.StringType)
	}

	// Users: preserve original format (email or ID) as user provided
	if len(teamResp.Users) > 0 {
		usersSet, err := utils.RestoreUserInputFormatForUsers(ctx, apiManager, teamResp.Users, state.Users)
		if err != nil {
			return fmt.Errorf("failed to convert users to original format: %w", err)
		}
		state.Users = usersSet
	} else {
		state.Users = types.SetNull(types.StringType)
	}

	// Node: convert to name (API always returns node name, but user may provide name or ID)
	// Always update from API response to detect external changes

	// TODO: Update state with node name/ID what user provided based on that,
	state.Node = types.StringValue(utils.ExtractNodeName(teamResp.Node))

	return nil
}
