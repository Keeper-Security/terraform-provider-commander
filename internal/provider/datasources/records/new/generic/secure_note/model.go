// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package securenote

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordsecurenote "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/secure_note"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SecureNoteDataSourceModel maps a Keeper `encryptedNotes` vault record for read-only access.
type SecureNoteDataSourceModel struct {
	SecureNote types.String `tfsdk:"secure_note"`
	commonrecordsecurenote.SecureNoteModel
	new_share.ShareModel
}
