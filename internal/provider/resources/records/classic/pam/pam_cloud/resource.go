// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamcloud

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &PamCloudResource{}
var _ resource.ResourceWithConfigure = &PamCloudResource{}
var _ resource.ResourceWithImportState = &PamCloudResource{}

type PamCloudResource struct {
	utils.BaseResource
}

func (r *PamCloudResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_classic_pam_cloud"
}

func (r *PamCloudResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewPamCloudResource() resource.Resource {
	return &PamCloudResource{}
}
