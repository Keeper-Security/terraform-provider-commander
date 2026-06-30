// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpammachine

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &PamMachineResource{}
var _ resource.ResourceWithConfigure = &PamMachineResource{}
var _ resource.ResourceWithImportState = &PamMachineResource{}

type PamMachineResource struct {
	utils.BaseResource
}

func (r *PamMachineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_pam_machine"
}

func (r *PamMachineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewPamMachineResource() resource.Resource {
	return &PamMachineResource{}
}
