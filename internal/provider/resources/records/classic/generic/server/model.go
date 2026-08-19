// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordserver "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/server"
)

// ServerResourceModel is the classic serverCredentials resource state model.
type ServerResourceModel struct {
	commonrecordserver.ServerModel
	classic_share.ShareModel
}
