// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
)

// gatewayOfflineIndicator is a substring (matched case-insensitively) Commander includes in
// the message when the gateway can't be reached to list its plugins - the response still
// reports "status":"success", but with this notice and no actual plugin entries instead of the
// usual list, e.g.:
//
//	["This Gateway currently is not online.", "Did not get router response.", "Available SaaS Plugins"]
const gatewayOfflineIndicator = "not online"

// ErrGatewayOffline is returned by FetchSaasPluginList when the configured gateway didn't
// respond, so no real plugin list is available to validate "SaaS Type" against.
var ErrGatewayOffline = errors.New("unable to load available plugins (SaaS Types). The gateway is offline")

// SaasPluginEntry describes one plugin parsed from
// `pam action saas config --gateway <gw> --list`.
type SaasPluginEntry struct {
	// Name is the value practitioners set on the "SaaS Type" custom field, e.g. "Snowflake".
	Name string
	// IntegrationType is "Builtin", "Catalog", or "Custom".
	IntegrationType string
	Description     string
}

// saasPluginLinePattern matches one plugin line from the --list output, e.g.:
//
//	"* Snowflake (Builtin) - For Snowflake, rotate the password for a user."
//	"* My Custom Plugin (Custom) - Some description."
var saasPluginLinePattern = regexp.MustCompile(`^\*\s*(.+?)\s*\((Builtin|Catalog|Custom)\)\s*-\s*(.*)$`)

// FetchSaasPluginList runs `pam action saas config --gateway <gw> --list` and parses the
// plugin name/integration type/description out of each line of its response.
//
// Commander currently returns this as a `message` array of free-text lines, with the first
// line being a header ("Available SaaS Plugins") rather than a plugin entry; this parses that
// textual format. If Commander starts returning structured JSON for this command instead, only
// this function's parsing needs to change - ValidateSaasType and its callers are unaffected.
func FetchSaasPluginList(ctx context.Context, apiManager *api.ApiManager, gateway string) ([]SaasPluginEntry, error) {
	command := fmt.Sprintf("%s --gateway %s --list", utils.CmdPamActionSaaSConfig, utils.QuoteShellSingle(gateway))
	errSummary := ErrOpListSaasPluginsForGateway(gateway)
	apiResp, err := apiManager.ExecuteCommand(ctx, command, errSummary)
	if err != nil {
		return nil, err
	}

	lines, err := messageToLines(apiResp.Message.String())
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), gatewayOfflineIndicator) {
			return nil, ErrGatewayOffline
		}
	}

	var entries []SaasPluginEntry
	for i, line := range lines {
		if i == 0 {
			// Header line (e.g. "Available SaaS Plugins"), not a plugin entry.
			continue
		}
		match := saasPluginLinePattern.FindStringSubmatch(strings.TrimSpace(line))
		if match == nil {
			continue
		}
		entries = append(entries, SaasPluginEntry{
			Name:            strings.TrimSpace(match[1]),
			IntegrationType: match[2],
			Description:     strings.TrimSpace(match[3]),
		})
	}
	return entries, nil
}

// messageToLines parses a Commander response's Message field back into its original []string
// form. Message.String() re-serializes an array-shaped message into compact JSON text (see
// api.FlexibleMessage.UnmarshalJSON); this reverses that for commands, like
// `pam action saas config --list`, that return one line per array element.
func messageToLines(message string) ([]string, error) {
	trimmed := strings.TrimSpace(message)
	if !strings.HasPrefix(trimmed, "[") {
		return nil, fmt.Errorf("expected a list-shaped response, got: %s", message)
	}
	var lines []string
	if err := json.Unmarshal([]byte(trimmed), &lines); err != nil {
		return nil, fmt.Errorf("unable to parse response as a list of strings: %w", err)
	}
	return lines, nil
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
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name)
	}
	return fmt.Errorf("SaaS Type %q is not supported for this gateway. Supported values: %s", saasType, strings.Join(names, ", "))
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
