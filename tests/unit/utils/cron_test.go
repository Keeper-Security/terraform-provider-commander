// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils_test

import (
	"strings"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils/cronvalidate"
)

func TestValidateKeeperRotationCron_ValidKeeperExamples(t *testing.T) {
	t.Parallel()

	valid := []string{
		"0 28 17 ? * *",
		"0 0 12 * * ?",
		"0 15 10 ? * *",
		"0 15 10 * * ?",
		"0 15 10 * * ? *",
		"0 15 10 * * ? 2005",
		"0 10,44 14 ? 3 WED",
		"0 15 10 ? * MON-FRI",
		"0 15 10 15 * ?",
		"0 15 10 L * ?",
		"0 15 10 L-2 * ?",
		"0 15 10 ? * 6L",
		"0 15 10 ? * 6L 2002-2005",
		"0 15 10 ? * 6#3",
		"0 0 12 1/5 * ?",
		"0 11 11 11 11 ?",
	}

	for _, expr := range valid {
		t.Run(expr, func(t *testing.T) {
			t.Parallel()
			if err := cronvalidate.ValidateKeeperRotationCron(expr); err != nil {
				t.Fatalf("expected valid cron %q, got error: %v", expr, err)
			}
		})
	}
}

func TestValidateKeeperRotationCron_Invalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expr     string
		contains string
	}{
		{
			name:     "empty",
			expr:     "",
			contains: "cannot be empty",
		},
		{
			name:     "five field unix cron",
			expr:     "56 17 * * *",
			contains: "6 or 7",
		},
		{
			name:     "descriptor preset",
			expr:     "@daily",
			contains: "6 or 7",
		},
		{
			name:     "both dom and dow specified",
			expr:     "0 15 10 15 * MON",
			contains: "cannot both be specified",
		},
		{
			name:     "both dom and dow question mark",
			expr:     "0 15 10 ? * ?",
			contains: "cannot both be '?'",
		},
		{
			name:     "every minute",
			expr:     "0 * * ? * *",
			contains: "1-hour interval",
		},
		{
			name:     "half hour increment",
			expr:     "0 */30 * ? * *",
			contains: "1-hour interval",
		},
		{
			name:     "seconds increment",
			expr:     "0/15 0 12 ? * *",
			contains: "single value",
		},
		{
			name:     "invalid month name",
			expr:     "0 0 12 ? FOO *",
			contains: "month",
		},
		{
			name:     "invalid hour",
			expr:     "0 0 25 ? * *",
			contains: "hours",
		},
		{
			name:     "invalid dom character",
			expr:     "0 0 12 # * ?",
			contains: "invalid character",
		},
		{
			name:     "year out of range",
			expr:     "0 0 12 ? * MON 1800",
			contains: "year",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := cronvalidate.ValidateKeeperRotationCron(tc.expr)
			if err == nil {
				t.Fatalf("expected error for %q", tc.expr)
			}
			if tc.contains != "" && !strings.Contains(err.Error(), tc.contains) {
				t.Fatalf("error %q should contain %q", err.Error(), tc.contains)
			}
		})
	}
}
