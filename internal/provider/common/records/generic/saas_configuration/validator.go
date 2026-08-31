// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// requiredSaasTypeCustomFieldValidator ensures custom includes a text field labeled "SaaS Type".
type requiredSaasTypeCustomFieldValidator struct{}

func RequiredSaasTypeCustomFieldValidator() requiredSaasTypeCustomFieldValidator {
	return requiredSaasTypeCustomFieldValidator{}
}

func (requiredSaasTypeCustomFieldValidator) Description(_ context.Context) string {
	return fmt.Sprintf(`custom must include a %s field with label %q`, SaaSTypeCustomFieldType, SaaSTypeCustomFieldLabel)
}

func (requiredSaasTypeCustomFieldValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("`custom` must include a **%s** field with label **%s**.", SaaSTypeCustomFieldType, SaaSTypeCustomFieldLabel)
}

func (requiredSaasTypeCustomFieldValidator) ValidateList(_ context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		resp.Diagnostics.AddError(
			"Missing required custom field",
			fmt.Sprintf(`custom must include a %s field with label %q.`, SaaSTypeCustomFieldType, SaaSTypeCustomFieldLabel),
		)
		return
	}

	for _, elem := range req.ConfigValue.Elements() {
		if elem.IsNull() || elem.IsUnknown() {
			continue
		}
		obj, ok := elem.(types.Object)
		if !ok {
			continue
		}
		fieldType, label, ok := customFieldTypeAndLabel(obj.Attributes())
		if !ok {
			continue
		}
		if fieldType == SaaSTypeCustomFieldType && label == SaaSTypeCustomFieldLabel {
			return
		}
	}

	resp.Diagnostics.AddError(
		"Missing required custom field",
		fmt.Sprintf(`custom must include a %s field with label %q.`, SaaSTypeCustomFieldType, SaaSTypeCustomFieldLabel),
	)
}

func customFieldTypeAndLabel(attrs map[string]attr.Value) (fieldType, label string, ok bool) {
	t, tOk := attrs["type"].(types.String)
	l, lOk := attrs["label"].(types.String)
	if !tOk || !lOk || t.IsNull() || t.IsUnknown() || l.IsNull() || l.IsUnknown() {
		return "", "", false
	}
	return strings.TrimSpace(t.ValueString()), strings.TrimSpace(l.ValueString()), true
}
