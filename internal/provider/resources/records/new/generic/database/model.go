// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package database

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecorddatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/database"
)

// DatabaseResourceModel is the classic databaseCredentials resource state model.
type DatabaseResourceModel struct {
	commonrecorddatabase.DatabaseModel
	new_share.ShareModel
}
