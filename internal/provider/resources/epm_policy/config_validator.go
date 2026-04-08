// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import (
	"context"

	commonepm "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/epm_policy"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// epmPolicyConfigValidator validates policy-type allowed fields (e.g. least_privilege allows only certain attributes).
type epmPolicyConfigValidator struct{}

func (epmPolicyConfigValidator) Description(ctx context.Context) string {
	return "Validates fields allowed and required for the policy type and status (e.g. " + commonepm.PolicyTypeLeastPrivilege + " allows only policy name, type, status " + commonepm.StatusDescriptionForLeastPrivilege() + ", machine_collections; elevation/file_access enforce require control and three collections; command enforce requires control and user/machine collections only—applications are never allowed for command)."
}

func (epmPolicyConfigValidator) MarkdownDescription(ctx context.Context) string {
	return "Validates that only fields allowed for the given policy type are set."
}

func (v epmPolicyConfigValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model EpmPolicyResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyType := model.PolicyType.ValueString()
	status := model.Status.ValueString()

	commonepm.ValidatePolicyTypeAllowedFields(
		policyType, status,
		model.Control, model.DayFilter,
		model.UserGroups, model.MachineCollections, model.Applications, model.TimeFilter, model.DateFilter,
		path.Root("status"), path.Root("control"), path.Root("day_filter"), path.Root("user_groups"), path.Root("machine_collections"), path.Root("applications"),
		path.Root("time_filter"), path.Root("date_filter"),
		&resp.Diagnostics,
	)
}
