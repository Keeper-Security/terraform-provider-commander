// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &EnterpriseNodeResource{}
var _ resource.ResourceWithConfigure = &EnterpriseNodeResource{}
var _ resource.ResourceWithImportState = &EnterpriseNodeResource{}

type EnterpriseNodeResource struct {
	utils.BaseResource
}

func (r *EnterpriseNodeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_node"
}

func (r *EnterpriseNodeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.BaseResource.ConfigureResource(ctx, req, resp)
}

func NewEnterpriseNodeResource() resource.Resource {
	return &EnterpriseNodeResource{}
}
