// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package folder

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type colorValidator struct{}

func (colorValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("value must be one of: %s", strings.Join(ValidColors, ", "))
}

func (colorValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("value must be one of: `%s`", strings.Join(ValidColors, "`, `"))
}

func (v colorValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	val := req.ConfigValue.ValueString()
	for _, c := range ValidColors {
		if val == c {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid folder color",
		fmt.Sprintf("%q is not a valid color. Must be one of: %s", val, strings.Join(ValidColors, ", ")),
	)
}

// ColorValidator returns a validator that checks folder color values.
func ColorValidator() colorValidator {
	return colorValidator{}
}
