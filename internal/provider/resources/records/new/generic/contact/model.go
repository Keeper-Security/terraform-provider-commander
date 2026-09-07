// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordcontact "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/contact"
)

// ContactResourceModel is the New (NSF) contact resource state model: shared contact
// fields plus the `share` attribute reconciled via new_share.
type ContactResourceModel struct {
	commonrecordcontact.ContactModel
	new_share.ShareModel
}
