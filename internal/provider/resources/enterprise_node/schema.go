// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisenode

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *EnterpriseNodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Enterprise Node Name", 1, false),
				},
			},
			"parent": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Enterprise Node Parent Name", 1, true),
				},
			},
			"toggle_isolated": schema.BoolAttribute{
				Optional: true,
			},
			"managed_company": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					utils.ManagedCompanyValidator,
				},
			},
		},
	}
}
