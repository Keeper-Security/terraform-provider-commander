// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

	// TODO: Need to update validation for keeper_endpoint_privilege_manager to validate number that is supported in admin console - this validation also needs to be done in commander cli
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

// GetAllValidAddOns returns a list of all valid add-on names
// This includes base add-ons and add-ons that support number suffixes
func GetAllValidAddOns() []string {
	validOptions := []string{}

	// Add base add-ons
	for k := range baseAddOns {
		validOptions = append(validOptions, k)
	}

	// Add add-ons with number support (show both with and without :N format)
	for k := range addOnsWithNumber {
		validOptions = append(validOptions, k, k+":N")
	}

	return validOptions
}
