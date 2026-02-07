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
	if err := d.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	err := utils.ExecuteWithManagedCompanyContext(ctx, d.apiManager, data.ManagedCompany, func() error {
		userInfo, err := utils.FetchEnterpriseUserByEmailOrId(ctx, d.apiManager, data.User.ValueString())

		if err != nil {
			return err
		}
		if userInfo == nil {
			return fmt.Errorf("Enterprise user: '%s' not found", data.User.ValueString())
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

		data.Status = types.StringValue(userInfo.Status)
		data.ManagedCompany = types.StringNull() // Bec we dotn want to return the managed company in the result

		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Enterprise User Failed",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
