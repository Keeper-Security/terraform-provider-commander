// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TeamInfo represents a team from the API response
type TeamInfo struct {
	TeamUid string `json:"team_uid"`
	Name    string `json:"name"`
}

// ParseTeamsResponse parses the JSON response from enterprise-info -t command
func ParseTeamsResponse(data interface{}) ([]TeamInfo, error) {
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

// BuildTeamLookupMaps creates lookup maps from API response
func BuildTeamLookupMaps(teamsRespData []TeamInfo) LookupMaps {
	identifierToId := make(map[string]string)
	idToIdentifier := make(map[string]string)

	for _, team := range teamsRespData {
		if team.TeamUid != "" && team.Name != "" {
			identifierToId[team.Name] = team.TeamUid
			idToIdentifier[team.TeamUid] = team.Name
		}
	}

	return LookupMaps{
		IdentifierToId: identifierToId,
		IdToIdentifier: idToIdentifier,
	}
}

// ConvertTeamsToIdMap converts a types.Set of teams to a map of team_uid -> original input
func ConvertTeamsToIdMap(teams types.Set, lookup LookupMaps, teamsRespData []TeamInfo) (map[string]string, error) {
	validateTeam := func(userInput string) (bool, string) {
		for _, team := range teamsRespData {
			if team.Name == userInput && team.TeamUid == "" {
				return false, "team '" + userInput + "' exists but has no valid team_uid. This team cannot be used"
			}
		}
		return true, ""
	}

	return ConvertItemsToIdMap(teams, lookup, "team", validateTeam)
}

// FetchAndProcessTeams processes teams for both create and update operations
// For create: stateTeams should be null/empty, planTeams contains teams to add
// For update: compares stateTeams (old) with planTeams (new) to determine additions and removals
// Returns a string with -at "team_uid" for additions and -rt "team_uid" for removals
func FetchAndProcessTeams(ctx context.Context, apiManager *api.ApiManager, stateTeams types.Set, planTeams types.Set) (string, error) {
	// Early return if both are empty/null
	if (stateTeams.IsNull() || len(stateTeams.Elements()) == 0) &&
		(planTeams.IsNull() || len(planTeams.Elements()) == 0) {
		return "", nil
	}

	// Fetch teams from API
	teamsResp, err := apiManager.ExecuteCommand(ctx, "enterprise-info -t --format json", "Unable to fetch enterprise teams")
	if err != nil {
		return "", err
	}

	// Parse the teams response
	teamsRespData, err := ParseTeamsResponse(teamsResp.Data)
	if err != nil {
		return "", err
	}

	// Build lookup maps
	lookup := BuildTeamLookupMaps(teamsRespData)

	// Convert state teams to team_uid map (old teams)
	stateTeamIdMap, err := ConvertTeamsToIdMap(stateTeams, lookup, teamsRespData)
	if err != nil {
		return "", err
	}

	// Convert plan teams to team_uid map (new teams)
	planTeamIdMap, err := ConvertTeamsToIdMap(planTeams, lookup, teamsRespData)
	if err != nil {
		return "", err
	}

	// Early return if no changes
	if len(stateTeamIdMap) == 0 && len(planTeamIdMap) == 0 {
		return "", nil
	}

	// Find teams to add and remove
	var parts []string

	// Add teams that are in plan but not in state
	for teamId := range planTeamIdMap {
		if _, exists := stateTeamIdMap[teamId]; !exists {
			parts = append(parts, fmt.Sprintf("-at '%s'", teamId))
		}
	}

	// Remove teams that are in state but not in plan
	for teamId := range stateTeamIdMap {
		if _, exists := planTeamIdMap[teamId]; !exists {
			parts = append(parts, fmt.Sprintf("-rt '%s'", teamId))
		}
	}

	if len(parts) == 0 {
		return "", nil
	}

	return strings.Join(parts, " "), nil
}

// RestoreUserInputFormatForTeams converts team names from API response back to the format
// that the user originally provided in their Terraform configuration.
//
// This function preserves the original user input format to prevent false diffs in Terraform plans.
// If a user specified teams by team_uid (e.g., "abc123xyz"), the function will return team_uids.
// If they specified by name (e.g., "Dev Team"), it will return names.
//
// Parameters:
//   - teamNames: Team names returned by the API (from enterprise-info command)
//   - currentState: Current Terraform state containing teams (what user originally provided)
//
// Returns:
//   - types.Set: Set of teams in the original user input format (names or team_uids)
//   - error: Error if fetching teams or building lookup maps fails
//
// Example:
//
//	User config: teams = ["abc123xyz", "def456uvw"]
//	API returns: ["Dev Team", "QA Team"]
//	Function returns: ["abc123xyz", "def456uvw"] (preserves original team_uids)
func RestoreUserInputFormatForTeams(ctx context.Context, apiManager *api.ApiManager, teamNames []string, currentState types.Set) (types.Set, error) {
	return RestoreUserInputFormatFromApiResponse(
		ctx,
		apiManager,
		teamNames,
		currentState,
		"team",
		"enterprise-info -t --format json",
		func(data interface{}) (interface{}, error) { return ParseTeamsResponse(data) },
		func(data interface{}) LookupMaps { return BuildTeamLookupMaps(data.([]TeamInfo)) },
	)
}
