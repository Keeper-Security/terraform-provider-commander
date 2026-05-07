// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package managedcompany

import (
	"context"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *ManagedCompanyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data ManagedCompanyResourceModel

	// Get planned data from Terraform
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate ApiManager is configured
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR,
			err.Error(),
		)
		return
	}

	if err := utils.RunWithMspContext(ctx, r.ApiManager, func() error {
		return addManagedCompany(ctx, r.ApiManager, &data)
	}, "Create Managed Company Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the ID in the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// buildManagedCompanyAddCommand builds the Commander CLI command for adding a managed company.
func addManagedCompany(ctx context.Context, apiManager *api.ApiManager, data *ManagedCompanyResourceModel) error {
	var parts []string

	parts = append(parts, "msp-add")

	// Required parameters
	parts = append(parts, fmt.Sprintf("'%s'", data.Name.ValueString()))
	parts = append(parts, fmt.Sprintf("--node '%s'", data.Node.ValueString()))
	parts = append(parts, fmt.Sprintf("--plan '%s'", data.Plan.ValueString()))

	// Optional parameters
	if !data.Seats.IsNull() {
		parts = append(parts, fmt.Sprintf("--seats %d", data.Seats.ValueInt64()))
	}

	if !data.FilePlan.IsNull() {
		parts = append(parts, fmt.Sprintf("--file-plan '%s'", data.FilePlan.ValueString()))
	}

	// Add-ons
	if !data.AddOns.IsNull() && !data.AddOns.IsUnknown() {
		for _, addOn := range normalizeAddOns(data.AddOns) {
			parts = append(parts, fmt.Sprintf("--addon '%s'", addOn))
		}
	}

	command := strings.Join(parts, " ")

	apiResp, err := apiManager.ExecuteCommand(ctx, command, "Unable to create managed company")
	if err != nil {
		return err
	}

	data.Id = types.StringValue(apiResp.Message.String())

	return nil

}
