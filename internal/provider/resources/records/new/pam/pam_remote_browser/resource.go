// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamremotebrowser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &PamRemoteBrowserResource{}
var _ resource.ResourceWithConfigure = &PamRemoteBrowserResource{}
var _ resource.ResourceWithImportState = &PamRemoteBrowserResource{}

type PamRemoteBrowserResource struct {
	utils.BaseResource
}

func (r *PamRemoteBrowserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_pam_remote_browser"
}

func (r *PamRemoteBrowserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewPamRemoteBrowserResource() resource.Resource {
	return &PamRemoteBrowserResource{}
}
