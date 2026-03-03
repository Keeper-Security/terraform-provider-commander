// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharedfolder

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &SharedFolderResource{}
var _ resource.ResourceWithConfigure = &SharedFolderResource{}
var _ resource.ResourceWithImportState = &SharedFolderResource{}

type SharedFolderResource struct {
	utils.BaseResource
}

func (r *SharedFolderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shared_folder"
}

func (r *SharedFolderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewSharedFolderResource() resource.Resource {
	return &SharedFolderResource{}
}
