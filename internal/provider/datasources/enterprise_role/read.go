// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"
	"fmt"
	"strconv"

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
	if err := d.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	err := utils.ExecuteWithManagedCompanyContext(ctx, d.apiManager, data.ManagedCompany, func() error {
		roleInfo, err := utils.FetchEnterpriseRoleByNameOrId(ctx, d.apiManager, data.Role.ValueString())

		if err != nil {
			return err
		}
		if roleInfo == nil {
			return fmt.Errorf("Enterprise role: '%s' not found", data.Role.ValueString())
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

		data.ManagedCompany = types.StringNull() // Bec we dotn want to return the managed company in the result

		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Enterprise Role Failed",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
