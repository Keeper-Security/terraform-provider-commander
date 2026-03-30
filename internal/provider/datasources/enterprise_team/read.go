// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *EnterpriseTeamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EnterpriseTeamDataSourceModel

	// Get configuration data
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	if err := utils.RunWithManagedCompanyContext(ctx, d.ApiManager, data.ManagedCompany, func() error {
		teamInfo, err := utils.FetchEnterpriseTeamByNameOrId(ctx, d.ApiManager, data.Team.ValueString())
		if err != nil {
			return err
		}
		if teamInfo == nil {
			return fmt.Errorf("enterprise team: '%s' not found", data.Team.ValueString())
		}
		data.Id = types.StringValue(teamInfo.TeamUid)
		data.Name = types.StringValue(teamInfo.Name)
		users, usersDiags := types.SetValueFrom(ctx, types.StringType, teamInfo.Users)
		if usersDiags.HasError() {
			return fmt.Errorf("failed to create users set: %v", usersDiags.Errors())
		}
		data.Users = users
		roles, rolesDiags := types.SetValueFrom(ctx, types.StringType, teamInfo.Roles)
		if rolesDiags.HasError() {
			return fmt.Errorf("failed to create roles set: %v", rolesDiags.Errors())
		}
		data.Roles = roles
		data.ManagedCompany = types.StringNull()
		return nil
	}, "Read Enterprise Team Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
