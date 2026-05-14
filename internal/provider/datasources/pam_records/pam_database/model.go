// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdatabase

import (
	commonpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_database"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PamDatabaseDataSourceModel struct {
	PamDatabase types.String `tfsdk:"pam_database"`
	commonpamdatabase.PamDatabaseResourceModel
}
