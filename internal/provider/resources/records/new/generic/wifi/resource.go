// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &WifiResource{}
var _ resource.ResourceWithConfigure = &WifiResource{}
var _ resource.ResourceWithImportState = &WifiResource{}

type WifiResource struct {
	utils.BaseResource
}

func (r *WifiResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_wifi"
}

func (r *WifiResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewWifiResource() resource.Resource {
	return &WifiResource{}
}
