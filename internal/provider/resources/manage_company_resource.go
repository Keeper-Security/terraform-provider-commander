// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
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

// ----------------------------------------------------------

// ----- NAME VALIDATOR --------------------------------
type nameValidator struct{}

func (v nameValidator) Description(ctx context.Context) string {
	return "Managed Company Name must be at least 1 character long."
}

func (v nameValidator) MarkdownDescription(ctx context.Context) string {
	return "Managed Company Name must be at least 1 character long."
}

func (v nameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the value is unknown (e.g., from data source during plan)
	// But still validate null and empty strings from user input
	if req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if len(value) < 1 {
		resp.Diagnostics.AddError(
			"Invalid Managed Company Name",
			"Managed Company Name must be at least 1 character long.")
	}
}

// ----- NODE VALIDATOR --------------------------------
type nodeValidator struct{}

func (v nodeValidator) Description(ctx context.Context) string {
	return "Managing Node name or ID."
}

func (v nodeValidator) MarkdownDescription(ctx context.Context) string {
	return "Managing Node name or ID."
}

func (v nodeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the value is unknown (e.g., from data source during plan)
	if req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if len(value) < 1 {
		resp.Diagnostics.AddError(
			"Invalid Managing Node Name or ID",
			"Managing Node name or ID must be at least 1 character long.")
	}
}

// ----- PLAN VALIDATOR --------------------------------
type planValidator struct{}

const (
	PlanBusiness       = "business"
	PlanBusinessPlus   = "businessPlus"
	PlanEnterprise     = "enterprise"
	PlanEnterprisePlus = "enterprisePlus"
)

var planOptions = []string{PlanBusiness, PlanBusinessPlus, PlanEnterprise, PlanEnterprisePlus}

func (v planValidator) Description(ctx context.Context) string {
	return "Must be one of: " + strings.Join(planOptions, ", ")
}

func (v planValidator) MarkdownDescription(ctx context.Context) string {
	return "Must be one of: `" + strings.Join(planOptions, "`, `") + "`"
}

func (v planValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the value is unknown (e.g., from data source during plan)
	if req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value != PlanBusiness && value != PlanBusinessPlus && value != PlanEnterprise && value != PlanEnterprisePlus {
		resp.Diagnostics.AddError(
			"Invalid Plan Value",
			fmt.Sprintf("Must be one of: `%s`. Got: %s", strings.Join(planOptions, "`, `"), value),
		)
	}
}

// -------------------------------- SEATS VALIDATOR --------------------------------
type seatsValidator struct{}

func (v seatsValidator) Description(ctx context.Context) string {
	return "You must enter a license count or enter -1 for unlimited. Enter \"0\" if you do not want users to be provisioned yet."
}

func (v seatsValidator) MarkdownDescription(ctx context.Context) string {
	return "You must enter a license count or enter -1 for unlimited. Enter \"0\" if you do not want users to be provisioned yet."
}

func (v seatsValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	// Skip validation if the value is null or unknown (optional field)
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueInt64()
	if value < -1 || (value < 0 && value != -1) {
		resp.Diagnostics.AddError(
			"Invalid Seats Value",
			fmt.Sprintf("You must enter a license count or enter -1 for unlimited. Enter \"0\" if you do not want users to be provisioned yet. Got: %d", value),
		)
	}
}

// ----- FILE PLAN VALIDATOR --------------------------------
type filePlanValidator struct{}

const (
	FilePlan100GB = "100gb"
	FilePlan1TB   = "1tb"
	FilePlan10TB  = "10tb"
)

// TODO: HERE WE NEED TO ADD PLAN BASED FILE PLAN OPTIONS BEC. DIFFERENT PLANS HAVE DIFFERENT FILE PLAN OPTIONS
/*
1. business plan - 100gb, 1tb, 10tb
2. businessPlus plan - 1tb, 10tb
3. enterprise plan - 100gb, 1tb, 10tb
4. enterprisePlus plan - 1tb, 10tb
*/

var filePlanOptions = []string{FilePlan100GB, FilePlan1TB, FilePlan10TB}
var filePlanPlusOptions = []string{FilePlan1TB, FilePlan10TB}

func (v filePlanValidator) Description(ctx context.Context) string {
	return "Must be one of: " + strings.Join(filePlanOptions, ", ")
}

func (v filePlanValidator) MarkdownDescription(ctx context.Context) string {
	return "Must be one of: `" + strings.Join(filePlanOptions, "`, `") + "`"
}

func (v filePlanValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the value is null or unknown (optional field)
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	// Get the plan value from config to determine valid file plan options
	var planValue types.String
	diags := req.Config.GetAttribute(ctx, req.Path.ParentPath().AtName("plan"), &planValue)
	if diags.HasError() || planValue.IsNull() || planValue.IsUnknown() {
		// Skip validation if plan is not available - plan validator will handle plan validation
		return
	}

	plan := planValue.ValueString()

	// Determine valid options based on plan type
	var validOptions []string
	var validOptionsStr string

	// Plus plans (businessPlus, enterprisePlus) only allow 1tb and 10tb
	if plan == PlanBusinessPlus || plan == PlanEnterprisePlus {
		validOptions = filePlanPlusOptions
		validOptionsStr = strings.Join(filePlanPlusOptions, "`, `")
	} else {
		// Non-plus plans (business, enterprise) allow 100gb, 1tb, and 10tb
		validOptions = filePlanOptions
		validOptionsStr = strings.Join(filePlanOptions, "`, `")
	}

	// Validate the value against the appropriate options
	isValid := false
	for _, option := range validOptions {
		if value == option {
			isValid = true
			break
		}
	}

	if !isValid {
		resp.Diagnostics.AddError(
			"Invalid File Storage Plan Value",
			fmt.Sprintf("For plan `%s`, file plan must be one of: `%s`. Got: %s", plan, validOptionsStr, value),
		)
	}
}

// ----- ADD-ONS VALIDATOR --------------------------------
type addOnsValidator struct{}

var (
	// Base add-ons (no number suffix)

	/*
		Note:
		consumer_breach_watch - when we use this add-on we get error from commander cli
		professional_services_silver_add_on - when we use this add-on it don't make any changes to the company.
		gold_professional_services_add_on - when we use this add-on it don't make any changes to the company.
		platinum_professional_services_add_on - when we use this add-on it don't make any changes to the company.

	*/
	baseAddOns = map[string]bool{
		"chat":                           true,
		"enterprise_audit_and_reporting": true,
		// "professional_services_silver_add_on":   true,
		// "gold_professional_services_add_on":     true,
		// "platinum_professional_services_add_on": true,
		"msp_service_and_support": true,
		// "consumer_breach_watch":    true,
		"enterprise_breach_watch":  true,
		"compliance_report":        true,
		"secrets_manager":          true,
		"password_rotation":        true,
		"remote_browser_isolation": true,
	}

	// Add-ons that can have :N suffix
	addOnsWithNumber = map[string]bool{
		"connection_manager":                true,
		"privileged_access_manager":         true,
		"keeper_endpoint_privilege_manager": true,
	}

	// Regex to match add-on with number: "connection_manager:5"
	addOnWithNumberRegex = regexp.MustCompile(`^([a-z_]+):(\d+)$`)
)

/*
COMMANDER ADD-ONS TO KEEPER ADMIN CONSOLE NAMING MAPPING ---> FOR REFERENCE ONLY

keeper_endpoint_privilege_manager -> Endpoint Manager

enterprise_breach_watch -> Breach Watch
compliance_report -> Compliance Reporting
enterprise_audit_and_reporting -> Advanced Reporting & Alerts Module
msp_service_and_support -> Dedicated Service & Support
privileged_access_manager -> Privileged Access Manager
secrets_manager -> Keeper Secrets Manager (KSM)
connection_manager -> Keeper Connection Manager (On-Prem)
remote_browser_isolation -> Remote Browser Isolation
chat -> KeeperChat

*/

func (v addOnsValidator) Description(ctx context.Context) string {
	return "Each add-on must be a valid add-on name. Add-ons with :N suffix can include a number (defaults to 1 if omitted)."
}

func (v addOnsValidator) MarkdownDescription(ctx context.Context) string {
	return "Each add-on must be a valid add-on name. Add-ons with `:N` suffix can include a number (defaults to 1 if omitted)."
}

func (v addOnsValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	// Skip validation if the value is null or unknown (optional field)
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()
	for idx, elem := range elements {
		elemPath := req.Path.AtListIndex(idx)

		// Get the string value
		strValue, ok := elem.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				elemPath,
				"Invalid Secure Add-On Type",
				fmt.Sprintf("Expected string, got: %T", elem),
			)
			continue
		}

		value := strValue.ValueString()

		// Check if it's a base add-on
		if baseAddOns[value] {
			continue // Valid
		}

		// Check if it's an add-on with number format (e.g., "connection_manager:5")
		if matches := addOnWithNumberRegex.FindStringSubmatch(value); matches != nil {
			addOnName := matches[1]
			numberStr := matches[2]

			// Validate the add-on name
			if !addOnsWithNumber[addOnName] {
				resp.Diagnostics.AddAttributeError(
					elemPath,
					"Invalid Secure Add-On",
					fmt.Sprintf("Secure Add-On '%s' does not support number suffix. Got: %s", addOnName, value),
				)
				continue
			}

			// Validate the number is positive
			number, err := strconv.Atoi(numberStr)
			if err != nil || number < 1 {
				resp.Diagnostics.AddAttributeError(
					elemPath,
					"Invalid Secure Add-On Number",
					fmt.Sprintf("Number suffix must be a positive integer. Got: %s", value),
				)
				continue
			}

			continue // Valid
		}

		// Check if it's an add-on with number that's missing the suffix (e.g., "connection_manager")
		// This is valid - will default to :1 in Create/Update
		if addOnsWithNumber[value] {
			continue // Valid - will be normalized to :1 later
		}

		// Invalid add-on
		validOptions := []string{}
		for k := range baseAddOns {
			validOptions = append(validOptions, k)
		}
		for k := range addOnsWithNumber {
			validOptions = append(validOptions, k, k+":N")
		}

		resp.Diagnostics.AddAttributeError(
			elemPath,
			"Invalid Secure Add-On",
			fmt.Sprintf("Invalid Secure Add-On '%s'. Must be one of: %s", value, strings.Join(validOptions, ", ")),
		)
	}
}

// ----------------------------------------------------------

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
				Description:         "Secure Add-Ons to apply to the Managed Company.",
				MarkdownDescription: "Secure Add-Ons to apply to the Managed Company.",
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

	// Node: keep existing state value (skip for now as mentioned)
	state.Plan = types.StringValue(companyInfo.Plan)
	state.Seats = types.Int64Value(int64(companyInfo.Allocated))

	// Convert storage format: "100GB" -> "100gb", "1TB" -> "1tb", "10TB" -> "10tb"
	storageLower := strings.ToLower(companyInfo.Storage)
	state.FilePlan = types.StringValue(storageLower)

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

// ------------------HELPER FUNCTIONS----------------------------------------

func buildManageCompanyAddCommand(data ManageCompanyResourceModel) string {
	var parts []string

	parts = append(parts, "msp-add")

	// Required parameters

	if !data.Name.IsNull() {
		parts = append(parts, fmt.Sprintf("'%s'", data.Name.ValueString()))
	}

	if !data.Node.IsNull() {
		parts = append(parts, fmt.Sprintf("--node '%s'", data.Node.ValueString()))
	}

	if !data.Plan.IsNull() {
		parts = append(parts, fmt.Sprintf("--plan '%s'", data.Plan.ValueString()))
	}

	// Optional parameters

	if !data.Seats.IsNull() {
		parts = append(parts, fmt.Sprintf("--seats %d", data.Seats.ValueInt64()))
	}

	if !data.FilePlan.IsNull() {
		parts = append(parts, fmt.Sprintf("--file-plan '%s'", data.FilePlan.ValueString()))
	}

	if !data.AddOns.IsNull() && !data.AddOns.IsUnknown() {
		addOnsList := normalizeAddOns(data.AddOns)

		for _, addOn := range addOnsList {
			parts = append(parts, fmt.Sprintf("--addon '%s'", addOn))
		}
	}

	return strings.Join(parts, " ")
}

func buildManageCompanyUpdateCommand(plan *ManageCompanyResourceModel, state *ManageCompanyResourceModel) string {
	var parts []string

	parts = append(parts, "msp-update")
	parts = append(parts, fmt.Sprintf("'%s'", state.Id.ValueString()))

	// TODO: CHECK IN COMMANDER FOR UPDATING NAME - syntax
	// if !state.Name.Equal(plan.Name) {
	// 	parts = append(parts, fmt.Sprintf("--name '%s'", plan.Name.ValueString()))
	// }

	// Required fields - no null check needed (validation ensures they exist)
	if !state.Node.Equal(plan.Node) {
		parts = append(parts, fmt.Sprintf("--node '%s'", plan.Node.ValueString()))
	}

	if !state.Plan.Equal(plan.Plan) {
		parts = append(parts, fmt.Sprintf("--plan '%s'", plan.Plan.ValueString()))
	}

	// Optional fields
	if !state.Seats.Equal(plan.Seats) {
		// check null (user might have removed them)
		if plan.Seats.IsNull() {
			parts = append(parts, fmt.Sprintf("--seats %d", 0))

		} else {
			parts = append(parts, fmt.Sprintf("--seats %d", plan.Seats.ValueInt64()))
		}
	}

	if !state.FilePlan.Equal(plan.FilePlan) && !plan.FilePlan.IsNull() {
		if plan.FilePlan.IsNull() {
			parts = append(parts, fmt.Sprintf("--file-plan '%s'", FilePlan100GB))
		} else {
			parts = append(parts, fmt.Sprintf("--file-plan '%s'", plan.FilePlan.ValueString()))
		}
	}

	// Add-ons

	// if there is add-on in state but not in plan, means we need to remove it by adding --remove-addon flag and appending add-ons
	// if there is add-on in plan but not in state, means we need to add it
	// if there is add-on with same base name but different number (e.g., connection_manager:2 -> connection_manager:1), we need to update it

	// Normalize add-ons from state and plan
	stateAddOns := normalizeAddOns(state.AddOns)
	planAddOns := normalizeAddOns(plan.AddOns)

	// Convert to maps for easier comparison
	stateAddOnsMap := make(map[string]bool)
	for _, addOn := range stateAddOns {
		stateAddOnsMap[addOn] = true
	}

	planAddOnsMap := make(map[string]bool)
	for _, addOn := range planAddOns {
		planAddOnsMap[addOn] = true
	}

	// Helper function to extract base add-on name (without number suffix)
	getBaseAddOnName := func(addOn string) string {
		if matches := addOnWithNumberRegex.FindStringSubmatch(addOn); matches != nil {
			return matches[1]
		}
		return addOn
	}

	// Find add-ons to remove (in state but not in plan)
	var addOnsToRemove []string
	// Track base names that are being updated (same base, different number)
	updatedBaseNames := make(map[string]bool)

	for _, stateAddOn := range stateAddOns {
		if !planAddOnsMap[stateAddOn] {
			// Check if this add-on has a number suffix and if there's a matching base name in plan
			baseName := getBaseAddOnName(stateAddOn)
			hasUpdate := false

			// Check if there's an add-on in plan with the same base name but different number
			for _, planAddOn := range planAddOns {
				if getBaseAddOnName(planAddOn) == baseName && planAddOn != stateAddOn {
					// This is an update case (same base, different number)
					updatedBaseNames[baseName] = true
					hasUpdate = true
					break
				}
			}

			// Only add to remove list if it's not an update case
			if !hasUpdate {
				addOnsToRemove = append(addOnsToRemove, stateAddOn)
			}
		}
	}

	// Find add-ons to add (in plan but not in state, or updates with different numbers)
	var addOnsToAdd []string
	for _, planAddOn := range planAddOns {
		if !stateAddOnsMap[planAddOn] {
			// Check if this is an update case (same base name exists in state with different number)
			baseName := getBaseAddOnName(planAddOn)
			if updatedBaseNames[baseName] {
				// This is an update - include it in add list
				addOnsToAdd = append(addOnsToAdd, planAddOn)
			} else {
				// This is a new add-on
				addOnsToAdd = append(addOnsToAdd, planAddOn)
			}
		}
	}

	// Add remove-addon flag for each add-on to remove
	for _, addOn := range addOnsToRemove {
		parts = append(parts, fmt.Sprintf("--remove-addon '%s'", addOn))
	}

	// Add add-addon flag for each add-on to add
	for _, addOn := range addOnsToAdd {
		parts = append(parts, fmt.Sprintf("--add-addon '%s'", addOn))
	}

	return strings.Join(parts, " ")
}

// function to normalize add-ons (add :1 if missing for add-ons that support it)
func normalizeAddOns(addOns types.List) []string {
	if addOns.IsNull() || addOns.IsUnknown() {
		return nil
	}

	elements := addOns.Elements()
	result := make([]string, 0, len(elements))

	for _, elem := range elements {
		if strValue, ok := elem.(types.String); ok {
			value := strValue.ValueString()
			// Add :1 if missing for add-ons that support it
			if addOnsWithNumber[value] && !strings.Contains(value, ":") {
				value = value + ":1"
			}
			result = append(result, value)
		}
	}

	return result
}

func NewManageCompanyResource() resource.Resource {
	return &ManageCompanyResource{}
}
