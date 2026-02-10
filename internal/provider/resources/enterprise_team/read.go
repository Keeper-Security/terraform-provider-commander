// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import (
	"context"
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
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	managedCompany := state.ManagedCompany
	err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, managedCompany, func() error {
		teamInfo, err := utils.FetchEnterpriseTeamByNameOrId(ctx, r.ApiManager, state.Id.ValueString())
		if err != nil {
			return err
		}
		if teamInfo == nil {
			resp.State.RemoveResource(ctx)
			return utils.ErrResourceRemoved
		}
		if err := mapTeamReadResponseToModel(ctx, r.ApiManager, *teamInfo, &state); err != nil {
			return fmt.Errorf("Failed to map team response to model: %w", err)
		}
		return nil
	}, "Read Enterprise Team Failed", &resp.Diagnostics)
	if err != nil && errors.Is(err, utils.ErrResourceRemoved) {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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

	// Node: preserve original format (name or ID) as user provided
	nodeVal, err := utils.RestoreUserInputFormatForNode(ctx, apiManager, teamResp.Node, state.Node)
	if err != nil {
		return fmt.Errorf("failed to convert node to original format: %w", err)
	}
	state.Node = nodeVal

	return nil
}
