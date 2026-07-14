// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamdatabase

import (
	"context"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &PamDatabaseResource{}
var _ resource.ResourceWithConfigure = &PamDatabaseResource{}
var _ resource.ResourceWithImportState = &PamDatabaseResource{}
var _ resource.ResourceWithConfigValidators = &PamDatabaseResource{}

type PamDatabaseResource struct {
	utils.BaseResource
}

func (r *PamDatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_pam_database"
}

func (r *PamDatabaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func (r *PamDatabaseResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		commonpamrecords.NewAllowSupplyHostHostnameConfigValidator(func(config PamDatabaseResourceModel) (types.Bool, *commonpamrecords.HostnameOrIPModel) {
			return commonpamrecords.AllowSupplyHostFromDatabasePamSettings(config.PamSettings), config.HostnameOrIP
		}),
	}
}

func NewPamDatabaseResource() resource.Resource {
	return &PamDatabaseResource{}
}
