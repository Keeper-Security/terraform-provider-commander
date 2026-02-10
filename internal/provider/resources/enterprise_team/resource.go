// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &EnterpriseTeamResource{}
var _ resource.ResourceWithConfigure = &EnterpriseTeamResource{}
var _ resource.ResourceWithImportState = &EnterpriseTeamResource{}

type EnterpriseTeamResource struct {
	utils.BaseResource
}

func (r *EnterpriseTeamResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_team"
}

func (r *EnterpriseTeamResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.BaseResource.ConfigureResource(ctx, req, resp)
}

func NewEnterpriseTeamResource() resource.Resource {
	return &EnterpriseTeamResource{}
}
