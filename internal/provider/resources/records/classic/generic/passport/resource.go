// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package passport

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &PassportResource{}
var _ resource.ResourceWithConfigure = &PassportResource{}
var _ resource.ResourceWithImportState = &PassportResource{}

type PassportResource struct {
	utils.BaseResource
}

func (r *PassportResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_classic_passport"
}

func (r *PassportResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewPassportResource() resource.Resource {
	return &PassportResource{}
}
