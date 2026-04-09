// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &SecretsManagerAppResource{}
var _ resource.ResourceWithConfigure = &SecretsManagerAppResource{}
var _ resource.ResourceWithImportState = &SecretsManagerAppResource{}

type SecretsManagerAppResource struct {
	utils.BaseResource
}

func (r *SecretsManagerAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secrets_manager"
}

func (r *SecretsManagerAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewSecretsManagerAppResource() resource.Resource {
	return &SecretsManagerAppResource{}
}
