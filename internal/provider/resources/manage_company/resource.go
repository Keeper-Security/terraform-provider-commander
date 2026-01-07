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
	AddOns   types.List   `tfsdk:"add_ons"`
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
				Description:         "Base plan. Must be one of: " + strings.Join(planOptions, ", "),
				MarkdownDescription: "Base plan. Must be one of: `" + strings.Join(planOptions, "`, `") + "`",
			},
			"file_plan": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					filePlanValidator{},
				},
				Description:         "File storage plan. Must be one of: " + strings.Join(filePlanOptions, ", "),
				MarkdownDescription: "File storage plan. Must be one of: `" + strings.Join(filePlanOptions, "`, `") + "`",
			},
			"add_ons": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
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

	apiResp, err := r.apiManager.ExecuteCommand(ctx, command, "Unable to create managed company")
	if err != nil {
		resp.Diagnostics.AddError(
			"Create Managed Company Failed",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(apiResp.Message)

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

	// Execute msp-down command (setup/initialization)
	_, err := r.apiManager.ExecuteCommand(ctx, "msp-down", "Unable to initialize managed company service")
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Managed Company Failed",
			err.Error(),
		)
		return
	}

	/*
		TODO: We will pass the company id to the command to get the company info when that is implemente in commander cli
	*/

	// Build command to get all companies info
	command := "msp-info --format json -v"

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

	/*
		TODO: we will update node name when we get get in response from commander cli
	*/
	// Node: keep existing state value (skip for now as mentioned)

	state.Plan = types.StringValue(companyInfo.Plan)
	state.Seats = types.Int64Value(int64(companyInfo.Allocated))

	// Convert storage format: "100GB" -> "100gb", "1TB" -> "1tb", "10TB" -> "10tb"
	storageLower := strings.ToLower(companyInfo.Storage)
	state.FilePlan = types.StringValue(storageLower)

	/*
		TODO: we will update add-ons list when we get get in response from commander cli, we dont need below complete logic
	*/
	// Process addons: preserve number suffix from state and maintain state order
	// Step 1: Create a map of state add-ons with their base names (without number suffix)
	stateAddOnsMap := make(map[string]string) // base name -> full name with suffix
	stateAddOnsOrder := make([]string, 0)     // preserve order of state add-ons
	if !state.AddOns.IsNull() && !state.AddOns.IsUnknown() {
		for _, elem := range state.AddOns.Elements() {
			if strValue, ok := elem.(types.String); ok {
				addonValue := strValue.ValueString()
				// Extract base name (remove :N suffix if present)
				baseName := addonValue
				if matches := addOnWithNumberRegex.FindStringSubmatch(addonValue); matches != nil {
					baseName = matches[1]
				}
				stateAddOnsMap[baseName] = addonValue
				stateAddOnsOrder = append(stateAddOnsOrder, baseName)
			}
		}
	}

	// Step 2: Create a map of API add-ons (base name -> processed value with suffix)
	apiAddOnsMap := make(map[string]string) // base name -> processed value
	for _, addon := range companyInfo.Addons {
		// Check if we have this add-on in state with a number suffix
		if stateValue, exists := stateAddOnsMap[addon]; exists {
			// Use the value from state (preserves the number suffix like :100)
			apiAddOnsMap[addon] = stateValue
		} else if addOnsWithNumber[addon] && !strings.Contains(addon, ":") {
			// If add-on supports numbers and doesn't have a suffix, add ":1"
			apiAddOnsMap[addon] = addon + ":1"
		} else {
			// Use as-is
			apiAddOnsMap[addon] = addon
		}
	}

	// Step 3: Build processed addons list maintaining state order, then append new ones
	processedAddons := make([]string, 0)
	processedSet := make(map[string]bool) // track which add-ons we've added

	// First, add add-ons from state in state order (if they exist in API)
	for _, baseName := range stateAddOnsOrder {
		if processedValue, exists := apiAddOnsMap[baseName]; exists {
			processedAddons = append(processedAddons, processedValue)
			processedSet[baseName] = true
		}
		// If not in API, skip it (add-on was removed)
	}

	// Then, append any new add-ons from API that weren't in state
	for _, addon := range companyInfo.Addons {
		if !processedSet[addon] {
			processedAddons = append(processedAddons, apiAddOnsMap[addon])
		}
	}

	// Convert addons slice to types.List
	addOnsElements := make([]types.String, len(processedAddons))
	for i, addon := range processedAddons {
		addOnsElements[i] = types.StringValue(addon)
	}
	addOnsList, diags := types.ListValueFrom(ctx, types.StringType, addOnsElements)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	state.AddOns = addOnsList

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
