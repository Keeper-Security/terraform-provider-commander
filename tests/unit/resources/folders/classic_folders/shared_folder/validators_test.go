// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sharedfolder_test

import (
	"context"
	"testing"

	sharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/classic_folders/shared_folder"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExpirationValidator(t *testing.T) {
	ctx := context.Background()
	v := sharedfolder.ExpirationValidator()
	p := path.Root("users").AtMapKey("u").AtName("expiration")

	t.Run("never_ok", func(t *testing.T) {
		var resp validator.StringResponse
		v.ValidateString(ctx, validator.StringRequest{
			Path:        p,
			ConfigValue: types.StringValue("never"),
		}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
	})

	t.Run("never_case_insensitive_ok", func(t *testing.T) {
		var resp validator.StringResponse
		v.ValidateString(ctx, validator.StringRequest{
			Path:        p,
			ConfigValue: types.StringValue("NeVeR"),
		}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
	})

	t.Run("datetime_ok", func(t *testing.T) {
		var resp validator.StringResponse
		v.ValidateString(ctx, validator.StringRequest{
			Path:        p,
			ConfigValue: types.StringValue("2030-06-01T12:30:45"),
		}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
	})

	t.Run("empty_error", func(t *testing.T) {
		var resp validator.StringResponse
		v.ValidateString(ctx, validator.StringRequest{
			Path:        p,
			ConfigValue: types.StringValue(""),
		}, &resp)
		if !resp.Diagnostics.HasError() {
			t.Fatal("expected error for empty expiration")
		}
	})

	t.Run("invalid_error", func(t *testing.T) {
		var resp validator.StringResponse
		v.ValidateString(ctx, validator.StringRequest{
			Path:        p,
			ConfigValue: types.StringValue("30d"),
		}, &resp)
		if !resp.Diagnostics.HasError() {
			t.Fatal("expected error for relative expiration")
		}
	})

	t.Run("null_skipped", func(t *testing.T) {
		var resp validator.StringResponse
		v.ValidateString(ctx, validator.StringRequest{
			Path:        p,
			ConfigValue: types.StringNull(),
		}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
	})
}

func TestUserExpirationManageUsersValidator(t *testing.T) {
	ctx := context.Background()
	v := sharedfolder.UserExpirationManageUsersValidator()
	objPath := path.Root("users").AtMapKey("alice")

	attrTypes := map[string]attr.Type{
		"manage_users":   types.BoolType,
		"manage_records": types.BoolType,
		"expiration":     types.StringType,
	}

	t.Run("datetime_and_manage_users_true_errors", func(t *testing.T) {
		obj := types.ObjectValueMust(attrTypes, map[string]attr.Value{
			"manage_users":   types.BoolValue(true),
			"manage_records": types.BoolValue(false),
			"expiration":     types.StringValue("2030-06-01T12:30:45"),
		})
		var resp validator.ObjectResponse
		v.ValidateObject(ctx, validator.ObjectRequest{
			Path:        objPath,
			ConfigValue: obj,
		}, &resp)
		if !resp.Diagnostics.HasError() {
			t.Fatal("expected error when manage_users true with time-limited expiration")
		}
	})

	t.Run("never_and_manage_users_true_ok", func(t *testing.T) {
		obj := types.ObjectValueMust(attrTypes, map[string]attr.Value{
			"manage_users":   types.BoolValue(true),
			"manage_records": types.BoolValue(true),
			"expiration":     types.StringValue("never"),
		})
		var resp validator.ObjectResponse
		v.ValidateObject(ctx, validator.ObjectRequest{
			Path:        objPath,
			ConfigValue: obj,
		}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
	})

	t.Run("expiration_null_and_manage_users_true_ok", func(t *testing.T) {
		obj := types.ObjectValueMust(attrTypes, map[string]attr.Value{
			"manage_users":   types.BoolValue(true),
			"manage_records": types.BoolValue(false),
			"expiration":     types.StringNull(),
		})
		var resp validator.ObjectResponse
		v.ValidateObject(ctx, validator.ObjectRequest{
			Path:        objPath,
			ConfigValue: obj,
		}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
	})

	t.Run("datetime_and_manage_users_false_ok", func(t *testing.T) {
		obj := types.ObjectValueMust(attrTypes, map[string]attr.Value{
			"manage_users":   types.BoolValue(false),
			"manage_records": types.BoolValue(true),
			"expiration":     types.StringValue("2030-06-01T12:30:45"),
		})
		var resp validator.ObjectResponse
		v.ValidateObject(ctx, validator.ObjectRequest{
			Path:        objPath,
			ConfigValue: obj,
		}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
	})

	t.Run("invalid_expiration_skips_combo", func(t *testing.T) {
		obj := types.ObjectValueMust(attrTypes, map[string]attr.Value{
			"manage_users":   types.BoolValue(true),
			"manage_records": types.BoolValue(false),
			"expiration":     types.StringValue("not-a-datetime"),
		})
		var resp validator.ObjectResponse
		v.ValidateObject(ctx, validator.ObjectRequest{
			Path:        objPath,
			ConfigValue: obj,
		}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("combo validator should not run for invalid expiration; got %v", resp.Diagnostics)
		}
	})
}
