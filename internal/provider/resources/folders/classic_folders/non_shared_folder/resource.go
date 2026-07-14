// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package nonsharedfolder

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &NonSharedFolderResource{}
var _ resource.ResourceWithConfigure = &NonSharedFolderResource{}
var _ resource.ResourceWithImportState = &NonSharedFolderResource{}

type NonSharedFolderResource struct {
	utils.BaseResource
}

func (r *NonSharedFolderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_non_shared_folder"
}

func (r *NonSharedFolderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewNonSharedFolderResource() resource.Resource {
	return &NonSharedFolderResource{}
}
