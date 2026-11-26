// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamcloud

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonpamcloud "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_cloud"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PamCloudDataSourceModel struct {
	PamCloud types.String `tfsdk:"pam_cloud"`
	commonpamcloud.PamCloudResourceModel
	classic_share.ShareModel
}
