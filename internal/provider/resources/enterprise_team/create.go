package enterpriseteam

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseTeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EnterpriseTeamResourceModel

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
	err := utils.ExecuteWithManagedCompanyContext(ctx, r.apiManager, data.ManagedCompany, func() error {
		// Fetch and process users/teams before creating the role
		// For create, stateUsers and stateRoles are null/empty, only planUsers and planRoles have items to add
		users, err := utils.FetchAndProcessUsers(ctx, r.apiManager, types.SetNull(types.StringType), data.Users)
		if err != nil {
			return err
		}

		roles, err := utils.FetchAndProcessRoles(ctx, r.apiManager, types.SetNull(types.StringType), data.Roles)
		if err != nil {
			return err
		}

		if err := addTeamBasicAttributes(ctx, r.apiManager, data); err != nil {
			return err
		}

		// TODO: WE WILL REMOVE THIS fetchTeamUidByName FUNCTION AFTER WE ARE RECIEVING TEAM UID IN THE RESPONSE IN COMMANDER CLI
		// Fetch the team UID by name
		teamUid, err := fetchTeamUidByName(ctx, r.apiManager, data.Name.ValueString())
		if err != nil {
			return err
		}

		data.Id = types.StringValue(teamUid)

		// Combine users and roles flags
		if users != "" {
			// Add Users and Roles to the recently created team
			command := fmt.Sprintf("enterprise-team '%s' %s -v", teamUid, users)

			_, err = r.apiManager.ExecuteCommand(ctx, command, "Unable to add users to the enterprise team")
			if err != nil {
				return err
			}
		}
		if roles != "" {
			// Add Users and Roles to the recently created team
			command := fmt.Sprintf("enterprise-team '%s' %s -v", teamUid, roles)

			_, err = r.apiManager.ExecuteCommand(ctx, command, "Unable to add roles to the enterprise team")
			if err != nil {
				return err
			}
		}

		return nil

	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Create Enterprise Team Failed",
			err.Error(),
		)
		return
	}

	// Set the ID in the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func addTeamBasicAttributes(ctx context.Context, apiManager *api.ApiManager, data EnterpriseTeamResourceModel) error {
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

	command := strings.Join(parts, " ")

	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to create enterprise team")
	if err != nil {
		return err
	}

	return nil
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
