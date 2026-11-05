// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
)

// gatewayOfflineIndicator is a substring (matched case-insensitively) Commander includes in
// the message when the configured gateway can't be reached - the response still reports
// "status":"success", but with this notice and no real plugin entries, e.g.:
//
//	["This Gateway currently is not online.", "Did not get router response.", "Available SaaS Plugins"]
const gatewayOfflineIndicator = "not online"

// ErrGatewayOffline is returned by FetchSaasPluginList when the configured gateway didn't respond.
var ErrGatewayOffline = errors.New("unable to load available plugins (SaaS Types). The gateway is offline")

// SaasPluginEntry describes one plugin parsed from
// `pam action saas config --gateway <gw> --list --format json`.
type SaasPluginEntry struct {
	// Name is the value practitioners set on the "SaaS Type" custom field, e.g. "Snowflake".
	Name            string
	IntegrationType string
	Description     string
}

// saasPluginListEntry is one element of the `data` array, e.g.:
//
//	{"name": "Snowflake", "type": "builtin", "summary": "For Snowflake, rotate the password for a user.", ...}
type saasPluginListEntry struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
}

// FetchSaasPluginList runs `pam action saas config --gateway <gw> --list --format json`.
func FetchSaasPluginList(ctx context.Context, apiManager *api.ApiManager, gateway string) ([]SaasPluginEntry, error) {
	command := fmt.Sprintf("%s --gateway %s --list %s", utils.CmdPamActionSaaSConfig, utils.QuoteShellSingle(gateway), utils.FlagFormatJSON)
	errSummary := ErrOpListSaasPluginsForGateway(gateway)
	apiResp, err := apiManager.ExecuteCommand(ctx, command, errSummary)
	if err != nil {
		return nil, err
	}

	if strings.Contains(strings.ToLower(apiResp.Message.String()), gatewayOfflineIndicator) {
		return nil, ErrGatewayOffline
	}

	if apiResp.Data == nil {
		return nil, nil
	}

	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to process plugin list response: %w", errSummary, err)
	}

	var rawEntries []saasPluginListEntry
	if err := json.Unmarshal(dataBytes, &rawEntries); err != nil {
		return nil, fmt.Errorf("%s: unable to parse plugin list response: %w", errSummary, err)
	}

	entries := make([]SaasPluginEntry, 0, len(rawEntries))
	for _, e := range rawEntries {
		entries = append(entries, SaasPluginEntry{
			Name:            strings.TrimSpace(e.Name),
			IntegrationType: strings.TrimSpace(e.Type),
			Description:     strings.TrimSpace(e.Summary),
		})
	}
	return entries, nil
}

// ValidateSaasType reports an error if saasType doesn't exactly match any plugin name in
// entries. An empty saasType is not this function's concern - the existing
// requiredSaasTypeCustomFieldValidator already enforces that the field is present.
func ValidateSaasType(saasType string, entries []SaasPluginEntry) error {
	saasType = strings.TrimSpace(saasType)
	if saasType == "" {
		return nil
	}
	for _, e := range entries {
		if e.Name == saasType {
			return nil
		}
	}
	return fmt.Errorf("SaaS Type %q is not supported for this gateway. Supported values:\n\n%s", saasType, formatSaasPluginEntriesTable(entries))
}

// formatSaasPluginEntriesTable renders entries as an aligned table for a diagnostic error detail.
func formatSaasPluginEntriesTable(entries []SaasPluginEntry) string {
	var buf strings.Builder
	tw := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SAAS TYPE\tPLUGIN TYPE\tDESCRIPTION")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", e.Name, e.IntegrationType, e.Description)
	}
	_ = tw.Flush()
	return buf.String()
}

// SaasTypeFromCustom returns the value of the custom field labeled "SaaS Type", or "" if absent.
func SaasTypeFromCustom(custom []commonrecordsutils.CustomFieldModel) string {
	for _, c := range custom {
		if strings.TrimSpace(c.Label.ValueString()) == SaaSTypeCustomFieldLabel {
			return c.Value.ValueString()
		}
	}
	return ""
}
