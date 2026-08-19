// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration_test

import (
	"context"
	"testing"

	commonrecordsaasconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/saas_configuration"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var customFieldAttrTypes = map[string]attr.Type{
	"type":      types.StringType,
	"label":     types.StringType,
	"value":     types.StringType,
	"sensitive": types.BoolType,
}

func customFieldObject(t *testing.T, fieldType, label, value string) types.Object {
	t.Helper()
	obj, diags := types.ObjectValue(customFieldAttrTypes, map[string]attr.Value{
		"type":      types.StringValue(fieldType),
		"label":     types.StringValue(label),
		"value":     types.StringValue(value),
		"sensitive": types.BoolValue(false),
	})
	if diags.HasError() {
		t.Fatalf("building custom field object: %v", diags)
	}
	return obj
}

func validateCustomList(t *testing.T, elems []attr.Value) validator.ListResponse {
	t.Helper()
	list, diags := types.ListValue(types.ObjectType{AttrTypes: customFieldAttrTypes}, elems)
	if diags.HasError() {
		t.Fatalf("building custom list: %v", diags)
	}

	var resp validator.ListResponse
	req := validator.ListRequest{
		Path:        path.Root("custom"),
		ConfigValue: list,
	}
	commonrecordsaasconfiguration.RequiredSaasTypeCustomFieldValidator().ValidateList(context.Background(), req, &resp)
	return resp
}

func TestRequiredSaasTypeCustomFieldValidator_AcceptsRequiredField(t *testing.T) {
	t.Parallel()

	resp := validateCustomList(t, []attr.Value{
		customFieldObject(t, "text", "SaaS Type", "Okta"),
		customFieldObject(t, "text", "AppName", "Example App"),
	})
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected valid custom list, got %v", resp.Diagnostics)
	}
}

func TestRequiredSaasTypeCustomFieldValidator_RejectsMissingSaaSType(t *testing.T) {
	t.Parallel()

	resp := validateCustomList(t, []attr.Value{
		customFieldObject(t, "text", "AppName", "Example App"),
	})
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error when SaaS Type custom field is missing")
	}
}

func TestRequiredSaasTypeCustomFieldValidator_RejectsWrongLabel(t *testing.T) {
	t.Parallel()

	resp := validateCustomList(t, []attr.Value{
		customFieldObject(t, "text", "Saas Type", "Okta"),
	})
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error when SaaS Type label casing does not match")
	}
}

func TestRequiredSaasTypeCustomFieldValidator_RejectsWrongType(t *testing.T) {
	t.Parallel()

	resp := validateCustomList(t, []attr.Value{
		customFieldObject(t, "secret", "SaaS Type", "Okta"),
	})
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error when SaaS Type field type is not text")
	}
}

func TestRequiredSaasTypeCustomFieldValidator_RejectsNullCustom(t *testing.T) {
	t.Parallel()

	var resp validator.ListResponse
	req := validator.ListRequest{
		Path:        path.Root("custom"),
		ConfigValue: types.ListNull(types.ObjectType{AttrTypes: customFieldAttrTypes}),
	}
	commonrecordsaasconfiguration.RequiredSaasTypeCustomFieldValidator().ValidateList(context.Background(), req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error when custom is null")
	}
}
