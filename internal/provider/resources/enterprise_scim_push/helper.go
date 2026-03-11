// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// buildScimPushCommand builds the scim push command from the model.
// Format: scim push [scim_id] --source {google,ad,record} --record RECORD --auto-approve {on,off}.
func buildScimPushCommand(data *EnterpriseScimPushResourceModel) string {
	source := strings.TrimSpace(strings.ToLower(data.Source.ValueString()))
	autoApprove := AutoApproveOff
	if data.AutoApprove.ValueBool() {
		autoApprove = AutoApproveOn
	}
	parts := []string{
		CmdScimPush,
		data.ScimId.ValueString(),
		fmt.Sprintf("%s '%s'", FlagSource, source),
		fmt.Sprintf("%s '%s'", FlagRecord, data.Record.ValueString()),
		fmt.Sprintf("%s '%s'", FlagAutoApprove, autoApprove),
	}
	return strings.Join(parts, " ")
}

// computeID returns a deterministic ID for the one-time push (scim_id + source + record + auto_approve + managed_company).
func computeID(data *EnterpriseScimPushResourceModel) string {
	managedCompany := ""
	if !data.ManagedCompany.IsNull() {
		managedCompany = data.ManagedCompany.ValueString()
	}
	h := sha256.New()
	h.Write([]byte(data.ScimId.ValueString()))
	h.Write([]byte("\n"))
	h.Write([]byte(strings.TrimSpace(strings.ToLower(data.Source.ValueString()))))
	h.Write([]byte("\n"))
	h.Write([]byte(data.Record.ValueString()))
	h.Write([]byte("\n"))
	if data.AutoApprove.ValueBool() {
		h.Write([]byte("true"))
	} else {
		h.Write([]byte("false"))
	}
	h.Write([]byte("\n"))
	h.Write([]byte(managedCompany))
	return hex.EncodeToString(h.Sum(nil))
}
