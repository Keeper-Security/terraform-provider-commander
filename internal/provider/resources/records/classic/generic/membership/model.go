// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package membership

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordmembership "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/membership"
)

// MembershipResourceModel is the classic membership resource state model: shared
// membership fields plus the `share` attribute reconciled via classic_share.
type MembershipResourceModel struct {
	commonrecordmembership.MembershipModel
	classic_share.ShareModel
}
