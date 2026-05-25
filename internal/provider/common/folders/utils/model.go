// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import "github.com/hashicorp/terraform-plugin-framework/types"

// IdentityModel holds id and name attributes shared by folder resources and data sources.
type IdentityModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
