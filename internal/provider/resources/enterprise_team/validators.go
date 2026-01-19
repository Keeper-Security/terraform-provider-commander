// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ----- NAME VALIDATOR --------------------------------
type nameValidator struct{}

func (v nameValidator) Description(ctx context.Context) string {
	return "Enterprise Team Name must be at least 1 character long."
}

func (v nameValidator) MarkdownDescription(ctx context.Context) string {
	return "Enterprise Team Name must be at least 1 character long."
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
			"Invalid Enterprise Team Name",
			"Enterprise Team Name must be at least 1 character long.")
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

// ----- ROLES SET VALIDATOR --------------------------------
type rolesValidator struct{}

func (v rolesValidator) Description(ctx context.Context) string {
	return "Roles set must not contain empty strings."
}

func (v rolesValidator) MarkdownDescription(ctx context.Context) string {
	return "Roles set must not contain empty strings."
}

func (v rolesValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
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
				"Invalid Role Type",
				fmt.Sprintf("Expected string, got: %T", elem),
			)
			continue
		}

		value := strValue.ValueString()

		// Check for empty strings
		if value == "" {
			resp.Diagnostics.AddError(
				"Empty Role String",
				"Roles set cannot contain empty strings. Each role must be a non-empty string.",
			)
		}
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
