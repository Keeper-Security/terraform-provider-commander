// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharefolder

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &ShareFolderResource{}
var _ resource.ResourceWithConfigure = &ShareFolderResource{}
var _ resource.ResourceWithImportState = &ShareFolderResource{}

type ShareFolderResource struct {
	utils.BaseResource
}

func (r *ShareFolderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_share_folder"
}

func (r *ShareFolderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewShareFolderResource() resource.Resource {
	return &ShareFolderResource{}
}
