// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	"context"

	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_directory"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_records"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &PamDirectoryResource{}
var _ resource.ResourceWithConfigure = &PamDirectoryResource{}
var _ resource.ResourceWithModifyPlan = &PamDirectoryResource{}

type PamDirectoryResource struct {
	utils.BaseResource
}

func (r *PamDirectoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pam_directory"
}

func (r *PamDirectoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func (r *PamDirectoryResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var plan commonpamdirectory.PamDirectoryResourceModel
	var state commonpamdirectory.PamDirectoryResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(commonpamrecords.ValidatePamSettingsFieldsNotRemoved(plan.PamSettings, state.PamSettings)...)
}

func NewPamDirectoryResource() resource.Resource {
	return &PamDirectoryResource{}
}
