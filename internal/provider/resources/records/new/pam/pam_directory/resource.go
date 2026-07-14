// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamdirectory

import (
	"context"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &PamDirectoryResource{}
var _ resource.ResourceWithConfigure = &PamDirectoryResource{}
var _ resource.ResourceWithImportState = &PamDirectoryResource{}
var _ resource.ResourceWithConfigValidators = &PamDirectoryResource{}

type PamDirectoryResource struct {
	utils.BaseResource
}

func (r *PamDirectoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_pam_directory"
}

func (r *PamDirectoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func (r *PamDirectoryResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		commonpamrecords.NewAllowSupplyHostHostnameConfigValidator(func(config PamDirectoryResourceModel) (types.Bool, *commonpamrecords.HostnameOrIPModel) {
			return commonpamrecords.AllowSupplyHostFromMachineDirectoryPamSettings(config.PamSettings), config.HostnameOrIP
		}),
	}
}

func NewPamDirectoryResource() resource.Resource {
	return &PamDirectoryResource{}
}
