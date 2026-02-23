// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisepush

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnterprisePushResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	FilePath           types.String `tfsdk:"file_path"`
	FileContentSha256  types.String `tfsdk:"file_content_sha256"`
	Email              types.Set    `tfsdk:"email"`
	Team               types.Set    `tfsdk:"team"`
	ManagedCompany     types.String `tfsdk:"managed_company"`
}
