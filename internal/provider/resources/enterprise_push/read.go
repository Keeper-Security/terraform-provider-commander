// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisepush

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Read is a no-op: the API does not support read. State is left unchanged so Terraform does not recreate the resource.
func (r *EnterprisePushResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnterprisePushResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
