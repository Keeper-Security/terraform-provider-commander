// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpammachine

import (
	"context"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &PamMachineResource{}
var _ resource.ResourceWithConfigure = &PamMachineResource{}
var _ resource.ResourceWithImportState = &PamMachineResource{}
var _ resource.ResourceWithConfigValidators = &PamMachineResource{}

type PamMachineResource struct {
	utils.BaseResource
}

func (r *PamMachineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_new_pam_machine"
}

func (r *PamMachineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.ConfigureResource(ctx, req, resp)
}

func (r *PamMachineResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		commonpamrecords.NewAllowSupplyHostHostnameConfigValidator(func(config PamMachineResourceModel) (types.Bool, *commonpamrecords.HostnameOrIPModel) {
			return commonpamrecords.AllowSupplyHostFromMachineDirectoryPamSettings(config.PamSettings), config.HostnameOrIP
		}),
	}
}

func NewPamMachineResource() resource.Resource {
	return &PamMachineResource{}
}
