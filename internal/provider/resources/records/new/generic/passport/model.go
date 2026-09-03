// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package passport

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordpassport "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/passport"
)

// PassportResourceModel is the new passport resource state model: shared
// passport fields plus the `share` attribute reconciled via new_share.
type PassportResourceModel struct {
	commonrecordpassport.PassportModel
	new_share.ShareModel
}
