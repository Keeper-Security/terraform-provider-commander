// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &SaasConfigurationResource{}
var _ resource.ResourceWithConfigure = &SaasConfigurationResource{}
var _ resource.ResourceWithImportState = &SaasConfigurationResource{}

type SaasConfigurationResource struct {
	utils.BaseResource
}

func (r *SaasConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_saas_configuration"
}

func (r *SaasConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewSaasConfigurationResource() resource.Resource {
	return &SaasConfigurationResource{}
}
