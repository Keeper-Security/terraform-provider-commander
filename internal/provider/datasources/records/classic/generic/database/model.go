// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecorddatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/database"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DatabaseDataSourceModel maps a Keeper `databaseCredentials` vault record for read-only access.
type DatabaseDataSourceModel struct {
	Database types.String `tfsdk:"database"`
	commonrecorddatabase.DatabaseModel
	classic_share.ShareModel
}
