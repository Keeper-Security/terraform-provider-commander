// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescim

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &EnterpriseScimResource{}
var _ resource.ResourceWithConfigure = &EnterpriseScimResource{}
var _ resource.ResourceWithImportState = &EnterpriseScimResource{}

type EnterpriseScimResource struct {
	utils.BaseResource
}

func (r *EnterpriseScimResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_scim"
}

func (r *EnterpriseScimResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewEnterpriseScimResource() resource.Resource {
	return &EnterpriseScimResource{}
}
