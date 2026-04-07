// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package epmpolicy_test

import (
	"context"
	"strings"
	"testing"

	commonepm "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/epm_policy"
)

func TestDescriptionFunctions_NonEmpty(t *testing.T) {
	checks := []func() string{
		commonepm.PolicyTypeDescription,
		commonepm.StatusDescription,
		commonepm.StatusDescriptionForLeastPrivilege,
		commonepm.ControlDescription,
		commonepm.DayFilterDescription,
	}
	for i, fn := range checks {
		if s := fn(); strings.TrimSpace(s) == "" {
			t.Fatalf("check %d: empty string", i)
		}
	}
}

func TestMarkdownDescriptionFunctions(t *testing.T) {
	ctx := context.Background()
	if commonepm.PolicyTypeMarkdown() == "" {
		t.Fatal("PolicyTypeMarkdown empty")
	}
	if commonepm.StatusMarkdown() == "" {
		t.Fatal("StatusMarkdown empty")
	}
	if commonepm.StatusMarkdownForLeastPrivilege() == "" {
		t.Fatal("StatusMarkdownForLeastPrivilege empty")
	}
	if commonepm.ControlMarkdown() == "" {
		t.Fatal("ControlMarkdown empty")
	}
	if commonepm.DayFilterMarkdown() == "" {
		t.Fatal("DayFilterMarkdown empty")
	}
	var pt commonepm.PolicyTypeValidator
	if pt.Description(ctx) == "" || pt.MarkdownDescription(ctx) == "" {
		t.Fatal("PolicyTypeValidator descriptions")
	}
	var st commonepm.StatusValidator
	if st.Description(ctx) == "" || st.MarkdownDescription(ctx) == "" {
		t.Fatal("StatusValidator descriptions")
	}
	var c commonepm.ControlSetValidator
	if c.Description(ctx) == "" || c.MarkdownDescription(ctx) == "" {
		t.Fatal("ControlSetValidator descriptions")
	}
	var d commonepm.DayFilterSetValidator
	if d.Description(ctx) == "" || d.MarkdownDescription(ctx) == "" {
		t.Fatal("DayFilterSetValidator descriptions")
	}
	var tf commonepm.TimeFilterSetValidator
	if tf.Description(ctx) == "" || tf.MarkdownDescription(ctx) == "" {
		t.Fatal("TimeFilterSetValidator descriptions")
	}
	var df commonepm.DateFilterSetValidator
	if df.Description(ctx) == "" || df.MarkdownDescription(ctx) == "" {
		t.Fatal("DateFilterSetValidator descriptions")
	}
	v := commonepm.CollectionSetValidator{DisplayName: "User group"}
	if v.Description(ctx) == "" || v.MarkdownDescription(ctx) == "" {
		t.Fatal("CollectionSetValidator descriptions")
	}
}
