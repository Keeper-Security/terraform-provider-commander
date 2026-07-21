// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package passport

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordpassport "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/passport"
)

// PassportResourceModel is the classic passport resource state model: shared
// passport fields plus the `share` attribute reconciled via classic_share.
type PassportResourceModel struct {
	commonrecordpassport.PassportModel
	classic_share.ShareModel
}
