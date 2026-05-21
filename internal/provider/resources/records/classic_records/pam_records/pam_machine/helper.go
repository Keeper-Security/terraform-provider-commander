// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records"
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records/pam_machine"
)

func recordUpdateHasMutations(plan, state commonpammachine.PamMachineResourceModel) bool {
	return !plan.Title.Equal(state.Title) ||
		!commonpamrecords.HostnameOrIPEqual(plan.HostnameOrIP, state.HostnameOrIP) ||
		!plan.OperatingSystem.Equal(state.OperatingSystem) ||
		!plan.InstanceName.Equal(state.InstanceName) ||
		!plan.InstanceId.Equal(state.InstanceId) ||
		!plan.ProviderGroup.Equal(state.ProviderGroup) ||
		!plan.ProviderRegion.Equal(state.ProviderRegion) ||
		!plan.Notes.Equal(state.Notes)
}
