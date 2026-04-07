// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var shareEntryAttrTypes = map[string]attr.Type{
	"secret":   types.StringType,
	"editable": types.BoolType,
}

type ShareEntryModel struct {
	Secret   types.String `tfsdk:"secret"`
	Editable types.Bool   `tfsdk:"editable"`
}

type ShareResponse struct {
	UID       string `json:"uid"`
	Editable  bool   `json:"editable"`
	ShareType string `json:"share_type"`
	Title     string `json:"title"`
	Type      string `json:"type"`
}

type UserResponse struct {
	Username   string `json:"username"`
	Role       string `json:"role"`
	Editable   bool   `json:"editable"`
	ShareAdmin bool   `json:"share_admin"`
	Shareable  bool   `json:"shareable"`
}

type GetAppResponse struct {
	AppName       string          `json:"app_name"`
	AppUID        string          `json:"app_uid"`
	ClientDevices []interface{}   `json:"client_devices"`
	Shares        []ShareResponse `json:"shares"`
	Users         []UserResponse  `json:"users"`
}

const cmdPrefix = "secrets-manager app"

func (d *SecretsManagerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SecretsManagerDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError("Provider Configuration Error", err.Error())
		return
	}

	command := fmt.Sprintf("%s get '%s' --format json", cmdPrefix, data.Application.ValueString())
	getResp, err := d.ApiManager.ExecuteCommand(ctx, command, "Unable to read Secrets Manager application")
	if err != nil {
		resp.Diagnostics.AddError("Read Secrets Manager Data Source Failed", err.Error())
		return
	}

	if getResp.Data == nil {
		resp.Diagnostics.AddError(
			"Read Secrets Manager Data Source Failed",
			fmt.Sprintf("Application '%s' not found", data.Application.ValueString()),
		)
		return
	}

	var appInfo GetAppResponse
	if err := utils.UnmarshalApiResponse(getResp.Data, &appInfo); err != nil {
		resp.Diagnostics.AddError("Read Secrets Manager Data Source Failed", err.Error())
		return
	}

	if appInfo.AppUID == "" {
		resp.Diagnostics.AddError(
			"Read Secrets Manager Data Source Failed",
			fmt.Sprintf("Application '%s' not found", data.Application.ValueString()),
		)
		return
	}

	data.Id = types.StringValue(appInfo.AppUID)
	data.Name = types.StringValue(appInfo.AppName)

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
		data.Shares = sharesSet
	} else {
		data.Shares = types.SetValueMust(types.ObjectType{AttrTypes: shareEntryAttrTypes}, []attr.Value{})
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
		data.AppUsers = usersSet
	} else {
		data.AppUsers = types.SetValueMust(types.StringType, []attr.Value{})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
