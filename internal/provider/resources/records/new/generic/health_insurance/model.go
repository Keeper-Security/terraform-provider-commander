// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package healthinsurance

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordhealthinsurance "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/health_insurance"
)

// HealthInsuranceResourceModel is the New (NSF) health insurance resource state model.
type HealthInsuranceResourceModel struct {
	commonrecordhealthinsurance.HealthInsuranceModel
	new_share.ShareModel
}
