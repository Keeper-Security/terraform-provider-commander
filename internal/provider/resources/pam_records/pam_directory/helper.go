// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_directory"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_records"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func buildHostnameOrIPJSON(h *commonpamrecords.HostnameOrIPModel) string {
	m := map[string]string{
		"hostName": h.HostName.ValueString(),
	}
	if !h.AdministrativePort.IsNull() && !h.AdministrativePort.IsUnknown() {
		m["port"] = strconv.FormatInt(int64(h.AdministrativePort.ValueInt32()), 10)
	} else {
		m["port"] = ""
	}

	b, _ := json.Marshal(m)
	return string(b)
}

func appendHostnameOrIPField(parts *[]string, h *commonpamrecords.HostnameOrIPModel) {
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

func appendOptionalCheckboxField(parts *[]string, flag string, v types.Bool) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	*parts = append(*parts, fmt.Sprintf("'%s=%t'", flag, v.ValueBool()))
}

func appendChangedCheckboxField(parts *[]string, flag string, plan, state types.Bool) {
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
	*parts = append(*parts, fmt.Sprintf("'%s=%t'", flag, plan.ValueBool()))
}

func appendOptionalDirectoryTypeField(parts *[]string, v types.String) {
	if v.IsNull() || v.IsUnknown() {
		return
	}
	*parts = append(*parts, fmt.Sprintf("'%s=%s'", FlagDirectoryType, v.ValueString()))
}

// appendAlternativeIPsField joins the set elements with newlines and writes
// them as a single multiline field value.
func appendAlternativeIPsField(parts *[]string, s types.Set) {
	if s.IsNull() || s.IsUnknown() || len(s.Elements()) == 0 {
		return
	}
	var ips []string
	diags := s.ElementsAs(context.Background(), &ips, false)
	if diags.HasError() {
		return
	}
	joined := strings.Join(ips, "\\n")
	*parts = append(*parts, fmt.Sprintf("'%s=%s'", FlagAlternativeIPs, joined))
}

func hostnameOrIPEqual(a, b *commonpamrecords.HostnameOrIPModel) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.HostName.Equal(b.HostName) && a.AdministrativePort.Equal(b.AdministrativePort)
}

func recordUpdateHasMutations(plan, state commonpamdirectory.PamDirectoryResourceModel) bool {
	return !plan.Title.Equal(state.Title) ||
		!hostnameOrIPEqual(plan.HostnameOrIP, state.HostnameOrIP) ||
		!plan.UseSSL.Equal(state.UseSSL) ||
		!plan.DomainName.Equal(state.DomainName) ||
		!plan.AlternativeIPs.Equal(state.AlternativeIPs) ||
		!plan.DirectoryId.Equal(state.DirectoryId) ||
		!plan.DirectoryType.Equal(state.DirectoryType) ||
		!plan.UserMatch.Equal(state.UserMatch) ||
		!plan.ProviderGroup.Equal(state.ProviderGroup) ||
		!plan.ProviderRegion.Equal(state.ProviderRegion) ||
		!plan.Notes.Equal(state.Notes)
}
