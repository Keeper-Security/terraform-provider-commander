// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package healthinsurance

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/classic_share"
	commonrecordhealthinsurance "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/health_insurance"
)

// HealthInsuranceResourceModel is the classic health insurance resource state model:
// shared health insurance fields plus the `share` attribute reconciled via classic_share.
type HealthInsuranceResourceModel struct {
	commonrecordhealthinsurance.HealthInsuranceModel
	classic_share.ShareModel
}
