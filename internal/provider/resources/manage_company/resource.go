// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ManageCompanyResource{}
var _ resource.ResourceWithConfigure = &ManageCompanyResource{}

type ManageCompanyResource struct {
	apiManager *api.ApiManager
}

type ManageCompanyResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Node     types.String `tfsdk:"node"`
	Plan     types.String `tfsdk:"plan"`
	Seats    types.Int64  `tfsdk:"seats"`
	FilePlan types.String `tfsdk:"file_plan"`
	AddOns   types.Set    `tfsdk:"add_ons"`
}

func (r *ManageCompanyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_manage_company"
}

func (r *ManageCompanyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					nameValidator{},
				},
				Description:         "Managed Company Name.",
				MarkdownDescription: "Managed Company Name.",
			},
			"node": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					nodeValidator{},
				},
				Description:         "Managing Node name or ID.",
				MarkdownDescription: "Managing Node name or ID.",
			},
			"seats": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					seatsValidator{},
				},
				Description:         "Maximum Licenses Allowed.",
				MarkdownDescription: "Maximum Licenses Allowed.",
			},
			"plan": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					planValidator{},
				},
				Description:         "Base plan. Must be one of: " + strings.Join(PlanOptions, ", "),
				MarkdownDescription: "Base plan. Must be one of: `" + strings.Join(PlanOptions, "`, `") + "`",
			},
			"file_plan": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					filePlanValidator{},
				},
				Description:         "File storage plan. Must be one of: " + strings.Join(FilePlanOptions, ", "),
				MarkdownDescription: "File storage plan. Must be one of: `" + strings.Join(FilePlanOptions, "`, `") + "`",
			},
			"add_ons": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					addOnsValidator{},
				},
				Description:         "Secure Add-Ons to apply to the Managed Company. Must be one of: " + strings.Join(GetAllValidAddOns(), ", "),
				MarkdownDescription: "Secure Add-Ons to apply to the Managed Company. Must be one of: `" + strings.Join(GetAllValidAddOns(), "`, `") + "`",
			},
		},
	}
}

func (r *ManageCompanyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.apiManager = apiManager
}

// ensureApiManager validates that apiManager is configured and returns an error if not
func (r *ManageCompanyResource) ensureApiManager() error {
	if r.apiManager == nil {
		return fmt.Errorf("the Keeper Commander provider is not properly configured. Please ensure the provider is set up with valid service_mode_url and service_mode_api_key")
	}
	return nil
}

func (r *ManageCompanyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data ManageCompanyResourceModel

	// Get planned data from Terraform
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Build the Commander command string
	command := buildManageCompanyAddCommand(data)

	// Only switch to MSP if it's an MSP account
	if r.apiManager.IsMspAccount {
		if err := utils.SwitchToMsp(ctx, r.apiManager); err != nil {
			resp.Diagnostics.AddError(
				"Create Managed Company Failed",
				fmt.Sprintf("Failed to switch to MSP context: %s", err.Error()),
			)
			return
		}
	}

	apiResp, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to create managed company")
	if err != nil {
		resp.Diagnostics.AddError(
			"Create Managed Company Failed",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(apiResp.Message.String())

	// Set the ID in the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}
func (r *ManageCompanyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ManageCompanyResourceModel

	// Get current state (old values)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Only switch to MSP if it's an MSP account
	if r.apiManager.IsMspAccount {
		if err := utils.SwitchToMsp(ctx, r.apiManager); err != nil {
			resp.Diagnostics.AddError(
				"Read Managed Company Failed",
				fmt.Sprintf("Failed to switch to MSP context: %s", err.Error()),
			)
			return
		}
	}

	// Execute msp-down command (setup/initialization)
	_, err := r.apiManager.ExecuteCommand(ctx, "msp-down", "Unable to initialize managed company service")
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Managed Company Failed",
			err.Error(),
		)
		return
	}

	// Build command to get all companies info
	command := fmt.Sprintf("msp-info -m '%s' --format json -v", state.Id.ValueString())

	apiResp, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to retrieve managed company information")
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Managed Company Failed",
			err.Error(),
		)
		return
	}

	// Parse the JSON response - it's an array of company objects
	var companies []struct {
		CompanyId   int      `json:"company_id"`
		CompanyName string   `json:"company_name"`
		Node        string   `json:"node"`
		NodeName    string   `json:"node_name"`
		Plan        string   `json:"plan"`
		Storage     string   `json:"storage"`
		Addons      []string `json:"addons"`
		Allocated   int      `json:"allocated"`
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

	// Find the company matching state.Id
	stateIdStr := state.Id.ValueString()
	stateIdInt, err := strconv.Atoi(stateIdStr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Company ID",
			fmt.Sprintf("The company ID in the state is invalid: %s", err.Error()),
		)
		return
	}

	var companyInfo *struct {
		CompanyId   int      `json:"company_id"`
		CompanyName string   `json:"company_name"`
		Node        string   `json:"node"`
		NodeName    string   `json:"node_name"`
		Plan        string   `json:"plan"`
		Storage     string   `json:"storage"`
		Addons      []string `json:"addons"`
		Allocated   int      `json:"allocated"`
	}

	for i := range companies {
		if companies[i].CompanyId == stateIdInt {
			companyInfo = &companies[i]
			break
		}
	}

	if companyInfo == nil {
		// Resource not found - remove from state
		// Terraform will detect this and mark the resource for destruction
		resp.State.RemoveResource(ctx)
		return
	}

	// Map the response to the model
	state.Id = types.StringValue(strconv.Itoa(companyInfo.CompanyId))
	state.Name = types.StringValue(companyInfo.CompanyName)

	state.Node = types.StringValue(utils.ExtractNodeName(companyInfo.NodeName))

	state.Plan = types.StringValue(companyInfo.Plan)
	state.Seats = types.Int64Value(int64(companyInfo.Allocated))

	// Convert storage format: "100GB" -> "100gb", "1TB" -> "1tb", "10TB" -> "10tb"
	storageLower := strings.ToLower(companyInfo.Storage)
	state.FilePlan = types.StringValue(storageLower)

	// Convert addons array of strings to types.Set
	addOnsElements := make([]types.String, len(companyInfo.Addons))
	for i, addon := range companyInfo.Addons {
		addOnsElements[i] = types.StringValue(addon)
	}
	addOnsSet, diags := types.SetValueFrom(ctx, types.StringType, addOnsElements)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	state.AddOns = addOnsSet

	// Set the updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (r *ManageCompanyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ManageCompanyResourceModel
	var state ManageCompanyResourceModel

	// Get planned data (new values)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state (old values)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	command := buildManageCompanyUpdateCommand(&plan, &state)

	// Only switch to MSP if it's an MSP account
	if r.apiManager.IsMspAccount {
		if err := utils.SwitchToMsp(ctx, r.apiManager); err != nil {
			resp.Diagnostics.AddError(
				"Update Managed Company Failed",
				fmt.Sprintf("Failed to switch to MSP context: %s", err.Error()),
			)
			return
		}
	}

	_, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to update managed company")
	if err != nil {
		resp.Diagnostics.AddError(
			"Update Managed Company Failed",
			err.Error(),
		)
		return
	}

	// Keep the same ID
	plan.Id = state.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
func (r *ManageCompanyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state ManageCompanyResourceModel

	// Get planned data from Terraform
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.ensureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	// Build delete command
	command := fmt.Sprintf("msp-remove '%s' -f", state.Id.ValueString())

	// Only switch to MSP if it's an MSP account
	if r.apiManager.IsMspAccount {
		if err := utils.SwitchToMsp(ctx, r.apiManager); err != nil {
			resp.Diagnostics.AddError(
				"Delete Managed Company Failed",
				fmt.Sprintf("Failed to switch to MSP context: %s", err.Error()),
			)
			return
		}
	}

	_, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to delete managed company")
	if err != nil {
		resp.Diagnostics.AddError(
			"Delete Managed Company Failed",
			err.Error(),
		)
		return
	}
}

func NewManageCompanyResource() resource.Resource {
	return &ManageCompanyResource{}
}
