// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package membership

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &MembershipResource{}
var _ resource.ResourceWithConfigure = &MembershipResource{}
var _ resource.ResourceWithImportState = &MembershipResource{}

type MembershipResource struct {
	utils.BaseResource
}

func (r *MembershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_membership"
}

func (r *MembershipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewMembershipResource() resource.Resource {
	return &MembershipResource{}
}
