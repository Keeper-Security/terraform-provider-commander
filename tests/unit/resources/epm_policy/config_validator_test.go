// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package epmpolicy_test

import (
	"context"
	"testing"

	epmpolicy "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/epm_policy"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func epmPolicyResourceWithConfigValidators(t *testing.T) resource.ResourceWithConfigValidators {
	t.Helper()
	res := epmpolicy.NewEpmPolicyResource()
	withCV, ok := res.(resource.ResourceWithConfigValidators)
	if !ok {
		t.Fatalf("expected ResourceWithConfigValidators, got %T", res)
	}
	return withCV
}

func TestEpmPolicyConfigValidator_Descriptions(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	withCV := epmPolicyResourceWithConfigValidators(t)
	for i, v := range withCV.ConfigValidators(ctx) {
		if v.Description(ctx) == "" || v.MarkdownDescription(ctx) == "" {
			t.Fatalf("validator %d: empty description", i)
		}
	}
}

func TestEpmPolicyConfigValidator_ValidateResource(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	withCV := epmPolicyResourceWithConfigValidators(t)
	validators := withCV.ConfigValidators(ctx)
	if len(validators) != 1 {
		t.Fatalf("expected 1 config validator, got %d", len(validators))
	}

	r := epmpolicy.NewEpmPolicyResource().(*epmpolicy.EpmPolicyResource)
	var schResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schResp)
	sch := schResp.Schema

	objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"id":                  tftypes.String,
		"managed_company":     tftypes.String,
		"policy_name":         tftypes.String,
		"policy_type":         tftypes.String,
		"status":              tftypes.String,
		"control":             tftypes.Set{ElementType: tftypes.String},
		"user_groups":         tftypes.Set{ElementType: tftypes.String},
		"machine_collections": tftypes.Set{ElementType: tftypes.String},
		"applications":        tftypes.Set{ElementType: tftypes.String},
		"day_filter":          tftypes.Set{ElementType: tftypes.String},
		"time_filter":         tftypes.Set{ElementType: tftypes.String},
		"date_filter":         tftypes.Set{ElementType: tftypes.String},
	}}

	nullSet := tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil)
	controlSet := tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
		tftypes.NewValue(tftypes.String, "audit"),
	})
	raw := tftypes.NewValue(objType, map[string]tftypes.Value{
		"id":                  tftypes.NewValue(tftypes.String, nil),
		"managed_company":     tftypes.NewValue(tftypes.String, nil),
		"policy_name":         tftypes.NewValue(tftypes.String, "n"),
		"policy_type":         tftypes.NewValue(tftypes.String, "least_privilege"),
		"status":              tftypes.NewValue(tftypes.String, "enforce"),
		"control":             controlSet,
		"user_groups":         nullSet,
		"machine_collections": nullSet,
		"applications":        nullSet,
		"day_filter":          nullSet,
		"time_filter":         nullSet,
		"date_filter":         nullSet,
	})

	req := resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: sch, Raw: raw}}
	var resp resource.ValidateConfigResponse
	validators[0].ValidateResource(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("least_privilege + control must fail validation")
	}
}

func TestEpmPolicyConfigValidator_ValidateResource_ConfigGetError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	withCV := epmPolicyResourceWithConfigValidators(t)
	validators := withCV.ConfigValidators(ctx)

	r := epmpolicy.NewEpmPolicyResource().(*epmpolicy.EpmPolicyResource)
	var schResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schResp)
	sch := schResp.Schema
	objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"id": tftypes.String, "managed_company": tftypes.String, "policy_name": tftypes.String,
		"policy_type": tftypes.String, "status": tftypes.String,
		"control": tftypes.Set{ElementType: tftypes.String}, "user_groups": tftypes.Set{ElementType: tftypes.String},
		"machine_collections": tftypes.Set{ElementType: tftypes.String}, "applications": tftypes.Set{ElementType: tftypes.String},
		"day_filter": tftypes.Set{ElementType: tftypes.String}, "time_filter": tftypes.Set{ElementType: tftypes.String},
		"date_filter": tftypes.Set{ElementType: tftypes.String},
	}}
	req := resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: sch, Raw: tftypes.NewValue(objType, nil)}}
	var resp resource.ValidateConfigResponse
	validators[0].ValidateResource(ctx, req, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("null config should fail Get")
	}
}
