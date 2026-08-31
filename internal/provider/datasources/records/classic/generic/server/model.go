// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordserver "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/server"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ServerDataSourceModel maps a Keeper `serverCredentials` vault record for read-only access.
type ServerDataSourceModel struct {
	Server types.String `tfsdk:"server"`
	commonrecordserver.ServerModel
	classic_share.ShareModel
}
