// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisepush

import (
	"context"
	"fmt"
	"os"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &EnterprisePushResource{}
var _ resource.ResourceWithConfigure = &EnterprisePushResource{}
var _ resource.ResourceWithModifyPlan = &EnterprisePushResource{}

type EnterprisePushResource struct {
	utils.BaseResource
}

func (r *EnterprisePushResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_push"
}

func (r *EnterprisePushResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewEnterprisePushResource() resource.Resource {
	return &EnterprisePushResource{}
}

// ModifyPlan reads the file at file_path and sets file_content_sha256 in the plan so that
// file-only changes (same path, different content) trigger a replace and re-push.
func (r *EnterprisePushResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// No plan (e.g. destroy): nothing to do.
	if req.Plan.Raw.IsNull() {
		return
	}
	var plan EnterprisePushResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	filePath := plan.FilePath.ValueString()
	if filePath == "" || plan.FilePath.IsUnknown() {
		return
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		// File missing or unreadable at plan time: add warning and leave plan unchanged.
		resp.Diagnostics.AddWarning(
			"Cannot read file for plan",
			fmt.Sprintf("Could not read file at %q during plan: %v. File content change detection will not work until the file is present.", filePath, err),
		)
		return
	}
	hash := contentSHA256Hex(content)
	plan.FileContentSha256 = types.StringValue(hash)
	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}
