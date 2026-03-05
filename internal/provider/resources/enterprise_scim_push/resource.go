// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &EnterpriseScimPushResource{}
var _ resource.ResourceWithConfigure = &EnterpriseScimPushResource{}

type EnterpriseScimPushResource struct {
	utils.BaseResource
}

func (r *EnterpriseScimPushResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_scim_push"
}

func (r *EnterpriseScimPushResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewEnterpriseScimPushResource() resource.Resource {
	return &EnterpriseScimPushResource{}
}
