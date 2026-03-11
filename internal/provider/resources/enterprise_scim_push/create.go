// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisescimpush

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseScimPushResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EnterpriseScimPushResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(
			utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR,
			err.Error(),
		)
		return
	}

	source := strings.TrimSpace(strings.ToLower(data.Source.ValueString()))
	if source != SourceGoogle && source != SourceAD && source != SourceRecord {
		resp.Diagnostics.AddError(
			"Invalid source",
			"Source must be one of: google, ad, record. Got: "+data.Source.ValueString(),
		)
		return
	}

	if err := utils.RunWithManagedCompanyContext(ctx, r.ApiManager, data.ManagedCompany, func() error {
		command := buildScimPushCommand(&data)
		_, err := r.ApiManager.ExecuteCommand(ctx, command, "Enterprise SCIM push failed")
		return err
	}, "Enterprise SCIM Push Failed", &resp.Diagnostics); err != nil {
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue(computeID(&data))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
