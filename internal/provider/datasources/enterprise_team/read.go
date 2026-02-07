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
	if err := d.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	err := utils.ExecuteWithManagedCompanyContext(ctx, d.apiManager, data.ManagedCompany, func() error {
		teamInfo, err := utils.FetchEnterpriseTeamByNameOrId(ctx, d.apiManager, data.Team.ValueString())
		if err != nil {
			return err
		}
		if teamInfo == nil {
			return fmt.Errorf("Enterprise team: '%s' not found", data.Team.ValueString())
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

		data.ManagedCompany = types.StringNull() // Bec we dotn want to return the managed company in the result

		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Enterprise Team Failed",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
