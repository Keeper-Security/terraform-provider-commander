// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var appUIDRegex = regexp.MustCompile(`\(UID:\s*([A-Za-z0-9_-]+)\)`)

func (r *SecretsManagerAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SecretsManagerAppResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError("Provider Configuration Error", err.Error())
		return
	}

	command := fmt.Sprintf("%s create '%s'", CmdPrefix, data.Name.ValueString())
	createResp, err := r.ApiManager.ExecuteCommand(ctx, command, "Unable to create Secrets Manager application")
	if err != nil {
		resp.Diagnostics.AddError("Create Secrets Manager App Failed", err.Error())
		return
	}

	matches := appUIDRegex.FindStringSubmatch(string(createResp.Message))
	if len(matches) < 2 || matches[1] == "" {
		resp.Diagnostics.AddError(
			"Create Secrets Manager App Failed",
			fmt.Sprintf("Unable to extract app UID from response: %s", createResp.Message),
		)
		return
	}
	data.Id = types.StringValue(matches[1])

	if !data.Shares.IsNull() && len(data.Shares.Elements()) > 0 {
		var shares []ShareEntryModel
		resp.Diagnostics.Append(data.Shares.ElementsAs(ctx, &shares, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, share := range shares {
			if err := addShareToApp(ctx, r.ApiManager, data.Id.ValueString(), share); err != nil {
				resp.Diagnostics.AddError("Create Secrets Manager App Failed", err.Error())
				return
			}
		}
	}

	if !data.AppUsers.IsNull() && len(data.AppUsers.Elements()) > 0 {
		var emails []string
		resp.Diagnostics.Append(data.AppUsers.ElementsAs(ctx, &emails, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, email := range emails {
			if err := shareAppWithUser(ctx, r.ApiManager, data.Id.ValueString(), email); err != nil {
				resp.Diagnostics.AddError("Create Secrets Manager App Failed", err.Error())
				return
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func addShareToApp(ctx context.Context, apiManager *api.ApiManager, appUID string, share ShareEntryModel) error {
	command := fmt.Sprintf("%s add --app '%s' --secret '%s'", CmdSharePrefix, appUID, share.Secret.ValueString())
	if !share.Editable.IsNull() && share.Editable.ValueBool() {
		command += " --editable"
	}
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to add secret to application")
	return err
}

func removeShareFromApp(ctx context.Context, apiManager *api.ApiManager, appUID string, secretUID string) error {
	command := fmt.Sprintf("%s remove --app '%s' --secret '%s'", CmdSharePrefix, appUID, secretUID)
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to remove secret from application")
	return err
}

func updateSharePermission(ctx context.Context, apiManager *api.ApiManager, appUID string, secretUID string, editable bool) error {
	flag := "--readonly"
	if editable {
		flag = "--editable"
	}
	command := fmt.Sprintf("%s update --app '%s' --secret '%s' %s", CmdSharePrefix, appUID, secretUID, flag)
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to update secret permissions")
	return err
}

func shareAppWithUser(ctx context.Context, apiManager *api.ApiManager, appUID string, email string) error {
	command := fmt.Sprintf("%s share '%s' --email '%s'", CmdPrefix, appUID, email)
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to share application with user")
	return err
}

func unshareAppFromUser(ctx context.Context, apiManager *api.ApiManager, appUID string, email string) error {
	command := fmt.Sprintf("%s unshare '%s' --email '%s'", CmdPrefix, appUID, email)
	_, err := apiManager.ExecuteCommand(ctx, command, "Unable to unshare application from user")
	return err
}
