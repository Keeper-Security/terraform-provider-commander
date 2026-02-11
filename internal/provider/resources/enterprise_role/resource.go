// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriserole

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &EnterpriseRoleResource{}
var _ resource.ResourceWithConfigure = &EnterpriseRoleResource{}
var _ resource.ResourceWithImportState = &EnterpriseRoleResource{}

type EnterpriseRoleResource struct {
	utils.BaseResource
}

func (r *EnterpriseRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_role"
}

func (r *EnterpriseRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewEnterpriseRoleResource() resource.Resource {
	return &EnterpriseRoleResource{}
}
