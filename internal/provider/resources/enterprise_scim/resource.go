// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &EnterpriseScimResource{}
var _ resource.ResourceWithConfigure = &EnterpriseScimResource{}
var _ resource.ResourceWithImportState = &EnterpriseScimResource{}
var _ resource.ResourceWithModifyPlan = &EnterpriseScimResource{}

type EnterpriseScimResource struct {
	utils.BaseResource
}

func (r *EnterpriseScimResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_scim"
}

func (r *EnterpriseScimResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewEnterpriseScimResource() resource.Resource {
	return &EnterpriseScimResource{}
}

// ModifyPlan marks provisioning_token as unknown when prefix or unique_groups is being updated,
// so the update response (which returns a new token in res.message) can be stored without causing
// Terraform to report an inconsistent result.
func (r *EnterpriseScimResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}
	var plan, state EnterpriseScimResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Treat null and empty string as equivalent for prefix (API returns "" when unset; config may use null).
	prefixChanged := !prefixEqual(plan.Prefix, state.Prefix)
	uniqueGroupsChanged := !plan.UniqueGroups.Equal(state.UniqueGroups)
	if prefixChanged || uniqueGroupsChanged {
		plan.ProvisioningToken = types.StringUnknown()
		resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
	}
}

// prefixEqual returns true if the two prefix values are equivalent for change-detection.
// Null and empty string are treated as the same (no prefix).
func prefixEqual(a, b types.String) bool {
	if a.Equal(b) {
		return true
	}
	aEmpty := a.IsNull() || a.ValueString() == ""
	bEmpty := b.IsNull() || b.ValueString() == ""
	return aEmpty && bEmpty
}
