// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils_test

import (
	"strings"
	"testing"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestAppendCustomFieldsUpdate_LabelRenameClearsOldField(t *testing.T) {
	t.Parallel()

	plan := []commonrecordsutils.CustomFieldModel{
		{
			Type:  types.StringValue("text"),
			Label: types.StringValue("AppName updated"),
			Value: types.StringValue("Example SaaS App updated"),
		},
	}
	state := []commonrecordsutils.CustomFieldModel{
		{
			Type:  types.StringValue("text"),
			Label: types.StringValue("AppName"),
			Value: types.StringValue("Example SaaS App"),
		},
	}

	var parts []string
	commonrecordsutils.AppendCustomFieldsUpdate(&parts, plan, state)

	if len(parts) != 2 {
		t.Fatalf("parts = %#v, want clear old + set new", parts)
	}
	if !strings.Contains(parts[0], `'c.text.AppName'=''`) {
		t.Fatalf("first part = %q, want old label cleared", parts[0])
	}
	if !strings.Contains(parts[1], `'c.text.AppName updated'='Example SaaS App updated'`) {
		t.Fatalf("second part = %q, want new label set", parts[1])
	}
}

func TestAppendCustomFieldsUpdate_ValueChangeDoesNotClearField(t *testing.T) {
	t.Parallel()

	plan := []commonrecordsutils.CustomFieldModel{
		{
			Type:  types.StringValue("text"),
			Label: types.StringValue("AppName"),
			Value: types.StringValue("new value"),
		},
	}
	state := []commonrecordsutils.CustomFieldModel{
		{
			Type:  types.StringValue("text"),
			Label: types.StringValue("AppName"),
			Value: types.StringValue("old value"),
		},
	}

	var parts []string
	commonrecordsutils.AppendCustomFieldsUpdate(&parts, plan, state)

	if len(parts) != 1 {
		t.Fatalf("parts = %#v, want only value update", parts)
	}
	if !strings.Contains(parts[0], `'c.text.AppName'='new value'`) {
		t.Fatalf("part = %q, want updated value only", parts[0])
	}
}

func TestAppendCustomFieldsUpdate_RemovedFieldIsCleared(t *testing.T) {
	t.Parallel()

	state := []commonrecordsutils.CustomFieldModel{
		{
			Type:  types.StringValue("text"),
			Label: types.StringValue("AppName"),
			Value: types.StringValue("Example SaaS App"),
		},
	}

	var parts []string
	commonrecordsutils.AppendCustomFieldsUpdate(&parts, nil, state)

	if len(parts) != 1 {
		t.Fatalf("parts = %#v, want removed field cleared", parts)
	}
	if !strings.Contains(parts[0], `'c.text.AppName'=''`) {
		t.Fatalf("part = %q, want removed field cleared", parts[0])
	}
}

func TestAppendCustomFieldsUpdate_UnchangedCustomFields(t *testing.T) {
	t.Parallel()

	custom := []commonrecordsutils.CustomFieldModel{
		{
			Type:  types.StringValue("text"),
			Label: types.StringValue("AppName"),
			Value: types.StringValue("Example SaaS App"),
		},
	}

	var parts []string
	commonrecordsutils.AppendCustomFieldsUpdate(&parts, custom, custom)
	if len(parts) != 0 {
		t.Fatalf("parts = %#v, want no updates", parts)
	}
}
