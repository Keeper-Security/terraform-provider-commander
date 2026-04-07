// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *SecretsManagerAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SecretsManagerAppResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError("Provider Configuration Error", err.Error())
		return
	}

	command := fmt.Sprintf("%s remove '%s' --force", CmdPrefix, state.Id.ValueString())
	_, err := r.ApiManager.ExecuteCommand(ctx, command, "Unable to remove Secrets Manager application")
	if err != nil {
		resp.Diagnostics.AddError("Delete Secrets Manager App Failed", err.Error())
		return
	}
}
