// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *SecretsManagerAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError("Provider Configuration Error", err.Error())
		return
	}

	state := SecretsManagerAppResourceModel{
		Id:       types.StringValue(req.ID),
		Name:     types.StringNull(),
		Shares:   types.SetNull(types.ObjectType{AttrTypes: shareEntryAttrTypes}),
		AppUsers: types.SetNull(types.StringType),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
