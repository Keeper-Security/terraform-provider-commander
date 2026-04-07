// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import "github.com/hashicorp/terraform-plugin-framework/types"

type SecretsManagerAppResourceModel struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Shares   types.Set    `tfsdk:"shares"`
	AppUsers types.Set    `tfsdk:"app_users"`
}

type ShareEntryModel struct {
	Secret   types.String `tfsdk:"secret"`
	Editable types.Bool   `tfsdk:"editable"`
}
