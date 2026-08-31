// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package ssncard

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordssncard "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/ssn_card"
)

// SsnCardResourceModel is the classic SSN Card (Identity Card) resource state model:
// shared SSN Card fields plus the `share` attribute reconciled via classic_share.
type SsnCardResourceModel struct {
	commonrecordssncard.SsnCardModel
	classic_share.ShareModel
}
