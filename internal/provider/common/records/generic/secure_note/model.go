// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package securenote

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SecureNoteModel maps a Keeper `encryptedNotes` vault record.
// Shared between the resource and data source.
type SecureNoteModel struct {
	utils.BaseVaultRecordModel

	SecuredNote types.String `tfsdk:"secured_note"`
	Date        types.String `tfsdk:"date"`

	Custom []utils.CustomFieldModel `tfsdk:"custom"`
}
