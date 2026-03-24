// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &EpmPolicyResource{}
var _ resource.ResourceWithConfigure = &EpmPolicyResource{}
var _ resource.ResourceWithImportState = &EpmPolicyResource{}
var _ resource.ResourceWithConfigValidators = &EpmPolicyResource{}

type EpmPolicyResource struct {
	utils.BaseResource
}

func (r *EpmPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_epm_policy"
}

func (r *EpmPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func (r *EpmPolicyResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		epmPolicyConfigValidator{},
	}
}

func NewEpmPolicyResource() resource.Resource {
	return &EpmPolicyResource{}
}
