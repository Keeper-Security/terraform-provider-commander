// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy_test

import (
	"context"
	"testing"

	commonepm "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/epm_policy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func mustStringSet(t *testing.T, elems ...attr.Value) types.Set {
	t.Helper()
	if len(elems) == 0 {
		return types.SetNull(types.StringType)
	}
	s, diags := types.SetValue(types.StringType, elems)
	if diags.HasError() {
		t.Fatal(diags)
	}
	return s
}

func TestPolicyTypeValidator(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	p := path.Root("policy_type")
	var v commonepm.PolicyTypeValidator
	run := func(val types.String) *validator.StringResponse {
		resp := &validator.StringResponse{}
		v.ValidateString(ctx, validator.StringRequest{ConfigValue: val, Path: p}, resp)
		return resp
	}
	if run(types.StringUnknown()).Diagnostics.HasError() {
		t.Fatal("unknown")
	}
	if run(types.StringNull()).Diagnostics.HasError() {
		t.Fatal("null")
	}
	if run(types.StringValue("elevation")).Diagnostics.HasError() {
		t.Fatal("valid")
	}
	if !run(types.StringValue("nope")).Diagnostics.HasError() {
		t.Fatal("invalid")
	}
}

func TestStatusValidator(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	p := path.Root("status")
	var v commonepm.StatusValidator
	var unkResp validator.StringResponse
	v.ValidateString(ctx, validator.StringRequest{ConfigValue: types.StringUnknown(), Path: p}, &unkResp)
	if unkResp.Diagnostics.HasError() {
		t.Fatal("unknown")
	}
	var nullResp validator.StringResponse
	v.ValidateString(ctx, validator.StringRequest{ConfigValue: types.StringNull(), Path: p}, &nullResp)
	if nullResp.Diagnostics.HasError() {
		t.Fatal("null")
	}
	resp := &validator.StringResponse{}
	v.ValidateString(ctx, validator.StringRequest{ConfigValue: types.StringValue("bogus"), Path: p}, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestControlSetValidator(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	p := path.Root("control")
	var v commonepm.ControlSetValidator
	run := func(set types.Set) *validator.SetResponse {
		resp := &validator.SetResponse{}
		v.ValidateSet(ctx, validator.SetRequest{ConfigValue: set, Path: p}, resp)
		return resp
	}
	if run(types.SetUnknown(types.StringType)).Diagnostics.HasError() {
		t.Fatal("unknown")
	}
	nullEl, d := types.SetValue(types.StringType, []attr.Value{types.StringNull(), types.StringUnknown(), types.StringValue("")})
	if d.HasError() {
		t.Fatal(d)
	}
	if run(nullEl).Diagnostics.HasError() {
		t.Fatal("null/unknown/empty elements skip")
	}
	if run(mustStringSet(t, types.StringValue("audit"))).Diagnostics.HasError() {
		t.Fatal("valid")
	}
	if !run(mustStringSet(t, types.StringValue("badcontrol"))).Diagnostics.HasError() {
		t.Fatal("invalid value")
	}
}

func TestDayFilterSetValidator(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	p := path.Root("day_filter")
	var v commonepm.DayFilterSetValidator
	resp := &validator.SetResponse{}
	v.ValidateSet(ctx, validator.SetRequest{
		ConfigValue: mustStringSet(t, types.StringValue("monday")),
		Path:        p,
	}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}
	skip, d := types.SetValue(types.StringType, []attr.Value{types.StringNull(), types.StringValue("")})
	if d.HasError() {
		t.Fatal(d)
	}
	resp = &validator.SetResponse{}
	v.ValidateSet(ctx, validator.SetRequest{ConfigValue: skip, Path: p}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}
	resp = &validator.SetResponse{}
	v.ValidateSet(ctx, validator.SetRequest{
		ConfigValue: mustStringSet(t, types.StringValue("Funday")),
		Path:        p,
	}, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("want error")
	}
}

func TestTimeFilterSetValidator(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	p := path.Root("time_filter")
	var v commonepm.TimeFilterSetValidator
	badFmt := &validator.SetResponse{}
	v.ValidateSet(ctx, validator.SetRequest{
		ConfigValue: mustStringSet(t, types.StringValue("25-26")),
		Path:        p,
	}, badFmt)
	if !badFmt.Diagnostics.HasError() {
		t.Fatal("bad format")
	}
	oldClockFmt := &validator.SetResponse{}
	v.ValidateSet(ctx, validator.SetRequest{
		ConfigValue: mustStringSet(t, types.StringValue("09:00-17:00")),
		Path:        p,
	}, oldClockFmt)
	if !oldClockFmt.Diagnostics.HasError() {
		t.Fatal("HH:MM format should be rejected")
	}
	overlap := &validator.SetResponse{}
	v.ValidateSet(ctx, validator.SetRequest{
		ConfigValue: mustStringSet(t,
			types.StringValue("9-12"),
			types.StringValue("11-13"),
		),
		Path: p,
	}, overlap)
	if !overlap.Diagnostics.HasError() {
		t.Fatal("overlap")
	}
	ok := &validator.SetResponse{}
	v.ValidateSet(ctx, validator.SetRequest{
		ConfigValue: mustStringSet(t, types.StringValue("9-12")),
		Path:        p,
	}, ok)
	if ok.Diagnostics.HasError() {
		t.Fatal(ok.Diagnostics)
	}
	skipTime, d := types.SetValue(types.StringType, []attr.Value{types.StringNull(), types.StringUnknown(), types.StringValue("")})
	if d.HasError() {
		t.Fatal(d)
	}
	skipResp := &validator.SetResponse{}
	v.ValidateSet(ctx, validator.SetRequest{ConfigValue: skipTime, Path: p}, skipResp)
	if skipResp.Diagnostics.HasError() {
		t.Fatal(skipResp.Diagnostics)
	}
}

func TestDateFilterSetValidator(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	p := path.Root("date_filter")
	var v commonepm.DateFilterSetValidator
	bad := &validator.SetResponse{}
	v.ValidateSet(ctx, validator.SetRequest{
		ConfigValue: mustStringSet(t, types.StringValue("2025-01-01")),
		Path:        p,
	}, bad)
	if !bad.Diagnostics.HasError() {
		t.Fatal("bad format")
	}
	overlap := &validator.SetResponse{}
	v.ValidateSet(ctx, validator.SetRequest{
		ConfigValue: mustStringSet(t,
			types.StringValue("2025-01-01:2025-01-31"),
			types.StringValue("2025-01-15:2025-02-15"),
		),
		Path: p,
	}, overlap)
	if !overlap.Diagnostics.HasError() {
		t.Fatal("overlap")
	}
	skipDate, d := types.SetValue(types.StringType, []attr.Value{types.StringNull(), types.StringUnknown()})
	if d.HasError() {
		t.Fatal(d)
	}
	skipResp := &validator.SetResponse{}
	v.ValidateSet(ctx, validator.SetRequest{ConfigValue: skipDate, Path: p}, skipResp)
	if skipResp.Diagnostics.HasError() {
		t.Fatal(skipResp.Diagnostics)
	}
}

func TestCollectionSetValidator(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	p := path.Root("user_groups")
	v := commonepm.CollectionSetValidator{DisplayName: "User group"}
	mix := &validator.SetResponse{}
	v.ValidateSet(ctx, validator.SetRequest{
		ConfigValue: mustStringSet(t, types.StringValue("*"), types.StringValue("id1")),
		Path:        p,
	}, mix)
	if !mix.Diagnostics.HasError() {
		t.Fatal("mix wildcard")
	}
	emptyStr := &validator.SetResponse{}
	v.ValidateSet(ctx, validator.SetRequest{
		ConfigValue: mustStringSet(t, types.StringValue("  ")),
		Path:        p,
	}, emptyStr)
	if !emptyStr.Diagnostics.HasError() {
		t.Fatal("empty string element")
	}
	skip, d := types.SetValue(types.StringType, []attr.Value{types.StringNull(), types.StringUnknown()})
	if d.HasError() {
		t.Fatal(d)
	}
	skipResp := &validator.SetResponse{}
	v.ValidateSet(ctx, validator.SetRequest{ConfigValue: skip, Path: p}, skipResp)
	if skipResp.Diagnostics.HasError() {
		t.Fatal(skipResp.Diagnostics)
	}
}

func TestIsPolicyTypeAndCollectionsHelpers(t *testing.T) {
	t.Parallel()
	if !commonepm.IsPolicyTypeLeastPrivilege("  Least_Privilege ") {
		t.Fatal()
	}
	if !commonepm.IsPolicyTypeCommand("command") {
		t.Fatal()
	}
	emptyList := types.ListNull(types.StringType)
	if !commonepm.IsListEmptyOrNull(emptyList) {
		t.Fatal("null list")
	}
	unkList := types.ListUnknown(types.StringType)
	if !commonepm.IsListEmptyOrNull(unkList) {
		t.Fatal("unknown list")
	}
	blankOnly, d := types.ListValue(types.StringType, []attr.Value{types.StringValue("  "), types.StringNull()})
	if d.HasError() {
		t.Fatal(d)
	}
	if !commonepm.IsListEmptyOrNull(blankOnly) {
		t.Fatal("blank-only list")
	}
	withVal, d := types.ListValue(types.StringType, []attr.Value{types.StringValue("id")})
	if d.HasError() {
		t.Fatal(d)
	}
	if commonepm.IsListEmptyOrNull(withVal) {
		t.Fatal("non-empty list")
	}
	if commonepm.IsStringEmptyOrNull(types.StringValue("x")) {
		t.Fatal()
	}
	if !commonepm.IsStringEmptyOrNull(types.StringNull()) || !commonepm.IsStringEmptyOrNull(types.StringUnknown()) {
		t.Fatal("null/unknown string")
	}
	if !commonepm.IsStringEmptyOrNull(types.StringValue("  ")) {
		t.Fatal()
	}
	if !commonepm.IsSetEmptyOrNull(types.SetNull(types.StringType)) {
		t.Fatal()
	}
	emptyExplicit, d := types.SetValue(types.StringType, []attr.Value{})
	if d.HasError() {
		t.Fatal(d)
	}
	if !commonepm.IsSetPresent(emptyExplicit) {
		t.Fatal("explicit empty set should count as present")
	}
	s := mustStringSet(t, types.StringValue("a"))
	if !commonepm.IsSetPresent(s) {
		t.Fatal("present")
	}
	if !commonepm.HasAtLeastOneCollection(s, types.SetNull(types.StringType), types.SetNull(types.StringType)) {
		t.Fatal("one collection")
	}
	if commonepm.HasAtLeastOneMachineAndUser(s, types.SetNull(types.StringType)) {
		t.Fatal("need both")
	}
	if commonepm.HasAllThreeCollections(s, s, types.SetNull(types.StringType)) {
		t.Fatal("need three")
	}
}

func TestIsStatusAllowedForLeastPrivilege(t *testing.T) {
	t.Parallel()
	if !commonepm.IsStatusAllowedForLeastPrivilege("enforce") || !commonepm.IsStatusAllowedForLeastPrivilege("off") {
		t.Fatal()
	}
	if commonepm.IsStatusAllowedForLeastPrivilege("monitor") {
		t.Fatal()
	}
}
