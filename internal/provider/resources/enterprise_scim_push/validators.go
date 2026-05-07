// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// sourceValidator validates that source is one of google, ad, record.
type sourceValidator struct{}

func (sourceValidator) Description(ctx context.Context) string {
	return "Source must be one of: google, ad, record."
}

func (sourceValidator) MarkdownDescription(ctx context.Context) string {
	return "Source must be one of: `google`, `ad`, `record`."
}

func (v sourceValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	val := strings.TrimSpace(strings.ToLower(req.ConfigValue.ValueString()))
	switch val {
	case SourceGoogle, SourceAD, SourceRecord:
		return
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid source",
		"Source must be one of: google, ad, record. Got: "+req.ConfigValue.ValueString(),
	)
}
