// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &EnterpriseUserResource{}
var _ resource.ResourceWithConfigure = &EnterpriseUserResource{}
var _ resource.ResourceWithImportState = &EnterpriseUserResource{}

type EnterpriseUserResource struct {
	utils.BaseResource
}

func (r *EnterpriseUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_user"
}

func (r *EnterpriseUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewEnterpriseUserResource() resource.Resource {
	return &EnterpriseUserResource{}
}
