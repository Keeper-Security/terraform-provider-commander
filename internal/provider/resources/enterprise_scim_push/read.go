// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Read is a no-op: the API does not support read. State is left unchanged.
func (r *EnterpriseScimPushResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnterpriseScimPushResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
