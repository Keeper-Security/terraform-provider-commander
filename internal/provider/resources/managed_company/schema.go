// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcompany

import (
	"context"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *ManagedCompanyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages a managed company within your Keeper MSP account.<br><br>" +
			"Managed Companies (MCs) are the independent tenants that MSPs manage through the central console; each MC can be administered by the MSP (full service), by an MC administrator (reseller model), or both (hybrid model).<br><br>" +
			"This resource only works when the provider is configured for an MSP account.<br><br>" +
			"For more information, see https://docs.keeper.io/en/enterprise-guide/keeper-msp",
		MarkdownDescription: "Creates and manages a **managed company** within your Keeper MSP account.<br><br>" +
			"**Managed Companies (MCs)** are the independent tenants that MSPs manage through the central console; each MC can be administered by the MSP (full service), by an MC administrator (reseller model), or both (hybrid model).<br><br>" +
			"This resource **only works when the provider is configured for an MSP account**.<br><br>" +
			"For more information, see [Keeper MSP documentation](https://docs.keeper.io/en/enterprise-guide/keeper-msp).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Company ID assigned by Keeper to the managed company after it is created. " +
					"Use this value to import an existing managed company into Terraform state or to reference the company from other resources.",
				MarkdownDescription: "**Company ID** assigned by Keeper to the managed company after it is created. " +
					"Use this value to **import** an existing managed company into Terraform state or to reference the company from other resources.",
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Managed Company Name", 1, false),
				},
				Description:         "Set the display name for the managed company. Must be at least one character.",
				MarkdownDescription: "Set the **display name** for the managed company. Must be at least **one character**.",
			},
			"node": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					utils.NodeValidator,
				},
				Description:         "The node that will manage this managed company. Provide the node name or node ID. ",
				MarkdownDescription: "The **node** that will **manage** this managed company. Provide the **node name** or **node ID**. ",
			},
			"seats": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					seatsValidator{},
				},
				Description: "Maximum number of user licenses for this managed company. Use -1 for unlimited licenses, " +
					"0 if you do not want users to be provisioned yet, or a positive number for a fixed license cap. ",
				MarkdownDescription: "Maximum number of user **licenses** for this managed company. Use `-1` for unlimited licenses, " +
					"`0` to create the company without provisioning any users yet, or a positive number for a fixed license cap. ",
			},
			"plan": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					planValidator{},
				},
				Description:         "The Keeper base plan for this managed company. The plan determines which features and limits apply. Must be one of: " + strings.Join(PlanOptions, ", "),
				MarkdownDescription: "The Keeper **base plan** for this managed company. The plan determines which features and limits apply. Must be one of: `" + strings.Join(PlanOptions, "`, `") + "`",
			},
			"file_plan": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					filePlanValidator{},
				},
				Description: "Secure File Storage for the managed company. Must be one of: " + strings.Join(FilePlanOptions, ", ") + "<br>" +
					"Available options depend on the selected plan: " +
					"Plus plans (businessPlus, enterprisePlus) support only 1tb and 10tb, other plans support 100gb, 1tb, and 10tb.",
				MarkdownDescription: "**Secure File Storage** for the managed company. Must be one of: `" + strings.Join(FilePlanOptions, "`, `") + "`<br>" +
					"Available options depend on the selected plan: " +
					"Plus plans (`businessPlus`, `enterprisePlus`) support only `1tb` and `10tb`, other plans support `100gb`, `1tb`, and `10tb`.",
			},
			"add_ons": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					addOnsValidator{},
				},
				Description: "Secure Add-Ons products to enable for this managed company.<br>" +
					"Must be one of: " + strings.Join(GetAllValidAddOns(), ", ") + "<br><br>" +
					"Some add-ons accept a numeric suffix for quantity, such as connection_manager:2 or privileged_access_manager:1.<br>",
				MarkdownDescription: "**Secure Add-Ons** products to enable for this managed company.<br>" +
					"Must be one of: `" + strings.Join(GetAllValidAddOns(), "`, `") + "`<br><br>" +
					"Some add-ons accept a numeric suffix for quantity, such as `connection_manager:2` or `privileged_access_manager:1`.",
			},
		},
	}
}
