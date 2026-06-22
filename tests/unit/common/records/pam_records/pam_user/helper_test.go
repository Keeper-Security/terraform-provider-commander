// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser_test

import (
	"testing"

	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam_records/pam_user"
)

func TestParsePasswordComplexityData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		raw    string
		want   string
		isNull bool
	}{
		{
			name: "commander rotation info format",
			raw:  "Length: 32; Lowercase: 1; Uppercase: 5; Digits: 1; Symbols: 2; Symbols Chars: !@#$%^?();',.=+[]<>{}-_/\\*&:\"`~|",
			want: "32,5,1,1,2",
		},
		{
			name: "symbols chars contains semicolons and fake field names",
			raw:  "Length: 32; Lowercase: 1; Uppercase: 5; Digits: 1; Symbols: 2; Symbols Chars: !; Uppercase: 999; Lowercase: 888",
			want: "32,5,1,1,2",
		},
		{
			name: "different field order",
			raw:  "Symbols: 0; Digits: 3; Uppercase: 2; Lowercase: 4; Length: 24",
			want: "24,2,4,3,0",
		},
		{
			name:   "empty",
			raw:    "",
			isNull: true,
		},
		{
			name:   "missing symbols",
			raw:    "Length: 32; Lowercase: 1; Uppercase: 5; Digits: 1",
			isNull: true,
		},
		{
			name:   "non-numeric value",
			raw:    "Length: abc; Lowercase: 1; Uppercase: 5; Digits: 1; Symbols: 2",
			isNull: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := commonpamuser.ParsePasswordComplexityData(tc.raw)
			if tc.isNull {
				if !got.IsNull() {
					t.Fatalf("expected null, got %q", got.ValueString())
				}
				return
			}
			if got.IsNull() {
				t.Fatalf("expected %q, got null", tc.want)
			}
			if got.ValueString() != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got.ValueString())
			}
		})
	}
}

func TestParseRotationInfoMessage_Complexity(t *testing.T) {
	t.Parallel()

	var state commonpamuser.PamUserSharedModel
	commonpamuser.ParseRotationInfoMessage([]string{
		"Password Complexity Data: Length: 32; Lowercase: 1; Uppercase: 5; Digits: 1; Symbols: 2; Symbols Chars: !@#$",
	}, nil, &state)

	if state.RotationSettings == nil {
		t.Fatal("expected rotation settings")
	}
	if state.RotationSettings.Complexity.IsNull() {
		t.Fatal("expected complexity value")
	}
	if got := state.RotationSettings.Complexity.ValueString(); got != "32,5,1,1,2" {
		t.Fatalf("expected 32,5,1,1,2, got %q", got)
	}
}
