// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"

	providerutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// FolderLocationSemanticEquality keeps the prior state value when the planned
// value differs only by whitespace around "/" separators (e.g. "A / B" vs
// "A/B"), preventing spurious diffs for folder_location attributes.
func FolderLocationSemanticEquality() planmodifier.String {
	return folderLocationSemanticEquality{}
}

type folderLocationSemanticEquality struct{}

func (folderLocationSemanticEquality) Description(_ context.Context) string {
	return "Treats folder paths that differ only by whitespace around '/' separators as equal."
}

func (m folderLocationSemanticEquality) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (folderLocationSemanticEquality) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() || req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	if providerutils.NormalizeFolderPath(req.StateValue.ValueString()) ==
		providerutils.NormalizeFolderPath(req.PlanValue.ValueString()) {
		resp.PlanValue = req.StateValue
	}
}
