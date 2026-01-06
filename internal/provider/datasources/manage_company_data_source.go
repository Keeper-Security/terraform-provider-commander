// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ManageCompanyDataSource{}
var _ datasource.DataSourceWithConfigure = &ManageCompanyDataSource{}

type ManageCompanyDataSource struct {
	apiManager *api.ApiManager
}

type ManageCompanyDataSourceModel struct {
	// Input fields (optional)
	Id   types.Number `tfsdk:"id"`
	Name types.String `tfsdk:"name"`

	// Output fields
	Node     types.String `tfsdk:"node"`
	Plan     types.String `tfsdk:"plan"`
	FilePlan types.String `tfsdk:"file_plan"`
}

func (d *ManageCompanyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_manage_company"
}

func (d *ManageCompanyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.NumberAttribute{
				Optional:            true,
				Description:         "Managed Company ID.",
				MarkdownDescription: "Managed Company ID.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "Managed Company Name.",
				MarkdownDescription: "Managed Company Name.",
			},
			"node": schema.StringAttribute{
				Computed:            true,
				Description:         "Managing Node name or ID.",
				MarkdownDescription: "Managing Node name or ID.",
			},
			"plan": schema.StringAttribute{
				Computed:            true,
				Description:         "Base plan.",
				MarkdownDescription: "Base plan.",
			},
			"file_plan": schema.StringAttribute{
				Computed:            true,
				Description:         "File storage plan.",
				MarkdownDescription: "File storage plan.",
			},
		},
	}
}

func (d *ManageCompanyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiManager, ok := req.ProviderData.(*api.ApiManager)
	if !ok {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			fmt.Sprintf("The provider was not configured correctly. Expected API manager, but got: %T. Please check your provider configuration.", req.ProviderData),
		)
		return
	}

	d.apiManager = apiManager
}

// ensureApiManager validates that apiManager is configured and returns an error if not
func (d *ManageCompanyDataSource) ensureApiManager() error {
	if d.apiManager == nil {
		return fmt.Errorf("the Keeper Commander provider is not properly configured. Please ensure the provider is set up with valid service_mode_url and service_mode_api_key")
	}
	return nil
}

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

	// Parse the JSON response - it's an array of company objects
	var companies []struct {
		CompanyId   int    `json:"company_id"`
		CompanyName string `json:"company_name"`
		Node        string `json:"node"`
		Plan        string `json:"plan"`
		Storage     string `json:"storage"`
		Allocated   int    `json:"allocated"`
	}

	// Convert apiResp.Data to JSON bytes and unmarshal
	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid API Response",
			fmt.Sprintf("Unable to process the response from Keeper Commander API: %s", err.Error()),
		)
		return
	}

	if err := json.Unmarshal(dataBytes, &companies); err != nil {
		resp.Diagnostics.AddError(
			"Invalid API Response",
			fmt.Sprintf("Unable to parse managed companies list from API response: %s", err.Error()),
		)
		return
	}

	// Find the company matching the provided ID or name
	var companyInfo *struct {
		CompanyId   int    `json:"company_id"`
		CompanyName string `json:"company_name"`
		Node        string `json:"node"`
		Plan        string `json:"plan"`
		Storage     string `json:"storage"`
		Allocated   int    `json:"allocated"`
	}

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

func NewManageCompanyDataSource() datasource.DataSource {
	return &ManageCompanyDataSource{}
}
