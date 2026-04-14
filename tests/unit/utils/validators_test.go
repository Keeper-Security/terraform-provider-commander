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
