// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package ssncard

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &SsnCardResource{}
var _ resource.ResourceWithConfigure = &SsnCardResource{}
var _ resource.ResourceWithImportState = &SsnCardResource{}

type SsnCardResource struct {
	utils.BaseResource
}

func (r *SsnCardResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_classic_ssn_card"
}

func (r *SsnCardResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewSsnCardResource() resource.Resource {
	return &SsnCardResource{}
}
