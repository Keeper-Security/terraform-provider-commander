// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package classicsharedfolder

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &ClassicSharedFolderResource{}
var _ resource.ResourceWithConfigure = &ClassicSharedFolderResource{}
var _ resource.ResourceWithImportState = &ClassicSharedFolderResource{}

type ClassicSharedFolderResource struct {
	utils.BaseResource
}

func (r *ClassicSharedFolderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shared_folder"
}

func (r *ClassicSharedFolderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewClassicSharedFolderResource() resource.Resource {
	return &ClassicSharedFolderResource{}
}
