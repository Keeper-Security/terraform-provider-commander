// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
)

// SaasFieldSpec describes one plugin-specific field parsed from
// `pam action saas config --gateway <gw> --plugin <saasType> --info --format json`.
type SaasFieldSpec struct {
	Label string
	// Type is the Keeper custom field type (e.g. "text", "secret") this field's value must
	// use. ValidateSaasFields falls back to "text" if this is empty.
	Type string
}

// saasFieldSpecEntry is one element of the required_fields/optional_fields arrays, e.g.:
//
//	{"label": "User Pool ID", "type": "secret", "required": true, "description": "User Pool ID."}
//
// required/description aren't captured: every field is mandatory regardless (by product
// decision), and the description isn't surfaced anywhere today.
type saasFieldSpecEntry struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

// saasPluginInfoResponse is the `data` object from
// `pam action saas config --gateway <gw> --plugin <saasType> --info --format json`.
type saasPluginInfoResponse struct {
	RequiredFields []saasFieldSpecEntry `json:"required_fields"`
	OptionalFields []saasFieldSpecEntry `json:"optional_fields"`
}

// FetchSaasFieldSpecs runs `pam action saas config --gateway <gw> --plugin <saasType> --info
// --format json`. Returns an empty slice (no error) if the plugin reports no fields.
func FetchSaasFieldSpecs(ctx context.Context, apiManager *api.ApiManager, gateway, saasType string) ([]SaasFieldSpec, error) {
	command := fmt.Sprintf("%s --gateway %s --plugin %s --info %s", utils.CmdPamActionSaaSConfig, utils.QuoteShellSingle(gateway), utils.QuoteShellSingle(saasType), utils.FlagFormatJSON)
	apiResp, err := apiManager.ExecuteCommand(ctx, command, ErrOpGetSaasPluginInfo)
	if err != nil {
		return nil, err
	}

	if apiResp.Data == nil {
		return nil, nil
	}

	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to process plugin info response: %w", ErrOpGetSaasPluginInfo, err)
	}

	var info saasPluginInfoResponse
	if err := json.Unmarshal(dataBytes, &info); err != nil {
		return nil, fmt.Errorf("%s: unable to parse plugin info response: %w", ErrOpGetSaasPluginInfo, err)
	}

	rawEntries := make([]saasFieldSpecEntry, 0, len(info.RequiredFields)+len(info.OptionalFields))
	rawEntries = append(rawEntries, info.RequiredFields...)
	rawEntries = append(rawEntries, info.OptionalFields...)

	specs := make([]SaasFieldSpec, 0, len(rawEntries))
	for _, e := range rawEntries {
		specs = append(specs, SaasFieldSpec{
			Label: strings.TrimSpace(e.Label),
			Type:  strings.TrimSpace(e.Type),
		})
	}
	return specs, nil
}

// ValidateSaasFields checks custom against specs: every spec (required or optional - all
// mandatory here, by product decision) must be present in custom by label, with the spec's
// Type (defaulting to "text"), and a non-blank value. saasType only names the plugin in the
// error message - it's assumed already validated via ValidateSaasType. Custom fields with
// labels not in specs (e.g. the universal "SaaS Type"/"Active" fields) are ignored.
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
		wantType := spec.Type
		if wantType == "" {
			wantType = commonrecordsutils.FieldTypeText
		}
		example := fmt.Sprintf(`{ type = %q, label = %q, value = "<your value>" }`, wantType, spec.Label)

		info, present := byLabel[spec.Label]
		if !present {
			problems = append(problems, fmt.Sprintf("%q is missing. Add it to `custom`, e.g.: %s", spec.Label, example))
			continue
		}

		if info.fieldType != "" && info.fieldType != wantType {
			problems = append(problems, fmt.Sprintf("%q must be of type %q (got %q). Use: %s", spec.Label, wantType, info.fieldType, example))
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
