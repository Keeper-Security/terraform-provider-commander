// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *SecretsManagerAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state SecretsManagerAppResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError("Provider Configuration Error", err.Error())
		return
	}

	appUID := state.Id.ValueString()

	if plan.Name.ValueString() != state.Name.ValueString() {
		command := fmt.Sprintf("%s update '%s' --name '%s'", CmdPrefix, appUID, plan.Name.ValueString())
		if _, err := r.ApiManager.ExecuteCommand(ctx, command, "Unable to rename Secrets Manager application"); err != nil {
			resp.Diagnostics.AddError("Update Secrets Manager App Failed", err.Error())
			return
		}
	}

	planShares := make(map[string]ShareEntryModel)
	if !plan.Shares.IsNull() && len(plan.Shares.Elements()) > 0 {
		var entries []ShareEntryModel
		resp.Diagnostics.Append(plan.Shares.ElementsAs(ctx, &entries, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, e := range entries {
			planShares[e.Secret.ValueString()] = e
		}
	}

	stateShares := make(map[string]ShareEntryModel)
	if !state.Shares.IsNull() && len(state.Shares.Elements()) > 0 {
		var entries []ShareEntryModel
		resp.Diagnostics.Append(state.Shares.ElementsAs(ctx, &entries, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, e := range entries {
			stateShares[e.Secret.ValueString()] = e
		}
	}

	for uid := range stateShares {
		if _, exists := planShares[uid]; !exists {
			if err := removeShareFromApp(ctx, r.ApiManager, appUID, uid); err != nil {
				resp.Diagnostics.AddError("Update Secrets Manager App Failed", err.Error())
				return
			}
		}
	}

	for uid, planShare := range planShares {
		stateShare, exists := stateShares[uid]
		if !exists {
			if err := addShareToApp(ctx, r.ApiManager, appUID, planShare); err != nil {
				resp.Diagnostics.AddError("Update Secrets Manager App Failed", err.Error())
				return
			}
		} else if planShare.Editable.ValueBool() != stateShare.Editable.ValueBool() {
			if err := updateSharePermission(ctx, r.ApiManager, appUID, uid, planShare.Editable.ValueBool()); err != nil {
				resp.Diagnostics.AddError("Update Secrets Manager App Failed", err.Error())
				return
			}
		}
	}

	planUsers := make(map[string]bool)
	if !plan.AppUsers.IsNull() && len(plan.AppUsers.Elements()) > 0 {
		var emails []string
		resp.Diagnostics.Append(plan.AppUsers.ElementsAs(ctx, &emails, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, e := range emails {
			planUsers[e] = true
		}
	}

	stateUsers := make(map[string]bool)
	if !state.AppUsers.IsNull() && len(state.AppUsers.Elements()) > 0 {
		var emails []string
		resp.Diagnostics.Append(state.AppUsers.ElementsAs(ctx, &emails, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, e := range emails {
			stateUsers[e] = true
		}
	}

	for email := range stateUsers {
		if !planUsers[email] {
			if err := unshareAppFromUser(ctx, r.ApiManager, appUID, email); err != nil {
				resp.Diagnostics.AddError("Update Secrets Manager App Failed", err.Error())
				return
			}
		}
	}

	for email := range planUsers {
		if !stateUsers[email] {
			if err := shareAppWithUser(ctx, r.ApiManager, appUID, email); err != nil {
				resp.Diagnostics.AddError("Update Secrets Manager App Failed", err.Error())
				return
			}
		}
	}

	plan.Id = state.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
