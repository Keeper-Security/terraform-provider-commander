// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &PamConfigurationResource{}
var _ resource.ResourceWithConfigure = &PamConfigurationResource{}
var _ resource.ResourceWithImportState = &PamConfigurationResource{}
var _ resource.ResourceWithConfigValidators = &PamConfigurationResource{}

type PamConfigurationResource struct {
	utils.BaseResource
}

func (r *PamConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pam_configuration"
}

func (r *PamConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewPamConfigurationResource() resource.Resource {
	return &PamConfigurationResource{}
}

func (r *PamConfigurationResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		environmentBlocksValidator{},
	}
}
