// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *ManageCompanyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ManageCompanyDataSourceModel

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

	// For MSP accounts, ensure we are in MSP context before running msp-info
	if d.apiManager.IsMspAccount {
		if err := utils.SwitchToMsp(ctx, d.apiManager); err != nil {
			resp.Diagnostics.AddError(
				"Read Managed Company Failed",
				fmt.Sprintf("Failed to switch to MSP context: %s", err.Error()),
			)
			return
		}
	}

	if err := utils.MspDown(ctx, d.apiManager); err != nil {
		resp.Diagnostics.AddError(
			"Read Managed Company Failed",
			err.Error(),
		)
		return
	}

	companyInfo, err := utils.FetchManageCompanyByNameOrId(ctx, d.apiManager, data.ManagedCompany.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Managed Company Failed",
			err.Error(),
		)
		return
	}

	if companyInfo == nil {
		resp.Diagnostics.AddError(
			"Managed Company Not Found",
			fmt.Sprintf("managed company: '%s' not found", data.ManagedCompany.ValueString()),
		)
		return
	}

	// Map the response to the model
	data.Id = types.StringValue(strconv.Itoa(companyInfo.CompanyId))
	data.Name = types.StringValue(companyInfo.CompanyName)
	data.Node = types.StringValue(companyInfo.Node)
	data.Plan = types.StringValue(companyInfo.Plan)

	// Convert storage format: "100GB" -> "100gb", "1TB" -> "1tb", "10TB" -> "10tb"
	storageLower := strings.ToLower(companyInfo.Storage)
	data.FilePlan = types.StringValue(storageLower)

	// Set the data
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
