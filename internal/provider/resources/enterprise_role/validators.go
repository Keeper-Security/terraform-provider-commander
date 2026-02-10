// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ----- ENFORCEMENT POLICY KEY VALIDATOR --------------------------------
// ValidEnforcementPolicyKeys is defined in constants.go and contains all valid enforcement policy keys

type enforcementPolicyKeyValidator struct{}

func (v enforcementPolicyKeyValidator) Description(ctx context.Context) string {
	return "Enforcement policy key must be one of the valid policy keys."
}

func (v enforcementPolicyKeyValidator) MarkdownDescription(ctx context.Context) string {
	return "Enforcement policy key must be one of the valid policy keys. See documentation for the complete list."
}

func (v enforcementPolicyKeyValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the value is unknown (e.g., from data source during plan)
	if req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if value == "" {
		resp.Diagnostics.AddError(
			"Invalid Enforcement Policy Key",
			"Enforcement policy key cannot be empty.")
		return
	}

	// Check if the value is in the valid keys list
	isValid := false
	for _, validKey := range ValidEnforcementPolicyKeys {
		if value == validKey {
			isValid = true
			break
		}
	}

	if !isValid {
		resp.Diagnostics.AddError(
			"Invalid Enforcement Policy Key",
			fmt.Sprintf("Enforcement policy key '%s' is not valid. Must be one of: %s", value, strings.Join(ValidEnforcementPolicyKeys, ", ")))
	}
}

// ----- ENFORCEMENT POLICIES MAP KEY VALIDATOR --------------------------------
// Validates that all map keys (policy keys) in enforcement_policies are valid.
type enforcementPoliciesMapKeyValidator struct{}

func (v enforcementPoliciesMapKeyValidator) Description(ctx context.Context) string {
	return "All enforcement policy keys (map keys) must be valid enforcement policy keys."
}

func (v enforcementPoliciesMapKeyValidator) MarkdownDescription(ctx context.Context) string {
	return "All enforcement policy keys (map keys) must be valid enforcement policy keys."
}

func (v enforcementPoliciesMapKeyValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	// Skip validation if the value is null or unknown (optional field)
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()

	for key := range elements {
		// Validate value type and non-empty string (MapAttribute value is types.String)
		rawVal := elements[key]
		strVal, ok := rawVal.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Invalid Enforcement Policy Value Type",
				fmt.Sprintf("Expected string value for enforcement policy '%s', got: %T", key, rawVal),
			)
			continue
		}

		// Skip value validation for unknowns (common during plan), but fail fast on explicit null/empty.
		if strVal.IsNull() {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Invalid Enforcement Policy Value",
				fmt.Sprintf("Enforcement policy value for key '%s' cannot be null.", key),
			)
		} else if !strVal.IsUnknown() && strVal.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Invalid Enforcement Policy Value",
				fmt.Sprintf("Enforcement policy value for key '%s' cannot be an empty string.", key),
			)
		}

		if key == "" {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Enforcement Policy Key",
				"Enforcement policy key (map key) cannot be empty.",
			)
			continue
		}

		isValid := false
		for _, validKey := range ValidEnforcementPolicyKeys {
			if key == validKey {
				isValid = true
				break
			}
		}

		if !isValid {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Invalid Enforcement Policy Key",
				fmt.Sprintf("Enforcement policy key '%s' is not valid. Must be one of: %s", key, strings.Join(ValidEnforcementPolicyKeys, ", ")),
			)
		}
	}
}

// ----- PRIVILEGES SET VALIDATOR --------------------------------
type privilegesValidator struct{}

func (v privilegesValidator) Description(ctx context.Context) string {
	return "Privileges must be valid privilege values."
}

func (v privilegesValidator) MarkdownDescription(ctx context.Context) string {
	return "Privileges must be valid privilege values. Valid values: " + strings.Join(ValidPrivileges, ", ")
}

func (v privilegesValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	// Skip validation if the value is null or unknown (optional field)
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()
	for _, elem := range elements {
		// Get the string value
		strValue, ok := elem.(types.String)
		if !ok {
			// This shouldn't happen if ElementType is correct, but handle it anyway
			resp.Diagnostics.AddError(
				"Invalid Privilege Type",
				fmt.Sprintf("Expected string, got: %T", elem),
			)
			continue
		}

		// Skip unknown values (e.g. from data source references not yet evaluated at plan time)
		if strValue.IsUnknown() {
			continue
		}

		value := strValue.ValueString()

		// Check for empty strings
		if value == "" {
			resp.Diagnostics.AddError(
				"Empty Privilege String",
				"Privileges set cannot contain empty strings. Each privilege must be a non-empty string.",
			)
			continue
		}

		// Check if the value is in the valid privileges list
		isValid := false
		for _, validPrivilege := range ValidPrivileges {
			if value == validPrivilege {
				isValid = true
				break
			}
		}

		if !isValid {
			resp.Diagnostics.AddError(
				"Invalid Privilege",
				fmt.Sprintf("Privilege '%s' is not valid. Must be one of: %s", value, strings.Join(ValidPrivileges, ", ")),
			)
		}
	}
}
