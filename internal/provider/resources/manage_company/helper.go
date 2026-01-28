// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// normalizeAddOns extracts add-ons from the Set (no normalization needed - validation ensures correct format)
func normalizeAddOns(addOns types.Set) []string {
	if addOns.IsNull() || addOns.IsUnknown() {
		return nil
	}

	elements := addOns.Elements()
	result := make([]string, 0, len(elements))

	for _, elem := range elements {
		if strValue, ok := elem.(types.String); ok {
			value := strValue.ValueString()
			result = append(result, value)
		}
	}

	return result
}
