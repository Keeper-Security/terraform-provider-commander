// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package securenote

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordsecurenote "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/secure_note"
)

// SecureNoteResourceModel is the classic encryptedNotes resource state model.
type SecureNoteResourceModel struct {
	commonrecordsecurenote.SecureNoteModel
	classic_share.ShareModel
}
