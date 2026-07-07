// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package login

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordlogin "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/login"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LoginDataSourceModel adds a lookup key (`login_record`) to the shared login model.
// The lookup attribute is named login_record because login is the username field on the record.
type LoginDataSourceModel struct {
	LoginRecord types.String `tfsdk:"login_record"`
	commonrecordlogin.LoginModel
	classic_share.ShareModel
}
