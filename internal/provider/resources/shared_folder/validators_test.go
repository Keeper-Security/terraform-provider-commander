// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharedfolder

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExpirationValidator(t *testing.T) {
	v := ExpirationValidator()
	ctx := context.Background()

	tests := []struct {
		name      string
		value     string
		wantError bool
	}{
		{"never", "never", false},
		{"never case insensitive", "Never", false},
		{"ISO date", "2026-12-31", false},
		{"ISO datetime", "2026-12-31 23:59:59", false},
		{"ISO datetime no seconds", "2026-12-31 14:30", false},
		{"relative 30d", "30d", false},
		{"relative 1y", "1y", false},
		{"relative 6mo", "6mo", false},
		{"relative 24h", "24h", false},
		{"relative 90days", "90days", false},
		{"relative 30minutes", "30minutes", false},
		{"empty", "", true},
		{"invalid format", "invalid", true},
		{"bad date", "2026-13-45", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validator.StringRequest{ConfigValue: types.StringValue(tt.value)}
			var resp validator.StringResponse
			v.ValidateString(ctx, req, &resp)
			hasError := resp.Diagnostics.HasError()
			if hasError != tt.wantError {
				t.Errorf("value %q: wantError=%v, got diagnostics: %v", tt.value, tt.wantError, resp.Diagnostics)
			}
		})
	}
}
