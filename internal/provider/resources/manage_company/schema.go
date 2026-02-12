// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managecompany

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *ManageCompanyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Managed Company Name", 1, false),
				},
				Description:         "Managed Company Name.",
				MarkdownDescription: "Managed Company Name.",
			},
			"node": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					utils.NodeValidator,
				},
				Description:         "Managing Node name or ID.",
				MarkdownDescription: "Managing Node name or ID.",
			},
			"seats": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					seatsValidator{},
				},
				Description:         "Maximum Licenses Allowed.",
				MarkdownDescription: "Maximum Licenses Allowed.",
			},
			"plan": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					planValidator{},
				},
				Description:         "Base plan. Must be one of: " + strings.Join(PlanOptions, ", "),
				MarkdownDescription: "Base plan. Must be one of: `" + strings.Join(PlanOptions, "`, `") + "`",
			},
			"file_plan": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					filePlanValidator{},
				},
				Description:         "File storage plan. Must be one of: " + strings.Join(FilePlanOptions, ", "),
				MarkdownDescription: "File storage plan. Must be one of: `" + strings.Join(FilePlanOptions, "`, `") + "`",
			},
			"add_ons": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					addOnsValidator{},
				},
				Description:         "Secure Add-Ons to apply to the Managed Company. Must be one of: " + strings.Join(GetAllValidAddOns(), ", "),
				MarkdownDescription: "Secure Add-Ons to apply to the Managed Company. Must be one of: `" + strings.Join(GetAllValidAddOns(), "`, `") + "`",
			},
		},
	}
}
