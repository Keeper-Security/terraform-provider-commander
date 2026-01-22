// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ----- NAME VALIDATOR --------------------------------
type nameValidator struct{}

func (v nameValidator) Description(ctx context.Context) string {
	return "Enterprise Node Name must be at least 1 character long."
}

func (v nameValidator) MarkdownDescription(ctx context.Context) string {
	return "Enterprise Node Name must be at least 1 character long."
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
			"Invalid Enterprise Node Name",
			"Enterprise Node Name must be at least 1 character long.")
	}
}

// ----- PARENT VALIDATOR --------------------------------
type parentValidator struct{}

func (v parentValidator) Description(ctx context.Context) string {
	return "Enterprise Node Parent Name must be at least 1 character long."
}

func (v parentValidator) MarkdownDescription(ctx context.Context) string {
	return "Enterprise Node Parent Name must be at least 1 character long."
}

func (v parentValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Skip validation if the value is unknown (e.g., from data source during plan)
	// or null (optional field not provided)
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()
	if len(value) < 1 {
		resp.Diagnostics.AddError(
			"Invalid Enterprise Node Parent Name",
			"Enterprise Node Parent Name must be at least 1 character long.")
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
