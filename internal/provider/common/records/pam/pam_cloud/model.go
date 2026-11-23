// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamcloud

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
)

// PamCloudResourceModel is the Terraform state shared by the classic and new
// PAM cloud resources.
type PamCloudResourceModel struct {
	commonrecordsutils.BaseVaultRecordModel

	URL               types.String `tfsdk:"url"`
	AccountIdentifier types.String `tfsdk:"account_identifier"`
}
