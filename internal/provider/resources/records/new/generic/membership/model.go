// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package membership

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordmembership "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/membership"
)

// MembershipResourceModel is the New (NSF) membership resource state model.
type MembershipResourceModel struct {
	commonrecordmembership.MembershipModel
	new_share.ShareModel
}
