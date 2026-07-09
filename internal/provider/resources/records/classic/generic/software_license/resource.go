// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package softwarelicense

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &SoftwareLicenseResource{}
var _ resource.ResourceWithConfigure = &SoftwareLicenseResource{}
var _ resource.ResourceWithImportState = &SoftwareLicenseResource{}

type SoftwareLicenseResource struct {
	utils.BaseResource
}

func (r *SoftwareLicenseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_classic_software_license"
}

func (r *SoftwareLicenseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewSoftwareLicenseResource() resource.Resource {
	return &SoftwareLicenseResource{}
}
