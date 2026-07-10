// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package securenote

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &SecureNoteResource{}
var _ resource.ResourceWithConfigure = &SecureNoteResource{}
var _ resource.ResourceWithImportState = &SecureNoteResource{}

type SecureNoteResource struct {
	utils.BaseResource
}

func (r *SecureNoteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_classic_secure_note"
}

func (r *SecureNoteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewSecureNoteResource() resource.Resource {
	return &SecureNoteResource{}
}
