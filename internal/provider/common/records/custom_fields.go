// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package records

import (
	"encoding/json"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomFieldModel is one user-defined custom field (maps to API `custom` array).
type CustomFieldModel struct {
	Type      types.String `tfsdk:"type"`
	Label     types.String `tfsdk:"label"`
	Value     types.String `tfsdk:"value"`
	Sensitive types.Bool   `tfsdk:"sensitive"`
}

// CustomFieldBlockSchema returns the list nested block schema for `custom`.
func CustomFieldBlockSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		Description:         "Manage custom fields for the record.",
		MarkdownDescription: "Manage custom fields for the record.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:            true,
					Description:         "Keeper field type (e.g. text, email, secret, phone, name, date).",
					MarkdownDescription: "Keeper field type (e.g. `text`, `email`, `secret`, `phone`, `name`, `date`).",
				},
				"label": schema.StringAttribute{
					Required:            true,
					Description:         "Field label.",
					MarkdownDescription: "Field label.",
					Validators: []validator.String{
						utils.StringMinLengthValidator("Custom field label", 1, false),
					},
				},
				"value": schema.StringAttribute{
					Required:            true,
					Description:         "Field value; for complex types use JSON matching the Keeper field schema.",
					MarkdownDescription: "Field value; for complex types use JSON matching the [Keeper field schema](https://docs.keeper.io/en/keeperpam/secrets-manager/about/field-record-types).",
				},
				"sensitive": schema.BoolAttribute{
					Optional:            true,
					Computed:            true,
					Description:         "Whether to mark the value as sensitive in Terraform state display.",
					MarkdownDescription: "Whether to mark the value as sensitive in Terraform state display.",
					Default:             booldefault.StaticBool(false),
				},
			},
		},
	}
}

// ParseCustomFields maps API `custom` array onto Terraform models.
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
		// Re-serialize value as compact JSON string for stable Terraform strings.
		cf.Value = types.StringValue(string(f.Value))
		if strings.TrimSpace(cf.Value.ValueString()) == "" || cf.Value.ValueString() == "null" {
			cf.Value = types.StringValue("[]")
		}
		out = append(out, cf)
	}
	return out
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
		*parts = append(*parts, formatFieldAssignment(key, val, customValueNeedsJSON(c.Type.ValueString(), val)))
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
			*parts = append(*parts, formatFieldAssignment(customFieldKey(c.Type.ValueString(), c.Label.ValueString()), "", false))
			continue
		}
		val := strings.TrimSpace(c.Value.ValueString())
		key := customFieldKey(c.Type.ValueString(), c.Label.ValueString())
		*parts = append(*parts, formatFieldAssignment(key, val, customValueNeedsJSON(c.Type.ValueString(), val)))
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
