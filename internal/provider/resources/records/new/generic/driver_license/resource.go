// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package driverlicense

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &DriverLicenseResource{}
var _ resource.ResourceWithConfigure = &DriverLicenseResource{}
var _ resource.ResourceWithImportState = &DriverLicenseResource{}

type DriverLicenseResource struct {
	utils.BaseResource
}

func (r *DriverLicenseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_driver_license"
}

func (r *DriverLicenseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewDriverLicenseResource() resource.Resource {
	return &DriverLicenseResource{}
}
