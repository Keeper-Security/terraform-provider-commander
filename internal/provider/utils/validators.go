// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ----- Node VALIDATOR --------------------------------
type NodeValidator struct{}

func (v NodeValidator) Description(ctx context.Context) string {
	return "Node must be at least 1 character long."
}

func (v NodeValidator) MarkdownDescription(ctx context.Context) string {
	return "Node must be at least 1 character long."
}

func (v NodeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
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
type ManagedCompanyValidator struct{}

func (v ManagedCompanyValidator) Description(ctx context.Context) string {
	return "Managed Company Name must be at least 1 character long."
}

func (v ManagedCompanyValidator) MarkdownDescription(ctx context.Context) string {
	return "Managed Company Name must be at least 1 character long."
}

func (v ManagedCompanyValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
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

// ----- TEAMS SET VALIDATOR --------------------------------

type TeamsValidator struct{}

func (v TeamsValidator) Description(ctx context.Context) string {
	return "Teams set must not contain empty strings."
}

func (v TeamsValidator) MarkdownDescription(ctx context.Context) string {
	return "Teams set must not contain empty strings."
}

func (v TeamsValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
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

// ----- ROLES SET VALIDATOR --------------------------------
type RolesValidator struct{}

func (v RolesValidator) Description(ctx context.Context) string {
	return "Roles set must not contain empty strings."
}

func (v RolesValidator) MarkdownDescription(ctx context.Context) string {
	return "Roles set must not contain empty strings."
}

func (v RolesValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
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
