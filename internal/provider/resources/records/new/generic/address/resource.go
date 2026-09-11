// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package address

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &AddressResource{}
var _ resource.ResourceWithConfigure = &AddressResource{}
var _ resource.ResourceWithImportState = &AddressResource{}

type AddressResource struct {
	utils.BaseResource
}

func (r *AddressResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_address"
}

func (r *AddressResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewAddressResource() resource.Resource {
	return &AddressResource{}
}
