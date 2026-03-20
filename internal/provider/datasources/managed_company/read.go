// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcompany

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *ManagedCompanyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ManagedCompanyDataSourceModel

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

	if err := utils.RunWithMspContext(ctx, d.ApiManager, func() error {
		if err := utils.MspDown(ctx, d.ApiManager); err != nil {
			return err
		}
		companyInfo, err := utils.FetchManagedCompanyByNameOrId(ctx, d.ApiManager, data.ManagedCompany.ValueString())
		if err != nil {
			return err
		}
		if companyInfo == nil {
			return fmt.Errorf("managed company: '%s' not found", data.ManagedCompany.ValueString())
		}
		data.Id = types.StringValue(strconv.Itoa(companyInfo.CompanyId))
		data.Name = types.StringValue(companyInfo.CompanyName)
		data.Node = types.StringValue(companyInfo.Node)
		data.NodeName = types.StringValue(companyInfo.NodeName)
		data.Plan = types.StringValue(companyInfo.Plan)

		storageLower := strings.ToLower(companyInfo.Storage)
		data.FilePlan = types.StringValue(storageLower)

		data.Seats = types.Int64Value(int64(companyInfo.Allocated))

		addOnsSet, diags := types.SetValueFrom(ctx, types.StringType, companyInfo.Addons)
		if diags.HasError() {
			return fmt.Errorf("failed to build add_ons set: %v", diags.Errors())
		}
		data.AddOns = addOnsSet

		return nil
	}, "Read Managed Company Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
