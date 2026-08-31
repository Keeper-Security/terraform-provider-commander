// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package wifi

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordwifi "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/wifi"
)

// WifiResourceModel is the classic WiFi resource state model: shared WiFi fields
// plus the `share` attribute reconciled via classic_share.
type WifiResourceModel struct {
	commonrecordwifi.WifiModel
	classic_share.ShareModel
}
