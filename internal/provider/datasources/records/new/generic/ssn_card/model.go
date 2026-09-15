// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package ssncard

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordssncard "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/ssn_card"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SsnCardDataSourceModel adds a lookup key (`ssn_card`) to the shared SSN Card model.
type SsnCardDataSourceModel struct {
	SsnCard types.String `tfsdk:"ssn_card"`
	commonrecordssncard.SsnCardModel
	new_share.ShareModel
}
