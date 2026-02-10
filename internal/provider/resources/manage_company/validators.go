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

// ----- PLAN VALIDATOR --------------------------------
type planValidator struct{}

// PlanOptions contains all valid plan options
var PlanOptions = []string{PlanBusiness, PlanBusinessPlus, PlanEnterprise, PlanEnterprisePlus}

func (v planValidator) Description(ctx context.Context) string {
	return "Must be one of: " + strings.Join(PlanOptions, ", ")
}

func (v planValidator) MarkdownDescription(ctx context.Context) string {
	return "Must be one of: `" + strings.Join(PlanOptions, "`, `") + "`"
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
			fmt.Sprintf("Must be one of: `%s`. Got: %s", strings.Join(PlanOptions, "`, `"), value),
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

/*
1. business plan - 100gb, 1tb, 10tb
2. businessPlus plan - 1tb, 10tb
3. enterprise plan - 100gb, 1tb, 10tb
4. enterprisePlus plan - 1tb, 10tb
*/

// FilePlanOptions contains all valid file plan options
var FilePlanOptions = []string{FilePlan100GB, FilePlan1TB, FilePlan10TB}

// FilePlanPlusOptions contains file plan options for plus plans (businessPlus, enterprisePlus)
var FilePlanPlusOptions = []string{FilePlan1TB, FilePlan10TB}

func (v filePlanValidator) Description(ctx context.Context) string {
	return "Must be one of: " + strings.Join(FilePlanOptions, ", ")
}

func (v filePlanValidator) MarkdownDescription(ctx context.Context) string {
	return "Must be one of: `" + strings.Join(FilePlanOptions, "`, `") + "`"
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
		validOptions = FilePlanPlusOptions
		validOptionsStr = strings.Join(FilePlanPlusOptions, "`, `")
	} else {
		// Non-plus plans (business, enterprise) allow 100gb, 1tb, and 10tb
		validOptions = FilePlanOptions
		validOptionsStr = strings.Join(FilePlanOptions, "`, `")
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
	// Regex to match add-on with number: "connection_manager:5"
	addOnWithNumberRegex = regexp.MustCompile(`^([a-z_]+):(\d+)$`)
)

// Base add-ons (no number suffix)
//
// Note: these add-ons are not working in commander cli
// consumer_breach_watch - when we use this add-on we get error from commander cli
// professional_services_silver_add_on - when we use this add-on it don't make any changes to the company.
// gold_professional_services_add_on - when we use this add-on it don't make any changes to the company.
// platinum_professional_services_add_on - when we use this add-on it don't make any changes to the company.
var BaseAddOns = map[string]bool{
	AddOnChat:                        true,
	AddOnEnterpriseAuditAndReporting: true,
	// "professional_services_silver_add_on":   true,
	// "gold_professional_services_add_on":     true,
	// "platinum_professional_services_add_on": true,
	AddOnMspServiceAndSupport: true,
	// "consumer_breach_watch":    true,
	AddOnEnterpriseBreachWatch:  true,
	AddOnComplianceReport:       true,
	AddOnSecretsManager:         true,
	AddOnPasswordRotation:       true,
	AddOnRemoteBrowserIsolation: true,
}

// AddOnsWithNumber are add-ons that can have :N suffix
var AddOnsWithNumber = map[string]bool{
	AddOnConnectionManager:              true,
	AddOnPrivilegedAccessManager:        true,
	AddOnKeeperEndpointPrivilegeManager: true,
}

func (v addOnsValidator) Description(ctx context.Context) string {
	return "Each add-on must be a valid add-on name. Add-ons that support numbers must include the :N suffix (e.g., connection_manager:3)."
}

func (v addOnsValidator) MarkdownDescription(ctx context.Context) string {
	return "Each add-on must be a valid add-on name. Add-ons that support numbers must include the `:N` suffix (e.g., `connection_manager:3`)."
}

func (v addOnsValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	// Skip validation if the value is null or unknown (optional field)
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()

	// Track addons with numbers for cross-validation
	var privilegedAccessManagerNum *int
	var connectionManagerNum *int

	// Track required addons for privileged_access_manager dependency validation
	hasPrivilegedAccessManager := false
	hasSecretsManager := false
	hasConnectionManager := false
	hasRemoteBrowserIsolation := false

	// validate individual addons and track presence
	for _, elem := range elements {
		// Get the string value
		strValue, ok := elem.(types.String)
		if !ok {
			// This shouldn't happen if ElementType is correct, but handle it anyway
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Secure Add-On Type",
				fmt.Sprintf("Expected string, got: %T", elem),
			)
			continue
		}

		// Skip unknown values (e.g. from variable or data source reference not yet evaluated at plan time)
		if strValue.IsUnknown() {
			continue
		}

		value := strValue.ValueString()

		// Track presence of specific addons for dependency validation
		switch value {
		case AddOnSecretsManager:
			hasSecretsManager = true
		case AddOnRemoteBrowserIsolation:
			hasRemoteBrowserIsolation = true
		}

		// Check if it's a base add-on
		if BaseAddOns[value] {
			continue // Valid
		}

		// Check if it's an add-on with number format (e.g., "connection_manager:5")
		if matches := addOnWithNumberRegex.FindStringSubmatch(value); matches != nil {
			addOnName := matches[1]
			numberStr := matches[2]

			// Validate the add-on name
			if !AddOnsWithNumber[addOnName] {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid Secure Add-On",
					fmt.Sprintf("Secure Add-On '%s' does not support number suffix. Got: %s", addOnName, value),
				)
				continue
			}

			// Validate the number is positive
			number, err := strconv.Atoi(numberStr)
			if err != nil || number < 1 {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid Secure Add-On Number",
					fmt.Sprintf("Number suffix must be a positive integer. Got: %s", value),
				)
				continue
			}

			// Track numbers for privileged_access_manager and connection_manager
			switch addOnName {
			case AddOnPrivilegedAccessManager:
				privilegedAccessManagerNum = &number
				hasPrivilegedAccessManager = true
			case AddOnConnectionManager:
				connectionManagerNum = &number
				hasConnectionManager = true
			}

			continue // Valid
		}

		// Check if it's an add-on that requires number suffix but doesn't have it
		if AddOnsWithNumber[value] {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Missing Number Suffix",
				fmt.Sprintf("Secure Add-On '%s' requires a number suffix in the format '%s:N' where N is a positive integer (e.g., '%s:3').", value, value, value),
			)
			continue
		}

		// Invalid add-on
		validOptions := []string{}
		for k := range BaseAddOns {
			validOptions = append(validOptions, k)
		}
		for k := range AddOnsWithNumber {
			validOptions = append(validOptions, k+":N")
		}

		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Secure Add-On",
			fmt.Sprintf("Invalid Secure Add-On '%s'. Must be one of: %s", value, strings.Join(validOptions, ", ")),
		)
	}

	// validate that privileged_access_manager and connection_manager have matching N values
	if privilegedAccessManagerNum != nil && connectionManagerNum != nil {
		if *privilegedAccessManagerNum != *connectionManagerNum {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Mismatched Add-On Numbers",
				fmt.Sprintf("When both '%s' and '%s' are provided, their number suffixes must match. Got %s:%d and %s:%d", AddOnPrivilegedAccessManager, AddOnConnectionManager, AddOnPrivilegedAccessManager, *privilegedAccessManagerNum, AddOnConnectionManager, *connectionManagerNum),
			)
		}
	}

	// validate that if privileged_access_manager is present, required addons must also be present
	if hasPrivilegedAccessManager {
		var missingAddons []string

		if !hasSecretsManager {
			missingAddons = append(missingAddons, AddOnSecretsManager)
		}

		if !hasConnectionManager {
			missingAddons = append(missingAddons, AddOnConnectionManager+":N (with matching N value)")
		}

		if !hasRemoteBrowserIsolation {
			missingAddons = append(missingAddons, AddOnRemoteBrowserIsolation)
		}

		if len(missingAddons) > 0 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Missing Required Add-Ons",
				fmt.Sprintf("When '%s' is provided, the following add-ons are required: %s", AddOnPrivilegedAccessManager, strings.Join(missingAddons, ", ")),
			)
		}
	}
}

// GetAllValidAddOns returns a list of all valid add-on names
// This includes base add-ons and add-ons that support number suffixes
func GetAllValidAddOns() []string {
	validOptions := []string{}

	// Add base add-ons
	for k := range BaseAddOns {
		validOptions = append(validOptions, k)
	}

	// Add add-ons with number support (must include :N format)
	for k := range AddOnsWithNumber {
		validOptions = append(validOptions, k+":N")
	}

	return validOptions
}
