// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &ContactResource{}
var _ resource.ResourceWithConfigure = &ContactResource{}
var _ resource.ResourceWithImportState = &ContactResource{}

type ContactResource struct {
	utils.BaseResource
}

func (r *ContactResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact"
}

func (r *ContactResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewContactResource() resource.Resource {
	return &ContactResource{}
}
