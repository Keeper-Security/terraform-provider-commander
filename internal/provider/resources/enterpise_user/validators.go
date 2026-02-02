package enterpiseuser

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ----- EMAIL VALIDATOR --------------------------------

type emailValidator struct{}

func (v emailValidator) Description(ctx context.Context) string {
	return "Email must be a valid email address."
}

func (v emailValidator) MarkdownDescription(ctx context.Context) string {
	return "Email must be a valid email address."
}

func (v emailValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	val := req.ConfigValue.ValueString()
	if !regexp.MustCompile(EmailRegex).MatchString(val) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Email Address",
			fmt.Sprintf("Email must be a valid email address. Got: %s", val),
		)
	}
}

// ----- NAME VALIDATOR --------------------------------

type nameValidator struct{}

func (v nameValidator) Description(ctx context.Context) string {
	return "Name must be at least 1 character long."
}

func (v nameValidator) MarkdownDescription(ctx context.Context) string {
	return "Name must be at least 1 character long."
}

func (v nameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	value := req.ConfigValue.ValueString()
	if len(value) < 1 {
		resp.Diagnostics.AddError(
			"Invalid Name",
			fmt.Sprintf("Name must be at least 1 character long. Got: %s", value),
		)
	}
}

// ----- JOB TITLE VALIDATOR --------------------------------

type jobTitleValidator struct{}

func (v jobTitleValidator) Description(ctx context.Context) string {
	return "Job title must be at least 1 character long."
}

func (v jobTitleValidator) MarkdownDescription(ctx context.Context) string {
	return "Job title must be at least 1 character long."
}

func (v jobTitleValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	value := req.ConfigValue.ValueString()
	if len(value) < 1 {
		resp.Diagnostics.AddError(
			"Invalid Job Title",
			fmt.Sprintf("Job title must be at least 1 character long. Got: %s", value),
		)
	}
}
