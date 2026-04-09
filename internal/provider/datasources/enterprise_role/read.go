// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"
	"fmt"
	"strconv"

	enterpriseRoleResource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_role"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *EnterpriseRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EnterpriseRoleDataSourceModel

	// Get configuration data
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR,
			err.Error(),
		)
		return
	}

	if err := utils.RunWithManagedCompanyContext(ctx, d.ApiManager, data.ManagedCompany, func() error {
		roleInfo, err := utils.FetchEnterpriseRoleByNameOrId(ctx, d.ApiManager, data.Role.ValueString())
		if err != nil {
			return err
		}
		if roleInfo == nil {
			return fmt.Errorf("enterprise role: '%s' not found", data.Role.ValueString())
		}

		data.Id = types.StringValue(strconv.Itoa(roleInfo.RoleId))
		data.Name = types.StringValue(roleInfo.Name)

		users, usersDiags := types.SetValueFrom(ctx, types.StringType, roleInfo.Users)
		if usersDiags.HasError() {
			return fmt.Errorf("failed to create users set: %v", usersDiags.Errors())
		}
		data.Users = users

		teams, teamsDiags := types.SetValueFrom(ctx, types.StringType, roleInfo.Teams)
		if teamsDiags.HasError() {
			return fmt.Errorf("failed to create teams set: %v", teamsDiags.Errors())
		}
		data.Teams = teams

		managingNodesMap, err := utils.MapManagedNodesPermissionsToState(ctx, roleInfo.ManagedNodesPermissions)
		if err != nil {
			return fmt.Errorf("failed to map managing nodes: %w", err)
		}
		data.ManagingNodes = managingNodesMap

		enforcementPoliciesMap, err := utils.MapEnforcementsToState(roleInfo.Enforcements, enterpriseRoleResource.GeneratedPasswordComplexity)
		if err != nil {
			return fmt.Errorf("failed to map enforcement policies: %w", err)
		}
		data.EnforcementPolicies = enforcementPoliciesMap

		data.ManagedCompany = types.StringNull()
		return nil
	}, "Read Enterprise Role Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
