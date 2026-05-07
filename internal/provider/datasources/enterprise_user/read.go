// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *EnterpriseUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EnterpriseUserDataSourceModel

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
		userInfo, err := utils.FetchEnterpriseUserByEmailOrId(ctx, d.ApiManager, data.User.ValueString())
		if err != nil {
			return err
		}
		if userInfo == nil {
			return fmt.Errorf("enterprise user: '%s' not found", data.User.ValueString())
		}
		data.Id = types.StringValue(strconv.Itoa(userInfo.UserId))
		data.Name = types.StringValue(userInfo.Name)
		data.Email = types.StringValue(userInfo.Email)
		data.JobTitle = types.StringValue(userInfo.JobTitle)
		roles, rolesDiags := types.SetValueFrom(ctx, types.StringType, userInfo.Roles)
		if rolesDiags.HasError() {
			return fmt.Errorf("failed to create roles set: %v", rolesDiags.Errors())
		}
		data.Roles = roles
		teams, teamsDiags := types.SetValueFrom(ctx, types.StringType, userInfo.Teams)
		if teamsDiags.HasError() {
			return fmt.Errorf("failed to create teams set: %v", teamsDiags.Errors())
		}
		data.Teams = teams
		data.ManagedCompany = types.StringNull()
		return nil
	}, "Read Enterprise User Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
