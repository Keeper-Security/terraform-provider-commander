// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	"context"

	commonpamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_configuration"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &PamConfigurationResource{}
var _ resource.ResourceWithConfigure = &PamConfigurationResource{}
var _ resource.ResourceWithImportState = &PamConfigurationResource{}
var _ resource.ResourceWithConfigValidators = &PamConfigurationResource{}
var _ resource.ResourceWithModifyPlan = &PamConfigurationResource{}

type PamConfigurationResource struct {
	utils.BaseResource
}

func (r *PamConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pam_configuration"
}

func (r *PamConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewPamConfigurationResource() resource.Resource {
	return &PamConfigurationResource{}
}

func (r *PamConfigurationResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		environmentBlocksValidator{},
	}
}

// ModifyPlan rejects adding universal_secrets_sync for the first time with enabled = false.
// This needs prior state, so it can't be a ConfigValidator.
func (r *PamConfigurationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan commonpamconfiguration.PamConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	planUss := ussCommonModel(&plan)
	// Null enabled is sent as --enabled false, so it's rejected the same as an explicit false.
	if planUss == nil || planUss.Enabled.IsUnknown() || planUss.Enabled.ValueBool() {
		return
	}

	if !req.State.Raw.IsNull() {
		var state commonpamconfiguration.PamConfigurationResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() || ussCommonModel(&state) != nil {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(
		path.Root(plan.Environment.ValueString()).AtName("universal_secrets_sync").AtName("enabled"),
		ErrSummaryInvalidConfig,
		ErrMsgUssFirstTimeDisabled,
	)
}
