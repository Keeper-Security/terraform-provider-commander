// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *EnterpriseScimDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EnterpriseScimDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR,
			err.Error(),
		)
		return
	}

	if err := utils.RunWithManagedCompanyContext(ctx, d.ApiManager, data.ManagedCompany, func() error {
		scimInfo, err := utils.FetchEnterpriseScimById(ctx, d.ApiManager, data.Scim.ValueString())
		if err != nil {
			return err
		}
		if scimInfo == nil {
			return fmt.Errorf("enterprise SCIM configuration: '%s' not found", data.Scim.ValueString())
		}

		data.Scim = types.StringValue(data.Scim.ValueString())
		data.ScimId = types.StringValue(strconv.Itoa(scimInfo.ScimID))
		data.ScimUrl = types.StringValue(scimInfo.ScimURL)
		data.NodeId = types.StringValue(strconv.Itoa(scimInfo.NodeID))
		data.NodeName = types.StringValue(utils.ExtractNodeName(scimInfo.NodeName))
		data.Status = types.StringValue(scimInfo.Status)

		if scimInfo.Prefix == "" {
			data.Prefix = types.StringNull()
		} else {
			data.Prefix = types.StringValue(scimInfo.Prefix)
		}

		data.UniqueGroups = types.BoolValue(scimInfo.UniqueGroups)
		data.ManagedCompany = types.StringNull()
		return nil
	}, "Read Enterprise SCIM Configuration Failed", &resp.Diagnostics); err != nil && errors.Is(err, utils.ErrResourceRemoved) {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
