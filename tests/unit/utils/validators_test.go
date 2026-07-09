// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils_test

import (
	"context"
	"testing"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNodeValidator_Description(t *testing.T) {
	v := utils.NodeValidator
	ctx := context.Background()
	if v.Description(ctx) != "Node Name must be at least 1 character(s) long." {
		t.Errorf("unexpected Description: %s", v.Description(ctx))
	}
	if v.MarkdownDescription(ctx) != "Node Name must be at least 1 character(s) long." {
		t.Errorf("unexpected MarkdownDescription: %s", v.MarkdownDescription(ctx))
	}
}

func TestNodeValidator_ValidateString_Valid(t *testing.T) {
	v := utils.NodeValidator
	ctx := context.Background()
	req := validator.StringRequest{ConfigValue: types.StringValue("node1")}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics for valid node name")
	}
}

func TestNodeValidator_ValidateString_Empty(t *testing.T) {
	v := utils.NodeValidator
	ctx := context.Background()
	req := validator.StringRequest{ConfigValue: types.StringValue("")}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for empty node name")
	}
}

func TestNodeValidator_ValidateString_WhitespaceOnly(t *testing.T) {
	v := utils.NodeValidator
	ctx := context.Background()
	req := validator.StringRequest{ConfigValue: types.StringValue("   \t  ")}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for whitespace-only node name")
	}
}

func TestNodeValidator_ValidateString_Null(t *testing.T) {
	v := utils.NodeValidator
	ctx := context.Background()
	req := validator.StringRequest{ConfigValue: types.StringNull()}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics for null (skip validation)")
	}
}

func TestNodeValidator_ValidateString_Unknown(t *testing.T) {
	v := utils.NodeValidator
	ctx := context.Background()
	req := validator.StringRequest{ConfigValue: types.StringUnknown()}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics for unknown (skip validation)")
	}
}

func TestManagedCompanyValidator_Description(t *testing.T) {
	v := utils.ManagedCompanyValidator
	ctx := context.Background()
	if v.Description(ctx) != "Managed Company Name must be at least 1 character(s) long." {
		t.Errorf("unexpected Description: %s", v.Description(ctx))
	}
}

func TestManagedCompanyValidator_ValidateString_Valid(t *testing.T) {
	v := utils.ManagedCompanyValidator
	ctx := context.Background()
	req := validator.StringRequest{ConfigValue: types.StringValue("Company A")}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics for valid managed company name")
	}
}

func TestManagedCompanyValidator_ValidateString_Empty(t *testing.T) {
	v := utils.ManagedCompanyValidator
	ctx := context.Background()
	req := validator.StringRequest{ConfigValue: types.StringValue("")}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for empty managed company name")
	}
}

func TestManagedCompanyValidator_ValidateString_WhitespaceOnly(t *testing.T) {
	v := utils.ManagedCompanyValidator
	ctx := context.Background()
	req := validator.StringRequest{ConfigValue: types.StringValue("     ")}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for whitespace-only managed company name")
	}
}

func TestManagedCompanyValidator_ValidateString_Null(t *testing.T) {
	v := utils.ManagedCompanyValidator
	ctx := context.Background()
	req := validator.StringRequest{ConfigValue: types.StringNull()}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics for null")
	}
}

func TestTeamsValidator_Description(t *testing.T) {
	v := utils.TeamsValidator
	ctx := context.Background()
	if v.Description(ctx) != "Team set must not contain empty or whitespace-only strings." {
		t.Errorf("unexpected Description: %s", v.Description(ctx))
	}
}

func TestTeamsValidator_ValidateSet_Valid(t *testing.T) {
	v := utils.TeamsValidator
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue("team1"), types.StringValue("team2")})
	req := validator.SetRequest{ConfigValue: setVal}
	var resp validator.SetResponse
	v.ValidateSet(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics for valid teams set")
	}
}

func TestTeamsValidator_ValidateSet_EmptyString(t *testing.T) {
	v := utils.TeamsValidator
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue(""), types.StringValue("team1")})
	req := validator.SetRequest{ConfigValue: setVal}
	var resp validator.SetResponse
	v.ValidateSet(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for set containing empty string")
	}
}

func TestTeamsValidator_ValidateSet_WhitespaceOnlyString(t *testing.T) {
	v := utils.TeamsValidator
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue("  \t  "), types.StringValue("team1")})
	req := validator.SetRequest{ConfigValue: setVal}
	var resp validator.SetResponse
	v.ValidateSet(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for set containing whitespace-only string")
	}
}

func TestTeamsValidator_ValidateSet_Null(t *testing.T) {
	v := utils.TeamsValidator
	ctx := context.Background()
	req := validator.SetRequest{ConfigValue: types.SetNull(types.StringType)}
	var resp validator.SetResponse
	v.ValidateSet(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics for null set")
	}
}

func TestTeamsValidator_ValidateSet_Unknown(t *testing.T) {
	v := utils.TeamsValidator
	ctx := context.Background()
	req := validator.SetRequest{ConfigValue: types.SetUnknown(types.StringType)}
	var resp validator.SetResponse
	v.ValidateSet(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics for unknown set")
	}
}

func TestRolesValidator_Description(t *testing.T) {
	v := utils.RolesValidator
	ctx := context.Background()
	if v.Description(ctx) != "Role set must not contain empty or whitespace-only strings." {
		t.Errorf("unexpected Description: %s", v.Description(ctx))
	}
}

func TestRolesValidator_ValidateSet_Valid(t *testing.T) {
	v := utils.RolesValidator
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue("role1"), types.StringValue("role2")})
	req := validator.SetRequest{ConfigValue: setVal}
	var resp validator.SetResponse
	v.ValidateSet(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics for valid roles set")
	}
}

func TestRolesValidator_ValidateSet_EmptyString(t *testing.T) {
	v := utils.RolesValidator
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue("role1"), types.StringValue("")})
	req := validator.SetRequest{ConfigValue: setVal}
	var resp validator.SetResponse
	v.ValidateSet(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for set containing empty string")
	}
}

func TestRolesValidator_ValidateSet_WhitespaceOnlyString(t *testing.T) {
	v := utils.RolesValidator
	ctx := context.Background()
	setVal, _ := types.SetValueFrom(ctx, types.StringType, []types.String{types.StringValue("role1"), types.StringValue(" \n ")})
	req := validator.SetRequest{ConfigValue: setVal}
	var resp validator.SetResponse
	v.ValidateSet(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for set containing whitespace-only string")
	}
}

func TestRolesValidator_ValidateSet_Null(t *testing.T) {
	v := utils.RolesValidator
	ctx := context.Background()
	req := validator.SetRequest{ConfigValue: types.SetNull(types.StringType)}
	var resp validator.SetResponse
	v.ValidateSet(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics for null set")
	}
}

// Enterprise node schema uses StringMinLengthValidator("Enterprise Node Name", 1, false) and ("Enterprise Node Parent Name", 1, true).

func TestStringMinLengthValidator_EnterpriseNodeName_Description(t *testing.T) {
	v := utils.StringMinLengthValidator("Enterprise Node Name", 1, false)
	ctx := context.Background()
	if v.Description(ctx) != "Enterprise Node Name must be at least 1 character(s) long." {
		t.Errorf("unexpected Description: %s", v.Description(ctx))
	}
}

func TestStringMinLengthValidator_EnterpriseNodeName_Valid(t *testing.T) {
	v := utils.StringMinLengthValidator("Enterprise Node Name", 1, false)
	ctx := context.Background()
	req := validator.StringRequest{ConfigValue: types.StringValue("Node1")}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics for valid name")
	}
}

func TestStringMinLengthValidator_EnterpriseNodeName_Empty(t *testing.T) {
	v := utils.StringMinLengthValidator("Enterprise Node Name", 1, false)
	ctx := context.Background()
	req := validator.StringRequest{ConfigValue: types.StringValue("")}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for empty name")
	}
}

func TestStringMinLengthValidator_EnterpriseNodeName_WhitespaceOnly(t *testing.T) {
	v := utils.StringMinLengthValidator("Enterprise Node Name", 1, false)
	ctx := context.Background()
	// Without TrimSpace, len("   ") >= 1 would incorrectly pass.
	req := validator.StringRequest{ConfigValue: types.StringValue("   ")}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for whitespace-only name")
	}
}

func TestStringMinLengthValidator_EnterpriseNodeParent_Description(t *testing.T) {
	v := utils.StringMinLengthValidator("Enterprise Node Parent Name", 1, true)
	ctx := context.Background()
	if v.Description(ctx) != "Enterprise Node Parent Name must be at least 1 character(s) long." {
		t.Errorf("unexpected Description: %s", v.Description(ctx))
	}
}

func TestStringMinLengthValidator_EnterpriseNodeParent_NullAllowed(t *testing.T) {
	v := utils.StringMinLengthValidator("Enterprise Node Parent Name", 1, true)
	ctx := context.Background()
	req := validator.StringRequest{ConfigValue: types.StringNull()}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Error("expected no diagnostics for null (optional parent)")
	}
}

func TestStringMinLengthValidator_EnterpriseNodeParent_EmptyFails(t *testing.T) {
	v := utils.StringMinLengthValidator("Enterprise Node Parent Name", 1, true)
	ctx := context.Background()
	req := validator.StringRequest{ConfigValue: types.StringValue("")}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for empty parent")
	}
}

func TestStringMinLengthValidator_EnterpriseNodeParent_WhitespaceOnlyFails(t *testing.T) {
	v := utils.StringMinLengthValidator("Enterprise Node Parent Name", 1, true)
	ctx := context.Background()
	req := validator.StringRequest{ConfigValue: types.StringValue("\n  \t  ")}
	var resp validator.StringResponse
	v.ValidateString(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for whitespace-only parent")
	}
}

func TestMapNonEmptyValidator_Description(t *testing.T) {
	v := utils.MapNonEmptyValidator("Share User UID or Email")
	ctx := context.Background()
	want := "Share User UID or Email must have at least one entry; omit the block instead of setting it to an empty map."
	if v.Description(ctx) != want {
		t.Errorf("unexpected Description: %s", v.Description(ctx))
	}
	if v.MarkdownDescription(ctx) != v.Description(ctx) {
		t.Error("expected MarkdownDescription to equal Description")
	}
}

func TestMapNonEmptyValidator_NullOrUnknown(t *testing.T) {
	v := utils.MapNonEmptyValidator("Share User UID or Email")
	ctx := context.Background()
	for _, m := range []types.Map{types.MapNull(types.StringType), types.MapUnknown(types.StringType)} {
		req := validator.MapRequest{ConfigValue: m, Path: path.Root("share")}
		var resp validator.MapResponse
		v.ValidateMap(ctx, req, &resp)
		if resp.Diagnostics.HasError() {
			t.Errorf("expected no diagnostics for null/unknown, got: %v", resp.Diagnostics)
		}
	}
}

func TestMapNonEmptyValidator_EmptyMapRejected(t *testing.T) {
	v := utils.MapNonEmptyValidator("Share User UID or Email")
	ctx := context.Background()
	m, diags := types.MapValueFrom(ctx, types.StringType, map[string]string{})
	if diags.HasError() {
		t.Fatal(diags)
	}
	req := validator.MapRequest{ConfigValue: m, Path: path.Root("share")}
	var resp validator.MapResponse
	v.ValidateMap(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for explicit empty map ({})")
	}
}

func TestMapNonEmptyValidator_NonEmptyMapAccepted(t *testing.T) {
	v := utils.MapNonEmptyValidator("Share User UID or Email")
	ctx := context.Background()
	m, diags := types.MapValueFrom(ctx, types.StringType, map[string]string{"alice@example.com": "viewer"})
	if diags.HasError() {
		t.Fatal(diags)
	}
	req := validator.MapRequest{ConfigValue: m, Path: path.Root("share")}
	var resp validator.MapResponse
	v.ValidateMap(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("expected no diagnostics for non-empty map, got: %v", resp.Diagnostics)
	}
}

func TestMapKeysMinLengthValidator_WhitespaceOnlyKey(t *testing.T) {
	v := utils.MapKeysMinLengthValidator("managing node name", 1)
	ctx := context.Background()
	m, diags := types.MapValueFrom(ctx, types.StringType, map[string]string{" \t ": "v"})
	if diags.HasError() {
		t.Fatal(diags)
	}
	req := validator.MapRequest{ConfigValue: m, Path: path.Root("test")}
	var resp validator.MapResponse
	v.ValidateMap(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for whitespace-only map key")
	}
}

func TestMapKeysEmailValidator_Description(t *testing.T) {
	v := utils.MapKeysEmailValidator("Share User Email")
	ctx := context.Background()
	if v.Description(ctx) != "All Share User Email (map keys) must be non-empty email addresses (e.g. user@example.com)." {
		t.Errorf("unexpected Description: %s", v.Description(ctx))
	}
	if v.MarkdownDescription(ctx) != v.Description(ctx) {
		t.Error("expected MarkdownDescription to equal Description")
	}
}

func TestMapKeysEmailValidator_NullOrUnknown(t *testing.T) {
	v := utils.MapKeysEmailValidator("Share User Email")
	ctx := context.Background()
	for _, m := range []types.Map{types.MapNull(types.StringType), types.MapUnknown(types.StringType)} {
		req := validator.MapRequest{ConfigValue: m, Path: path.Root("share")}
		var resp validator.MapResponse
		v.ValidateMap(ctx, req, &resp)
		if resp.Diagnostics.HasError() {
			t.Errorf("expected no diagnostics for null/unknown, got: %v", resp.Diagnostics)
		}
	}
}

func TestMapKeysEmailValidator_ValidEmails(t *testing.T) {
	v := utils.MapKeysEmailValidator("Share User Email")
	ctx := context.Background()
	in := map[string]string{
		"user@example.com":           "viewer",
		"first.last@sub.example.com": "share-manager",
		"a+tag@host.co":              "content-manager",
		"single.char@x.io":           "full-manager",
	}
	m, diags := types.MapValueFrom(ctx, types.StringType, in)
	if diags.HasError() {
		t.Fatal(diags)
	}
	req := validator.MapRequest{ConfigValue: m, Path: path.Root("share")}
	var resp validator.MapResponse
	v.ValidateMap(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("expected no diagnostics for valid emails, got: %v", resp.Diagnostics)
	}
}

func TestMapKeysEmailValidator_InvalidEmails(t *testing.T) {
	v := utils.MapKeysEmailValidator("Share User Email")
	ctx := context.Background()
	cases := []string{
		"",
		"  ",
		"no-at-sign",
		"@no-local.com",
		"two@@example.com",
		"missing-domain-dot@example",
	}
	for _, key := range cases {
		t.Run(key, func(t *testing.T) {
			m, diags := types.MapValueFrom(ctx, types.StringType, map[string]string{key: "viewer"})
			if diags.HasError() {
				t.Fatal(diags)
			}
			req := validator.MapRequest{ConfigValue: m, Path: path.Root("share")}
			var resp validator.MapResponse
			v.ValidateMap(ctx, req, &resp)
			if !resp.Diagnostics.HasError() {
				t.Errorf("expected diagnostics for invalid email key %q", key)
			}
		})
	}
}

func TestMapValuesStringOneOfValidator_Description(t *testing.T) {
	v := utils.MapValuesStringOneOfValidator("Share Permission", []string{"viewer", "share-manager"})
	ctx := context.Background()
	want := "Share Permission (map values) must each be one of: viewer, share-manager."
	if v.Description(ctx) != want {
		t.Errorf("Description = %q, want %q", v.Description(ctx), want)
	}
	if v.MarkdownDescription(ctx) != want {
		t.Errorf("MarkdownDescription = %q, want %q", v.MarkdownDescription(ctx), want)
	}
}

func TestMapValuesStringOneOfValidator_AllAllowed(t *testing.T) {
	allowed := []string{"viewer", "share-manager", "content-manager", "content-share-manager", "full-manager"}
	v := utils.MapValuesStringOneOfValidator("Share Permission", allowed)
	ctx := context.Background()
	in := map[string]string{
		"a@x.com": "viewer",
		"b@x.com": "share-manager",
		"c@x.com": "content-manager",
		"d@x.com": "content-share-manager",
		"e@x.com": "full-manager",
	}
	m, diags := types.MapValueFrom(ctx, types.StringType, in)
	if diags.HasError() {
		t.Fatal(diags)
	}
	req := validator.MapRequest{ConfigValue: m, Path: path.Root("share")}
	var resp validator.MapResponse
	v.ValidateMap(ctx, req, &resp)
	if resp.Diagnostics.HasError() {
		t.Errorf("expected no diagnostics for all-allowed values, got: %v", resp.Diagnostics)
	}
}

func TestMapValuesStringOneOfValidator_Disallowed(t *testing.T) {
	allowed := []string{"viewer", "share-manager"}
	v := utils.MapValuesStringOneOfValidator("Share Permission", allowed)
	ctx := context.Background()
	in := map[string]string{
		"a@x.com": "viewer",
		"b@x.com": "owner", // disallowed
	}
	m, diags := types.MapValueFrom(ctx, types.StringType, in)
	if diags.HasError() {
		t.Fatal(diags)
	}
	req := validator.MapRequest{ConfigValue: m, Path: path.Root("share")}
	var resp validator.MapResponse
	v.ValidateMap(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected diagnostics for disallowed value")
	}
}

func TestMapValuesStringOneOfValidator_NullOrUnknown(t *testing.T) {
	v := utils.MapValuesStringOneOfValidator("Share Permission", []string{"viewer"})
	ctx := context.Background()
	for _, m := range []types.Map{types.MapNull(types.StringType), types.MapUnknown(types.StringType)} {
		req := validator.MapRequest{ConfigValue: m, Path: path.Root("share")}
		var resp validator.MapResponse
		v.ValidateMap(ctx, req, &resp)
		if resp.Diagnostics.HasError() {
			t.Errorf("expected no diagnostics for null/unknown map, got: %v", resp.Diagnostics)
		}
	}
}

func TestNumericStringValidator_Valid(t *testing.T) {
	v := utils.NumericStringValidator("Port", true)
	ctx := context.Background()
	for _, value := range []string{"1", "22", "8080"} {
		req := validator.StringRequest{
			ConfigValue: types.StringValue(value),
			Path:        path.Root("port"),
		}
		var resp validator.StringResponse
		v.ValidateString(ctx, req, &resp)
		if resp.Diagnostics.HasError() {
			t.Errorf("expected no diagnostics for %q, got: %v", value, resp.Diagnostics)
		}
	}
}

func TestNumericStringValidator_Invalid(t *testing.T) {
	v := utils.NumericStringValidator("Port", true)
	ctx := context.Background()
	for _, value := range []string{"", "abc", "22a", " 22", "22 ", "22.5"} {
		req := validator.StringRequest{
			ConfigValue: types.StringValue(value),
			Path:        path.Root("port"),
		}
		var resp validator.StringResponse
		v.ValidateString(ctx, req, &resp)
		if !resp.Diagnostics.HasError() {
			t.Errorf("expected diagnostics for %q", value)
		}
	}
}

func TestNumericStringValidator_NullOrUnknown(t *testing.T) {
	v := utils.NumericStringValidator("Port", true)
	ctx := context.Background()
	for _, value := range []types.String{types.StringNull(), types.StringUnknown()} {
		req := validator.StringRequest{ConfigValue: value, Path: path.Root("port")}
		var resp validator.StringResponse
		v.ValidateString(ctx, req, &resp)
		if resp.Diagnostics.HasError() {
			t.Errorf("expected no diagnostics for null/unknown, got: %v", resp.Diagnostics)
		}
	}
}

func TestDateStringValidator_Valid(t *testing.T) {
	v := utils.DateStringValidator("Expiration date", true)
	ctx := context.Background()
	for _, value := range []string{"2026-07-09", "2026-05-20", "2000-01-01"} {
		req := validator.StringRequest{
			ConfigValue: types.StringValue(value),
			Path:        path.Root("expiration_date"),
		}
		var resp validator.StringResponse
		v.ValidateString(ctx, req, &resp)
		if resp.Diagnostics.HasError() {
			t.Errorf("expected no diagnostics for %q, got: %v", value, resp.Diagnostics)
		}
	}
}

func TestDateStringValidator_Invalid(t *testing.T) {
	v := utils.DateStringValidator("Expiration date", true)
	ctx := context.Background()
	for _, value := range []string{"", "2026-7-9", "07/09/2026", "2026-07-09T00:00:00Z", "2026-13-01", "2026-02-30", "not-a-date"} {
		req := validator.StringRequest{
			ConfigValue: types.StringValue(value),
			Path:        path.Root("expiration_date"),
		}
		var resp validator.StringResponse
		v.ValidateString(ctx, req, &resp)
		if !resp.Diagnostics.HasError() {
			t.Errorf("expected diagnostics for %q", value)
		}
	}
}

func TestDateStringValidator_NullOrUnknown(t *testing.T) {
	v := utils.DateStringValidator("Expiration date", true)
	ctx := context.Background()
	for _, value := range []types.String{types.StringNull(), types.StringUnknown()} {
		req := validator.StringRequest{ConfigValue: value, Path: path.Root("expiration_date")}
		var resp validator.StringResponse
		v.ValidateString(ctx, req, &resp)
		if resp.Diagnostics.HasError() {
			t.Errorf("expected no diagnostics for null/unknown, got: %v", resp.Diagnostics)
		}
	}
}
