// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisepush

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Delete is a no-op: remove from state only. The API does not support delete; there is nothing to delete server-side for a push.
func (r *EnterprisePushResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
