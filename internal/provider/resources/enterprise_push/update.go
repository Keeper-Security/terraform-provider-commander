// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisepush

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Update handles email/team changes: push only to newly added emails/teams; on removal only update state.
func (r *EnterprisePushResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state EnterprisePushResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			err.Error(),
		)
		return
	}

	planEmails := sortedSetStrings(plan.Email)
	planTeams := sortedSetStrings(plan.Team)
	stateEmails := sortedSetStrings(state.Email)
	stateTeams := sortedSetStrings(state.Team)

	addedEmails := setDifference(planEmails, stateEmails)
	addedTeams := setDifference(planTeams, stateTeams)

	// Read file once (needed for push when adds, and for id/file_content_sha256).
	filePath := plan.FilePath.ValueString()
	content, fileData, err := readFileAndParseJSON(filePath)
	if err != nil {
		resp.Diagnostics.AddError("Read File Failed", err.Error())
		return
	}

	// If any new emails or teams were added, push file content only to those.
	if len(addedEmails) > 0 || len(addedTeams) > 0 {
		command := buildEnterprisePushCommandWithTargets(addedEmails, addedTeams)
		if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, plan.ManagedCompany, func() error {
			_, err := r.ApiManager.ExecuteCommand(ctx, command, "Enterprise push failed", fileData)
			return err
		}, "Enterprise Push Failed", &resp.Diagnostics); err != nil {
			return
		}
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Update state with new id and file_content_sha256 (id = f(content, plan emails/teams)).
	plan.Id = types.StringValue(computeID(content, &plan))
	plan.FileContentSha256 = types.StringValue(contentSHA256Hex(content))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
