// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
)

// SaasFieldSpec describes one plugin-specific field parsed from
// `pam action saas config --gateway <gw> --plugin <saasType> --info`.
type SaasFieldSpec struct {
	// Label is the custom field label practitioners must use, e.g. "REST Url".
	Label string
	// Required reflects whether the plugin's --info response listed this field as
	// "Required" or "Optional". ValidateSaasFields currently treats every field as
	// mandatory regardless of this flag (by product decision); it's kept here for
	// fidelity to the source response and potential future use.
	Required    bool
	Description string
}

// saasFieldSectionHeader marks the start of the field list in the --info response; lines
// before it (plugin name, Type:, Author:, Summary:) are not field entries.
const saasFieldSectionHeader = "Fields"

// saasFieldLinePattern matches one field line from the --info output, e.g.:
//
//	"* Required: REST Url - URL endpoint."
//	"* Optional: REST Method - HTTP method. Either 'POST' or 'PUT'"
var saasFieldLinePattern = regexp.MustCompile(`^\*\s*(Required|Optional):\s*(.+?)\s*-\s*(.*)$`)

// FetchSaasFieldSpecs runs `pam action saas config --gateway <gw> --plugin <saasType> --info`
// and parses the required/optional field labels out of the response's "Fields" section. Lines
// before "Fields" (plugin name/Type/Author/Summary) are ignored. Returns an empty slice (no
// error) if the plugin has no "Fields" section at all.
//
// Commander currently returns this as a `message` array of free-text lines; this parses that
// textual format. If Commander starts returning structured JSON for this command instead, only
// this function's parsing needs to change - ValidateSaasFields and its callers are unaffected.
func FetchSaasFieldSpecs(ctx context.Context, apiManager *api.ApiManager, gateway, saasType string) ([]SaasFieldSpec, error) {
	command := fmt.Sprintf("%s --gateway %s --plugin %s --info", utils.CmdPamActionSaaSConfig, utils.QuoteShellSingle(gateway), utils.QuoteShellSingle(saasType))
	apiResp, err := apiManager.ExecuteCommand(ctx, command, ErrOpGetSaasPluginInfo)
	if err != nil {
		return nil, err
	}

	lines, err := messageToLines(apiResp.Message.String())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrOpGetSaasPluginInfo, err)
	}

	fieldsIdx := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == saasFieldSectionHeader {
			fieldsIdx = i
			break
		}
	}
	if fieldsIdx == -1 {
		return nil, nil
	}

	var specs []SaasFieldSpec
	for _, line := range lines[fieldsIdx+1:] {
		match := saasFieldLinePattern.FindStringSubmatch(strings.TrimSpace(line))
		if match == nil {
			continue
		}
		specs = append(specs, SaasFieldSpec{
			Label:       strings.TrimSpace(match[2]),
			Required:    strings.EqualFold(match[1], "Required"),
			Description: strings.TrimSpace(match[3]),
		})
	}
	return specs, nil
}

// ValidateSaasFields checks custom against specs (the plugin's field list): every spec -
// whether the plugin's --info response reported it as "Required" or "Optional" - must be
// present in custom (by label, type "text") with a value at least 1 character long after
// trimming. All of the plugin's fields are mandatory for this check; a field missing from
// custom, or present with a blank value, is an error either way.
//
// saasType is only used to name the plugin in the returned error message - it's assumed
// already validated (e.g. via ValidateSaasType).
//
// Custom fields with labels not in specs (e.g. the universal "SaaS Type"/"Active" fields) are
// not this function's concern.
func ValidateSaasFields(saasType string, custom []commonrecordsutils.CustomFieldModel, specs []SaasFieldSpec) error {
	type fieldInfo struct {
		fieldType string
		value     string
	}
	byLabel := make(map[string]fieldInfo, len(custom))
	for _, c := range custom {
		label := strings.TrimSpace(c.Label.ValueString())
		if label == "" {
			continue
		}
		byLabel[label] = fieldInfo{
			fieldType: strings.TrimSpace(c.Type.ValueString()),
			value:     c.Value.ValueString(),
		}
	}

	var problems []string
	for _, spec := range specs {
		example := fmt.Sprintf(`{ type = "text", label = %q, value = "<your value>" }`, spec.Label)

		info, present := byLabel[spec.Label]
		if !present {
			problems = append(problems, fmt.Sprintf("%q is missing. Add it to `custom`, e.g.: %s", spec.Label, example))
			continue
		}

		if info.fieldType != "" && info.fieldType != commonrecordsutils.FieldTypeText {
			problems = append(problems, fmt.Sprintf("%q must be of type %q (got %q). Use: %s", spec.Label, commonrecordsutils.FieldTypeText, info.fieldType, example))
			continue
		}

		if strings.TrimSpace(info.value) == "" {
			problems = append(problems, fmt.Sprintf("%q has an empty value. Set a non-empty value, e.g.: %s", spec.Label, example))
		}
	}

	if len(problems) == 0 {
		return nil
	}
	return fmt.Errorf(
		"for SaaS Type (Plugin) %q, the following custom field(s) are missing or invalid:\n  - %s",
		saasType, strings.Join(problems, "\n  - "),
	)
}
