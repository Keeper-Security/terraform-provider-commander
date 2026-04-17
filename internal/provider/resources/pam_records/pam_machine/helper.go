// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"encoding/json"
	"fmt"

	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_machine"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func buildHostnameOrIPJSON(h *commonpammachine.HostnameOrIPModel) string {
	m := map[string]string{
		"hostName": h.HostName.ValueString(),
	}
	if !h.Port.IsNull() && !h.Port.IsUnknown() {
		m["port"] = h.Port.ValueString()
	} else {
		m["port"] = ""
	}

	b, _ := json.Marshal(m)
	return string(b)
}

func appendHostnameOrIPField(parts *[]string, h *commonpammachine.HostnameOrIPModel) {
	if h == nil {
		return
	}
	hostnameJSON := buildHostnameOrIPJSON(h)
	*parts = append(*parts, fmt.Sprintf("'%s=$JSON:%s'", FlagPamHostname, hostnameJSON))
}

func appendOptionalTextField(parts *[]string, flag string, v types.String) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	*parts = append(*parts, fmt.Sprintf("'%s=%s'", flag, v.ValueString()))
}

func appendChangedTextField(parts *[]string, flag string, plan, state types.String) {
	if plan.Equal(state) {
		return
	}
	if plan.IsUnknown() {
		return
	}
	if plan.IsNull() {
		*parts = append(*parts, fmt.Sprintf("'%s='", flag))
		return
	}
	*parts = append(*parts, fmt.Sprintf("'%s=%s'", flag, plan.ValueString()))
}

func hostnameOrIPEqual(a, b *commonpammachine.HostnameOrIPModel) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.HostName.Equal(b.HostName) && a.Port.Equal(b.Port)
}

func recordUpdateHasMutations(plan, state commonpammachine.PamMachineResourceModel) bool {
	return !plan.Title.Equal(state.Title) ||
		!hostnameOrIPEqual(plan.HostnameOrIP, state.HostnameOrIP) ||
		!plan.OperatingSystem.Equal(state.OperatingSystem) ||
		!plan.InstanceName.Equal(state.InstanceName) ||
		!plan.InstanceId.Equal(state.InstanceId) ||
		!plan.ProviderGroup.Equal(state.ProviderGroup) ||
		!plan.ProviderRegion.Equal(state.ProviderRegion) ||
		!plan.Notes.Equal(state.Notes)
}
