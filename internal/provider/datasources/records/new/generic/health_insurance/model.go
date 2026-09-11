// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package healthinsurance

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	commonrecordhealthinsurance "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/generic/health_insurance"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// HealthInsuranceDataSourceModel adds a lookup key (`health_insurance`) to the
// shared health insurance model.
type HealthInsuranceDataSourceModel struct {
	HealthInsurance types.String `tfsdk:"health_insurance"`
	commonrecordhealthinsurance.HealthInsuranceModel
	new_share.ShareModel
}
