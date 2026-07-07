// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	CustomDescription                  = "Manage custom fields for the record."
	CustomMarkdownDescription          = "Manage custom fields for the record."
	CustomTypeDescription              = "Keeper field type (e.g. text, email, secret, phone, name, date)."
	CustomTypeMarkdownDescription      = "Keeper field type (e.g. `text`, `email`, `secret`, `phone`, `name`, `date`)."
	CustomLabelDescription             = "Field label."
	CustomLabelMarkdownDescription     = "Field label."
	CustomValueDescription             = "Field value; for complex types use jsonencode(JSON) matching the Keeper field schema."
	CustomValueMarkdownDescription     = "Field value; for complex types use `jsonencode(JSON)` matching the [Keeper field schema](https://docs.keeper.io/en/keeperpam/secrets-manager/about/field-record-types)."
	CustomSensitiveDescription         = "Whether to mark the value as sensitive in Terraform state display."
	CustomSensitiveMarkdownDescription = "Whether to mark the value as sensitive in Terraform state display."

	DSCustomDescription                  = "Custom fields stored in the record's `custom` array."
	DSCustomMarkdownDescription          = "Custom fields stored in the record's `custom` array."
	DSCustomTypeDescription              = "Keeper field type (e.g. text, email, secret)."
	DSCustomTypeMarkdownDescription      = "Keeper field **type** (e.g. `text`, `email`, `secret`)."
	DSCustomValueDescription             = "Field value (JSON-encoded for complex types)."
	DSCustomValueMarkdownDescription     = "Field **value** (JSON-encoded for complex types)."
	DSCustomSensitiveDescription         = "Whether the value should be treated as sensitive."
	DSCustomSensitiveMarkdownDescription = "Whether the value should be treated as **sensitive**."
)

// CustomFieldAttributeSchema returns the list nested attribute schema for `custom` on resources.
func CustomFieldAttributeSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:            true,
		Description:         CustomDescription,
		MarkdownDescription: CustomMarkdownDescription,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:            true,
					Description:         CustomTypeDescription,
					MarkdownDescription: CustomTypeMarkdownDescription,
				},
				"label": schema.StringAttribute{
					Required:            true,
					Description:         CustomLabelDescription,
					MarkdownDescription: CustomLabelMarkdownDescription,
					Validators: []validator.String{
						utils.StringMinLengthValidator("Custom field label", 1, false),
					},
				},
				"value": schema.StringAttribute{
					Required:            true,
					Description:         CustomValueDescription,
					MarkdownDescription: CustomValueMarkdownDescription,
				},
				"sensitive": schema.BoolAttribute{
					Optional:            true,
					Computed:            true,
					Description:         CustomSensitiveDescription,
					MarkdownDescription: CustomSensitiveMarkdownDescription,
					Default:             booldefault.StaticBool(false),
				},
			},
		},
	}
}

// CustomFieldDataSourceAttributeSchema returns the list nested attribute schema for `custom` on data sources.
func CustomFieldDataSourceAttributeSchema() dschema.ListNestedAttribute {
	return dschema.ListNestedAttribute{
		Computed:            true,
		Description:         DSCustomDescription,
		MarkdownDescription: DSCustomMarkdownDescription,
		NestedObject: dschema.NestedAttributeObject{
			Attributes: map[string]dschema.Attribute{
				"type": dschema.StringAttribute{
					Computed:            true,
					Description:         DSCustomTypeDescription,
					MarkdownDescription: DSCustomTypeMarkdownDescription,
				},
				"label": dschema.StringAttribute{
					Computed:            true,
					Description:         CustomLabelDescription,
					MarkdownDescription: CustomLabelMarkdownDescription,
				},
				"value": dschema.StringAttribute{
					Computed:            true,
					Sensitive:           true,
					Description:         DSCustomValueDescription,
					MarkdownDescription: DSCustomValueMarkdownDescription,
				},
				"sensitive": dschema.BoolAttribute{
					Computed:            true,
					Description:         DSCustomSensitiveDescription,
					MarkdownDescription: DSCustomSensitiveMarkdownDescription,
				},
			},
		},
	}
}

// ParseCustomFields maps API `custom` array onto Terraform models.
// The API returns field values as JSON arrays (e.g. ["value"] or [{"key":"val"}]).
// We extract the first element so the state matches what the user writes in HCL.
func ParseCustomFields(raw []utils.VaultRecordFieldResponse) []CustomFieldModel {
	if len(raw) == 0 {
		return nil
	}
	out := make([]CustomFieldModel, 0, len(raw))
	for i := range raw {
		f := &raw[i]
		cf := CustomFieldModel{
			Type:      types.StringValue(f.Type),
			Label:     types.StringValue(f.Label),
			Sensitive: types.BoolValue(fieldTypeSensitive(f.Type)),
		}
		cf.Value = types.StringValue(extractValueFromArray(f.Value))
		out = append(out, cf)
	}
	return out
}

// extractValueFromArray converts a JSON array value from the API into a string
// suitable for Terraform state comparison with user-provided HCL values.
//
// Single-element arrays are unwrapped:
//   - ["hello"]    → "hello"
//   - [{"k":"v"}]  → {"k":"v"}
//
// Multi-element arrays are returned as compact JSON:
//   - ["a","b"]           → ["a","b"]
//   - [{"k":"v"},{"x":1}] → [{"k":"v"},{"x":1}]
func extractValueFromArray(raw json.RawMessage) string {
	s := strings.TrimSpace(string(raw))
	if s == "" || s == "null" || s == "[]" {
		return ""
	}

	var arr []json.RawMessage
	if err := json.Unmarshal(raw, &arr); err != nil || len(arr) == 0 {
		return s
	}

	if len(arr) == 1 {
		return unwrapSingleElement(arr[0])
	}

	// Multiple elements: return the full array as compact JSON.
	var compact bytes.Buffer
	if err := json.Compact(&compact, raw); err == nil {
		return compact.String()
	}
	return s
}

// unwrapSingleElement extracts a single JSON value from raw bytes.
// Quoted strings are unquoted; objects/numbers/bools are returned as compact JSON.
func unwrapSingleElement(elem json.RawMessage) string {
	var str string
	if err := json.Unmarshal(elem, &str); err == nil {
		return str
	}

	var compact bytes.Buffer
	if err := json.Compact(&compact, elem); err == nil {
		return compact.String()
	}
	return strings.TrimSpace(string(elem))
}

func fieldTypeSensitive(t string) bool {
	switch t {
	case FieldTypePassword, FieldTypeSecret, FieldTypeNote, FieldTypeOTP, FieldTypeOneTimeCode, FieldTypePaymentCard, FieldTypeKeyPair:
		return true
	default:
		return false
	}
}

// customFieldKey returns Commander argument key c.<type>.<label> (no leading segment).
func customFieldKey(typ, label string) string {
	return "c." + typ + "." + label
}

// AppendCustomFieldsAdd appends custom field arguments for record-add.
func AppendCustomFieldsAdd(parts *[]string, custom []CustomFieldModel) {
	for i := range custom {
		c := &custom[i]
		if c.Type.IsNull() || c.Type.IsUnknown() || c.Label.IsNull() || c.Label.IsUnknown() {
			continue
		}
		if c.Value.IsNull() || c.Value.IsUnknown() {
			continue
		}
		key := customFieldKey(c.Type.ValueString(), c.Label.ValueString())
		val := strings.TrimSpace(c.Value.ValueString())
		if val == "" {
			continue
		}
		*parts = append(*parts, FormatFieldAssignment(key, val, customValueNeedsJSON(c.Type.ValueString(), val)))
	}
}

// AppendCustomFieldsUpdate appends changed custom field arguments for record-update.
func AppendCustomFieldsUpdate(parts *[]string, plan, state []CustomFieldModel) {
	// Simple approach: compare serialized list; if custom block changed, re-emit all plan custom fields.
	if customFieldsEqual(plan, state) {
		return
	}
	for i := range plan {
		c := &plan[i]
		if c.Type.IsNull() || c.Type.IsUnknown() || c.Label.IsNull() || c.Label.IsUnknown() {
			continue
		}
		if c.Value.IsNull() || c.Value.IsUnknown() {
			*parts = append(*parts, FormatFieldAssignment(customFieldKey(c.Type.ValueString(), c.Label.ValueString()), "", false))
			continue
		}
		val := strings.TrimSpace(c.Value.ValueString())
		key := customFieldKey(c.Type.ValueString(), c.Label.ValueString())
		*parts = append(*parts, FormatFieldAssignment(key, val, customValueNeedsJSON(c.Type.ValueString(), val)))
	}
}

func customValueNeedsJSON(fieldType, val string) bool {
	switch fieldType {
	case FieldTypeName, FieldTypePhone, FieldTypeAddress, FieldTypeHost, FieldTypePaymentCard, FieldTypeBankAccount, FieldTypeSecurityQuestion, FieldTypeKeyPair:
		return true
	default:
		if strings.HasPrefix(strings.TrimSpace(val), "{") || strings.HasPrefix(strings.TrimSpace(val), "[") {
			return true
		}
		return false
	}
}

// CustomFieldsEqual returns true when custom field slices are equivalent.
func CustomFieldsEqual(a, b []CustomFieldModel) bool {
	return customFieldsEqual(a, b)
}

func customFieldsEqual(a, b []CustomFieldModel) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Type.Equal(b[i].Type) || !a[i].Label.Equal(b[i].Label) || !a[i].Value.Equal(b[i].Value) {
			return false
		}
	}
	return true
}

// NormalizeCustomFromPlan coerces nulls for optional sensitive default.
func NormalizeCustomFromPlan(custom []CustomFieldModel) []CustomFieldModel {
	if len(custom) == 0 {
		return nil
	}
	out := make([]CustomFieldModel, len(custom))
	for i := range custom {
		out[i] = custom[i]
		if out[i].Sensitive.IsNull() || out[i].Sensitive.IsUnknown() {
			out[i].Sensitive = types.BoolValue(fieldTypeSensitive(out[i].Type.ValueString()))
		}
	}
	return out
}

// MustCompactJSON trims JSON for comparison (optional).
func MustCompactJSON(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return s
	}
	b, err := json.Marshal(v)
	if err != nil {
		return s
	}
	return string(b)
}
