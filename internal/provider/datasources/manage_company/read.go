package managecompany

import (
	"context"
	"fmt"
	"math/big"
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

	// Validate that at least one of id or name is provided
	if data.Id.IsNull() && data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"Either `id` or `name` must be provided to query the managed company.",
		)
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

	// Build command to get all companies info
	command := "msp-info --format json -v"

	apiResp, err := d.apiManager.ExecuteCommand(ctx, command, "Unable to retrieve managed company information")
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Managed Company Failed",
			err.Error(),
		)
		return
	}

	// Parse the JSON response
	var companies []utils.ManageCompanyResponse

	if err := utils.UnmarshalApiResponse(apiResp.Data, &companies); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Manage Company API Response",
			fmt.Sprintf("Unable to process the response from Keeper Commander API: %s", err.Error()),
		)
		return
	}

	// Find the company matching the provided ID or name
	var companyInfo *utils.ManageCompanyResponse

	for i := range companies {
		matched := false

		// Match by ID if provided
		if !data.Id.IsNull() {
			requestedId, _ := data.Id.ValueBigFloat().Int64()
			if companies[i].CompanyId == int(requestedId) {
				matched = true
			}
		}

		// Match by name if provided (and not already matched by ID)
		if !matched && !data.Name.IsNull() {
			if companies[i].CompanyName == data.Name.ValueString() {
				matched = true
			}
		}

		if matched {
			companyInfo = &companies[i]
			break
		}
	}

	if companyInfo == nil {
		searchCriteria := "with provided criteria"
		if !data.Id.IsNull() {
			idValue, _ := data.Id.ValueBigFloat().Int64()
			searchCriteria = fmt.Sprintf("with ID %d", idValue)
		} else if !data.Name.IsNull() {
			searchCriteria = fmt.Sprintf("with name '%s'", data.Name.ValueString())
		}
		resp.Diagnostics.AddError(
			"Managed Company Not Found",
			fmt.Sprintf("No managed company found %s. Please verify the ID or name and try again.", searchCriteria),
		)
		return
	}

	// Map the response to the model
	data.Id = types.NumberValue(big.NewFloat(float64(companyInfo.CompanyId)))
	data.Name = types.StringValue(companyInfo.CompanyName)
	data.Node = types.StringValue(companyInfo.Node)
	data.Plan = types.StringValue(companyInfo.Plan)

	// Convert storage format: "100GB" -> "100gb", "1TB" -> "1tb", "10TB" -> "10tb"
	storageLower := strings.ToLower(companyInfo.Storage)
	data.FilePlan = types.StringValue(storageLower)

	// Set the data
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
