// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamuser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &PamUserResource{}
var _ resource.ResourceWithConfigure = &PamUserResource{}
var _ resource.ResourceWithImportState = &PamUserResource{}

type PamUserResource struct {
	utils.BaseResource
}

func (r *PamUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_pam_user"
}

func (r *PamUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewPamUserResource() resource.Resource {
	return &PamUserResource{}
}
