package enterpriserole

import (
	"context"
	"fmt"
	"strconv"

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
		var roles []RoleResponse

		// Unmarshal API response into roles struct
		if err := utils.UnmarshalApiResponse(apiResp.Data, &roles); err != nil {
			return fmt.Errorf("unable to parse enterprise roles list from API response: %w", err)
		}

		// Find the role matching state.Id
		var roleInfo *RoleResponse
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
			return nil
		}

		// Map the response to the state
		state.Id = types.StringValue(stateId)
		state.Name = types.StringValue(roleInfo.Name)
		state.Node = types.StringValue(utils.ExtractNodeName(roleInfo.Node))

		// TODO: Later we will add users, teams, enforcement policies, managing nodes to state

		// Set the updated state
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Read Enterprise Role Failed",
			err.Error(),
		)
		return
	}
}
