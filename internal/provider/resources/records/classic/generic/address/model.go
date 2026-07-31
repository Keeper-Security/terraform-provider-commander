// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package address

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordaddress "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/address"
)

// AddressResourceModel is the classic Address resource state model:
// shared Address fields plus the `share` attribute reconciled via classic_share.
type AddressResourceModel struct {
	commonrecordaddress.AddressModel
	classic_share.ShareModel
}
