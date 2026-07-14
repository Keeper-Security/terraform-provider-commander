// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newfolder

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &NewFolderResource{}
var _ resource.ResourceWithConfigure = &NewFolderResource{}
var _ resource.ResourceWithImportState = &NewFolderResource{}

type NewFolderResource struct {
	utils.BaseResource
}

func (r *NewFolderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_folder"
}

func (r *NewFolderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewNewFolderResource() resource.Resource {
	return &NewFolderResource{}
}
