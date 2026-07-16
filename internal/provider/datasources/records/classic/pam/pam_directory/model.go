// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_directory"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PamDirectoryDataSourceModel struct {
	PamDirectory types.String `tfsdk:"pam_directory"`
	commonpamdirectory.PamDirectoryResourceModel
	classic_share.ShareModel
}
