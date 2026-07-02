// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamremotebrowser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &PamRemoteBrowserResource{}
var _ resource.ResourceWithConfigure = &PamRemoteBrowserResource{}

type PamRemoteBrowserResource struct {
	utils.BaseResource
}

func (r *PamRemoteBrowserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_classic_pam_remote_browser"
}

func (r *PamRemoteBrowserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewPamRemoteBrowserResource() resource.Resource {
	return &PamRemoteBrowserResource{}
}
