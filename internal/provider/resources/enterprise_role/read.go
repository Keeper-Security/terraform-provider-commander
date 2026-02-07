// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
		roleInfo, err := utils.FetchEnterpriseRoleByNameOrId(ctx, r.apiManager, state.Id.ValueString())
		if err != nil {
			return err
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
	managingNodesMap, err := mapManagedNodesPermissionsToState(ctx, roleInfo.ManagedNodesPermissions)
	if err != nil {
		return fmt.Errorf("failed to map managing nodes to state: %w", err)
	}
	state.ManagingNodes = managingNodesMap

	// Enforcement policies: map API enforcements (keys only, lowercase) to state map
	// Use value from state if key exists there, else empty string
	enforcementPoliciesMap, err := mapEnforcementsToState(roleInfo.Enforcements, state.EnforcementPolicies)
	if err != nil {
		return fmt.Errorf("failed to map enforcement policies to state: %w", err)
	}
	state.EnforcementPolicies = enforcementPoliciesMap

	return nil
}

// managingNodesMapElemType is the object type for each entry in the managing_nodes map.
var managingNodesMapElemType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"privileges": types.SetType{ElemType: types.StringType},
		"cascade":    types.BoolType,
	},
}

// mapManagedNodesPermissionsToState converts API ManagedNodesPermissions to a types.Map for state.
// Map key is the node name (ExtractNodeName). Value is object with privileges (Set) and cascade (Bool).
func mapManagedNodesPermissionsToState(_ context.Context, perms []utils.ManagedNodePermission) (types.Map, error) {
	if len(perms) == 0 {
		return types.MapNull(managingNodesMapElemType), nil
	}

	elements := make(map[string]attr.Value)
	for _, p := range perms {
		key := p.NodeName

		//  this is edge case, we will get node name all time
		if key == "" {
			key = strconv.FormatInt(p.NodeId, 10)
		}

		privilegeElems := make([]attr.Value, len(p.Privileges))
		for i, pr := range p.Privileges {
			privilegeElems[i] = types.StringValue(pr)
		}

		privilegesSet := types.SetValueMust(types.StringType, privilegeElems)
		obj := types.ObjectValueMust(
			managingNodesMapElemType.AttrTypes,
			map[string]attr.Value{
				"privileges": privilegesSet,
				"cascade":    types.BoolValue(p.Cascade),
			},
		)
		elements[key] = obj
	}

	mapVal, diags := types.MapValue(managingNodesMapElemType, elements)
	if diags.HasError() {
		return types.MapNull(managingNodesMapElemType), fmt.Errorf("failed to build managing_nodes map: %v", diags)
	}
	return mapVal, nil
}

// mapEnforcementsToState converts API enforcements (array of keys in lowercase) to a types.Map for state.
// API returns only keys; we normalize to UPPER_SNAKE_CASE. If the key exists in state, use state's value; else use "".
func mapEnforcementsToState(enforcements []string, stateEnforcementPolicies types.Map) (types.Map, error) {
	if len(enforcements) == 0 {
		return types.MapNull(types.StringType), nil
	}
	stateValues := make(map[string]string)
	if !stateEnforcementPolicies.IsNull() && !stateEnforcementPolicies.IsUnknown() {
		for key, val := range stateEnforcementPolicies.Elements() {
			if s, ok := val.(types.String); ok && !s.IsNull() && !s.IsUnknown() {
				stateValues[key] = s.ValueString()
			}
		}
	}
	elements := make(map[string]attr.Value)
	for _, key := range enforcements {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		normalizedKey := strings.ToUpper(key)
		if v, ok := stateValues[normalizedKey]; ok {
			elements[normalizedKey] = types.StringValue(v)
		} else {
			elements[normalizedKey] = types.StringValue("")
		}
	}
	if len(elements) == 0 {
		return types.MapNull(types.StringType), nil
	}
	mapVal, diags := types.MapValue(types.StringType, elements)
	if diags.HasError() {
		return types.MapNull(types.StringType), fmt.Errorf("failed to build enforcement_policies map: %v", diags)
	}
	return mapVal, nil
}
