// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package healthinsurance

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &HealthInsuranceResource{}
var _ resource.ResourceWithConfigure = &HealthInsuranceResource{}
var _ resource.ResourceWithImportState = &HealthInsuranceResource{}

type HealthInsuranceResource struct {
	utils.BaseResource
}

func (r *HealthInsuranceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_classic_health_insurance"
}

func (r *HealthInsuranceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewHealthInsuranceResource() resource.Resource {
	return &HealthInsuranceResource{}
}
