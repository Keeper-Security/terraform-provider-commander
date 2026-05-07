// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Update is a no-op for this one-time action resource. All attributes use RequiresReplace,
// so Terraform will replace (destroy + create) on any change; Update is not called in that flow.
func (r *EnterpriseScimPushResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnterpriseScimPushResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
