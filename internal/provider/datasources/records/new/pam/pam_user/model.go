// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package newpamuser

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_user"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PamUserDataSourceModel struct {
	PamUser types.String `tfsdk:"pam_user"`
	commonpamuser.PamUserSharedModel
	new_share.ShareModel
}
