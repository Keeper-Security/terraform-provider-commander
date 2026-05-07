// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Delete is a no-op: remove from state only. The API does not support delete.
func (r *EnterpriseScimPushResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
