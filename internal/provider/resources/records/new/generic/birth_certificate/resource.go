// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package birthcertificate

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &BirthCertificateResource{}
var _ resource.ResourceWithConfigure = &BirthCertificateResource{}
var _ resource.ResourceWithImportState = &BirthCertificateResource{}

type BirthCertificateResource struct {
	utils.BaseResource
}

func (r *BirthCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_birth_certificate"
}

func (r *BirthCertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewBirthCertificateResource() resource.Resource {
	return &BirthCertificateResource{}
}
