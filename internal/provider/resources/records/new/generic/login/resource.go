// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package login

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &LoginResource{}
var _ resource.ResourceWithConfigure = &LoginResource{}
var _ resource.ResourceWithImportState = &LoginResource{}

type LoginResource struct {
	utils.BaseResource
}

func (r *LoginResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_login"
}

func (r *LoginResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewLoginResource() resource.Resource {
	return &LoginResource{}
}
