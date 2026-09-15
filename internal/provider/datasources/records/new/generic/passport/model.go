// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package passport

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordpassport "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/passport"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PassportDataSourceModel adds a lookup key (`passport`) to the shared passport model.
type PassportDataSourceModel struct {
	Passport types.String `tfsdk:"passport"`
	commonrecordpassport.PassportModel
	new_share.ShareModel
}
