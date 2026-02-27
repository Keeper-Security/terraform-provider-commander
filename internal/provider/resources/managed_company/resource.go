// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcompany

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &ManagedCompanyResource{}
var _ resource.ResourceWithConfigure = &ManagedCompanyResource{}
var _ resource.ResourceWithImportState = &ManagedCompanyResource{}

type ManagedCompanyResource struct {
	utils.BaseResource
}

func (r *ManagedCompanyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_company"
}

func (r *ManagedCompanyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewManagedCompanyResource() resource.Resource {
	return &ManagedCompanyResource{}
}
