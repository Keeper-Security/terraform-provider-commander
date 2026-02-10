// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &ManageCompanyResource{}
var _ resource.ResourceWithConfigure = &ManageCompanyResource{}
var _ resource.ResourceWithImportState = &ManageCompanyResource{}

type ManageCompanyResource struct {
	utils.BaseResource
}

func (r *ManageCompanyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_manage_company"
}

func (r *ManageCompanyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.BaseResource.ConfigureResource(ctx, req, resp)
}

func NewManageCompanyResource() resource.Resource {
	return &ManageCompanyResource{}
}
