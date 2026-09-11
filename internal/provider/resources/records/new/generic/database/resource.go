// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &DatabaseResource{}
var _ resource.ResourceWithConfigure = &DatabaseResource{}
var _ resource.ResourceWithImportState = &DatabaseResource{}

type DatabaseResource struct {
	utils.BaseResource
}

func (r *DatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_database"
}

func (r *DatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func NewDatabaseResource() resource.Resource {
	return &DatabaseResource{}
}
