// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package login

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordlogin "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/login"
)

// LoginResourceModel is an alias for the shared login model.
type LoginResourceModel struct {
	commonrecordlogin.LoginModel
	new_share.ShareModel
}
