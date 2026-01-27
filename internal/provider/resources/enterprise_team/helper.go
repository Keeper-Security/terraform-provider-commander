// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TeamInfo represents a team from the API response
type TeamInfo struct {
	TeamUid string `json:"team_uid"`
	Name    string `json:"name"`
}

// EnterpriseTeamReadResponse represents the team information from the read API response
type EnterpriseTeamReadResponse struct {
	TeamUid   string   `json:"team_uid"`
	Name      string   `json:"name"`
	Restricts string   `json:"restricts"`
	Node      string   `json:"node"`
	Users     []string `json:"users"`
	Roles     []string `json:"roles"`
}

type NodeInfo struct {
	NodeId int    `json:"node_id"`
	Name   string `json:"name"`
}

// parseTeamsResponse parses the JSON response from enterprise-info -t command
func parseTeamsResponse(data interface{}) ([]TeamInfo, error) {
	var teams []TeamInfo

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to process the response from Keeper Commander Service Mode API: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &teams); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise teams list from Service Mode API response: %w", err)
	}

	return teams, nil
}

// findTeamUidByName finds a team UID by name from the teams list
func findTeamUidByName(teams []TeamInfo, teamName string) (string, error) {
	for _, team := range teams {
		if team.Name == teamName {
			return team.TeamUid, nil
		}
	}
	return "", fmt.Errorf("enterprise team with name '%s' not found in the response", teamName)
}

// fetchTeamUidByName fetches the team UID by name using the API
func fetchTeamUidByName(ctx context.Context, apiManager *api.ApiManager, teamName string) (string, error) {
	teamsResp, err := apiManager.ExecuteCommand(ctx, fmt.Sprintf("enterprise-info -t --format json -q '%s'", teamName), "Unable to fetch enterprise team ID")
	if err != nil {
		return "", err
	}

	teams, err := parseTeamsResponse(teamsResp.Data)
	if err != nil {
		return "", err
	}

	return findTeamUidByName(teams, teamName)
}

func parseNodesResponse(data interface{}) ([]NodeInfo, error) {
	var nodes []NodeInfo

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to process the response from Keeper Commander Service Mode API: %w", err)
	}

	if err := json.Unmarshal(dataBytes, &nodes); err != nil {
		return nil, fmt.Errorf("unable to parse enterprise nodes list from Service Mode API response: %w", err)
	}

	return nodes, nil
}

func buildEnterpriseTeamAddCommand(data EnterpriseTeamResourceModel) string {
	var parts []string

	parts = append(parts, "enterprise-team")

	// Required parameters
	parts = append(parts, fmt.Sprintf("--add '%s'", data.Name.ValueString()))

	// Optional parameters
	if !data.RestrictEdit.IsNull() && data.RestrictEdit.ValueBool() {
		parts = append(parts, "--restrict-edit on")
	}

	if !data.RestrictShare.IsNull() && data.RestrictShare.ValueBool() {
		parts = append(parts, "--restrict-share on")
	}

	if !data.RestrictView.IsNull() && data.RestrictView.ValueBool() {
		parts = append(parts, "--restrict-view on")
	}

	if !data.Node.IsNull() {
		parts = append(parts, fmt.Sprintf("--node '%s'", data.Node.ValueString()))
	}

	return strings.Join(parts, " ")
}

func buildEnterpriseTeamUpdateCommand(ctx context.Context, apiManager *api.ApiManager, plan *EnterpriseTeamResourceModel, state *EnterpriseTeamResourceModel) (string, error) {
	var parts []string

	parts = append(parts, "enterprise-team")

	// Required parameters
	if !state.Name.Equal(plan.Name) {
		parts = append(parts, fmt.Sprintf("--name '%s'", plan.Name.ValueString()))
	}

	// Optional parameters
	if !state.RestrictEdit.Equal(plan.RestrictEdit) {
		if !plan.RestrictEdit.IsNull() && plan.RestrictEdit.ValueBool() {
			parts = append(parts, "--restrict-edit on")
		} else {
			parts = append(parts, "--restrict-edit off")
		}
	}

	if !state.RestrictShare.Equal(plan.RestrictShare) {
		if !plan.RestrictShare.IsNull() && plan.RestrictShare.ValueBool() {
			parts = append(parts, "--restrict-share on")
		} else {
			parts = append(parts, "--restrict-share off")
		}
	}

	if !state.RestrictView.Equal(plan.RestrictView) {
		if !plan.RestrictView.IsNull() && plan.RestrictView.ValueBool() {
			parts = append(parts, "--restrict-view on")
		} else {
			parts = append(parts, "--restrict-view off")
		}
	}

	// TODO: we will node with its id all time
	if !state.Node.Equal(plan.Node) {
		parts = append(parts, fmt.Sprintf("--node '%s'", plan.Node.ValueString()))
	}

	// Process users and roles changes
	if !state.Users.Equal(plan.Users) {
		users, err := utils.FetchAndProcessUsers(ctx, apiManager, state.Users, plan.Users)
		if err != nil {
			return "", err
		}
		if users != "" {
			parts = append(parts, users)
		}
	}

	if !state.Roles.Equal(plan.Roles) {
		roles, err := utils.FetchAndProcessRoles(ctx, apiManager, state.Roles, plan.Roles)
		if err != nil {
			return "", err
		}
		if roles != "" {
			parts = append(parts, roles)
		}
	}

	parts = append(parts, fmt.Sprintf("'%s'", state.Id.ValueString()))

	return strings.Join(parts, " "), nil
}

// parseEnterpriseTeamReadResponse parses the JSON response from enterprise-info -t command
func parseEnterpriseTeamReadResponse(data interface{}) ([]EnterpriseTeamReadResponse, error) {
	var teams []EnterpriseTeamReadResponse

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

// convertNodeToName converts node from API response to node name
// API may return node as name or ID. If it's an ID, we fetch nodes and convert to name.
// Final state always stores node name (not ID).
func convertNodeToName(ctx context.Context, apiManager *api.ApiManager, nodeFromApi string) (types.String, error) {
	// Handle empty/null cases
	nodeFromApi = strings.TrimSpace(nodeFromApi)
	if nodeFromApi == "" {
		return types.StringNull(), nil
	}

	// Check if nodeFromApi is numeric (likely an ID)
	// If it's not numeric, assume it's already a name and use it directly
	nodeIdInt, err := strconv.Atoi(nodeFromApi)
	if err != nil {
		// Not numeric - assume it's a name, use it directly
		return types.StringValue(nodeFromApi), nil
	}

	// It's numeric - need to convert ID to name
	// Fetch all nodes to build lookup
	nodesResp, err := apiManager.ExecuteCommand(ctx, "enterprise-info -n --format json -v -q", "Unable to fetch enterprise nodes")
	if err != nil {
		return types.StringNull(), fmt.Errorf("failed to fetch nodes: %w", err)
	}

	nodes, err := parseNodesResponse(nodesResp.Data)
	if err != nil {
		return types.StringNull(), fmt.Errorf("failed to parse nodes: %w", err)
	}

	// Find node by ID
	for _, node := range nodes {
		if node.NodeId == nodeIdInt {
			if node.Name != "" {
				nodeName, ok := utils.ExtractNodeIDFromCreateNodeResponse(node.Name)
				if ok {
					return types.StringValue(nodeName), nil
				}
				return types.StringValue(node.Name), nil
			}
		}
	}

	// Node ID not found - return null (node may have been deleted)
	// This handles edge case where node was deleted but team still references it
	return types.StringNull(), nil
}

// mapTeamReadResponseToModel maps the API response to the resource model
func mapTeamReadResponseToModel(ctx context.Context, apiManager *api.ApiManager, teamResp EnterpriseTeamReadResponse, state *EnterpriseTeamResourceModel) error {
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

	// TODO: we are not getting node in respose
	// Update state with node name from API (always a name, not ID)
	state.Node = types.StringValue(utils.ExtractNodeName(teamResp.Node))

	return nil
}
