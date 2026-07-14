// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import (
	"context"
	"errors"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var shareEntryAttrTypes = map[string]attr.Type{
	"secret":   types.StringType,
	"editable": types.BoolType,
}

func (r *SecretsManagerAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SecretsManagerAppResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError("Provider Configuration Error", err.Error())
		return
	}

	command := fmt.Sprintf("%s get '%s' --format json", CmdPrefix, state.Id.ValueString())
	getResp, err := r.ApiManager.ExecuteCommand(ctx, command, "Unable to read Secrets Manager application")
	if err != nil {
		if errors.Is(err, api.ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read Secrets Manager App Failed", err.Error())
		return
	}

	if getResp.Data == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	var appInfo GetAppResponse
	if err := utils.UnmarshalApiResponse(getResp.Data, &appInfo); err != nil {
		resp.Diagnostics.AddError("Read Secrets Manager App Failed", err.Error())
		return
	}

	if appInfo.AppUID == "" {
		resp.Diagnostics.AddError(
			"Read Secrets Manager App Failed",
			fmt.Sprintf("API returned empty app UID for application %s", state.Id.ValueString()),
		)
		return
	}

	state.Id = types.StringValue(appInfo.AppUID)
	state.Name = types.StringValue(appInfo.AppName)

	if len(appInfo.Shares) > 0 {
		shareEntries := make([]ShareEntryModel, 0, len(appInfo.Shares))
		for _, s := range appInfo.Shares {
			shareEntries = append(shareEntries, ShareEntryModel{
				Secret:   types.StringValue(s.UID),
				Editable: types.BoolValue(s.Editable),
			})
		}
		sharesSet, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: shareEntryAttrTypes}, shareEntries)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Shares = sharesSet
	} else {
		state.Shares = types.SetNull(types.ObjectType{AttrTypes: shareEntryAttrTypes})
	}

	var userEmails []string
	for _, u := range appInfo.Users {
		if u.Role != "owner" {
			userEmails = append(userEmails, u.Username)
		}
	}
	if len(userEmails) > 0 {
		usersSet, diags := types.SetValueFrom(ctx, types.StringType, userEmails)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.AppUsers = usersSet
	} else {
		state.AppUsers = types.SetNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
