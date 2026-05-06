// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	"context"
	"fmt"
	"strings"

	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_directory"
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_records"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

func recordUpdateHasMutations(plan, state commonpamdirectory.PamDirectoryResourceModel) bool {
	return !plan.Title.Equal(state.Title) ||
		!commonpamrecords.HostnameOrIPEqual(plan.HostnameOrIP, state.HostnameOrIP) ||
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
