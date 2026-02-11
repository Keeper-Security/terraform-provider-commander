// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser

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
