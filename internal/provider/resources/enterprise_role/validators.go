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

// ----- NAME VALIDATOR --------------------------------
type nameValidator struct{}

func (v nameValidator) Description(ctx context.Context) string {
	return "Enterprise Role Name must be at least 1 character long."
}

func (v nameValidator) MarkdownDescription(ctx context.Context) string {
	return "Enterprise Role Name must be at least 1 character long."
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
			"Invalid Enterprise Role Name",
			"Enterprise Role Name must be at least 1 character long.")
	}
}

// ----- MANAGED COMPANY VALIDATOR --------------------------------
type managedCompanyValidator struct{}

func (v managedCompanyValidator) Description(ctx context.Context) string {
	return "Managed Company Name must be at least 1 character long."
}

func (v managedCompanyValidator) MarkdownDescription(ctx context.Context) string {
	return "Managed Company Name must be at least 1 character long."
}

func (v managedCompanyValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the value is unknown (e.g., from data source during plan)
	// or null (optional field not provided)
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()
	if len(value) < 1 {
		resp.Diagnostics.AddError(
			"Invalid Managed Company Name",
			"Managed Company Name must be at least 1 character long.")
	}
}

// ----- Node VALIDATOR --------------------------------
type nodeValidator struct{}

func (v nodeValidator) Description(ctx context.Context) string {
	return "Node must be at least 1 character long."
}

func (v nodeValidator) MarkdownDescription(ctx context.Context) string {
	return "Node must be at least 1 character long."
}

func (v nodeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the value is unknown (e.g., from data source during plan)
	// or null (optional field not provided)
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()
	if len(value) < 1 {
		resp.Diagnostics.AddError(
			"Invalid Node Name",
			"Node must be at least 1 character long.")
	}
}

// ----- MANAGING NODES MAP VALIDATOR --------------------------------
// Validates that all map keys (node names) in managing_nodes are valid
type managingNodesMapValidator struct{}

func (v managingNodesMapValidator) Description(ctx context.Context) string {
	return "All managing node names (map keys) must be at least 1 character long."
}

func (v managingNodesMapValidator) MarkdownDescription(ctx context.Context) string {
	return "All managing node names (map keys) must be at least 1 character long."
}

func (v managingNodesMapValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	// Skip validation if the value is null or unknown (optional field)
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Get all map keys (node names)
	elements := req.ConfigValue.Elements()

	// Validate each key (node name)
	for key := range elements {
		if len(key) < 1 {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Invalid Managing Node Name",
				"Managing node name (map key) must be at least 1 character long.",
			)
		}
	}
}

// ----- USERS SET VALIDATOR --------------------------------

type usersValidator struct{}

func (v usersValidator) Description(ctx context.Context) string {
	return "Users set must not contain empty strings."
}

func (v usersValidator) MarkdownDescription(ctx context.Context) string {
	return "Users set must not contain empty strings."
}

func (v usersValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
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
				"Invalid User Type",
				fmt.Sprintf("Expected string, got: %T", elem),
			)
			continue
		}

		value := strValue.ValueString()

		// Check for empty strings
		if value == "" {
			resp.Diagnostics.AddError(
				"Empty User String",
				"Users set cannot contain empty strings. Each user must be a non-empty string.",
			)
		}
	}
}

// ----- TEAMS SET VALIDATOR --------------------------------

type teamsValidator struct{}

func (v teamsValidator) Description(ctx context.Context) string {
	return "Teams set must not contain empty strings."
}

func (v teamsValidator) MarkdownDescription(ctx context.Context) string {
	return "Teams set must not contain empty strings."
}

func (v teamsValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
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
				"Invalid Team Type",
				fmt.Sprintf("Expected string, got: %T", elem),
			)
			continue
		}

		value := strValue.ValueString()

		// Check for empty strings
		if value == "" {
			resp.Diagnostics.AddError(
				"Empty Team String",
				"Teams set cannot contain empty strings. Each team must be a non-empty string.",
			)
		}
	}
}

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
