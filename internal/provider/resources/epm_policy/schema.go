// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import (
	"context"

	commonepm "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/epm_policy"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EpmPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Creates and manages an EPM (Endpoint Policy Management) policy.<br><br>" + "Endpoint Privilege Manager can apply least privilege policies to applications, users and machines across the fleet of endpoints which are running the Keeper agent. Policies can be applied to any collections in the tenant.<br><br>" + "For more information, see https://docs.keeper.io/en/keeperpam/endpoint-privilege-manager/policies.",
		MarkdownDescription: "Creates and manages an **EPM (Endpoint Policy Management) policy**.<br><br>" + "Endpoint Privilege Manager can apply least privilege policies to applications, users and machines across the fleet of endpoints which are running the Keeper agent. Policies can be applied to any collections in the tenant.<br><br>" + "For more information, see [EPM Policies](https://docs.keeper.io/en/keeperpam/endpoint-privilege-manager/policies).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "EPM policy ID assigned by Keeper after the policy is created.",
				MarkdownDescription: "EPM policy **ID** assigned by Keeper after the policy is created.",
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         utils.EnterpriseManagedCompanySchemaAttributeDescription,
				MarkdownDescription: utils.EnterpriseManagedCompanySchemaAttributeMarkdownDescription,
				Validators: []validator.String{
					utils.ManagedCompanyValidator,
				},
			},
			"policy_name": schema.StringAttribute{
				Required:            true,
				Description:         "Set the display name for the EPM policy.",
				MarkdownDescription: "Set the **display name** for the EPM policy.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Policy name", 1, false),
				},
			},
			"policy_type": schema.StringAttribute{
				Required:            true,
				Description:         "Type of policy. One of: " + commonepm.PolicyTypeDescription() + ".",
				MarkdownDescription: "Type of policy. One of: " + commonepm.PolicyTypeMarkdown() + ".",
				Validators: []validator.String{
					commonepm.PolicyTypeValidator{},
				},
			},
			"status": schema.StringAttribute{
				Required:            true,
				Description:         "Policy status. One of: " + commonepm.StatusDescription() + ".<br>For " + commonepm.PolicyTypeLeastPrivilege + ", only " + commonepm.StatusDescriptionForLeastPrivilege() + " allowed. Required for some policy types (e.g. enforce requires control for " + commonepm.PolicyTypeElevation + ", " + commonepm.PolicyTypeFileAccess + ", " + commonepm.PolicyTypeCommand + ").",
				MarkdownDescription: "Policy **status**. One of: " + commonepm.StatusMarkdown() + ".<br>For **" + commonepm.PolicyTypeLeastPrivilege + "**, only " + commonepm.StatusMarkdownForLeastPrivilege() + " allowed.<br>Required for some policy types (e.g. enforce requires control for " + commonepm.PolicyTypeElevation + ", " + commonepm.PolicyTypeFileAccess + ", " + commonepm.PolicyTypeCommand + ").",
				Validators: []validator.String{
					commonepm.StatusValidator{},
				},
			},
			"control": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Control actions. Set of values, each one of: " + commonepm.ControlDescription() + ". At least one required when status is enforce for " + commonepm.PolicyTypeElevation + ", " + commonepm.PolicyTypeFileAccess + ", or " + commonepm.PolicyTypeCommand + ". Not allowed for " + commonepm.PolicyTypeLeastPrivilege + ".",
				MarkdownDescription: "**Control** actions. Set of values, each one of: " + commonepm.ControlMarkdown() + ". At least one required when status is **enforce** for " + commonepm.PolicyTypeElevation + ", " + commonepm.PolicyTypeFileAccess + ", or " + commonepm.PolicyTypeCommand + ". Not allowed for **" + commonepm.PolicyTypeLeastPrivilege + "**.",
				Validators: []validator.Set{
					commonepm.ControlSetValidator{},
				},
			},
			"user_groups": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "User groups: \"*\" to select all users, or a set of user collection IDs. Cannot use \"*\" with other IDs.<br> Required (with machine_collections and/or applications) for " + commonepm.PolicyTypeElevation + "/" + commonepm.PolicyTypeFileAccess + " monitor/monitor_and_notify; required with machine_collections for " + commonepm.PolicyTypeCommand + " monitor/monitor_and_notify. Not allowed for " + commonepm.PolicyTypeLeastPrivilege + ".",
				MarkdownDescription: "**User groups**: `\"*\"` to select all users, or a set of user collection IDs. Cannot use **\"*\"** with other IDs.<br> Required (with machine_collections and/or applications) for " + commonepm.PolicyTypeElevation + "/" + commonepm.PolicyTypeFileAccess + " monitor/monitor_and_notify; required with machine_collections for " + commonepm.PolicyTypeCommand + " monitor/monitor_and_notify. Not allowed for **" + commonepm.PolicyTypeLeastPrivilege + "**.",
				Validators: []validator.Set{
					commonepm.CollectionSetValidator{DisplayName: "User group"},
				},
			},
			"machine_collections": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Machine collections: \"*\" to select all machines, or a set of machine collection IDs. Cannot use \"*\" with other IDs.<br> Required (with user_groups and/or applications) for " + commonepm.PolicyTypeElevation + "/" + commonepm.PolicyTypeFileAccess + " monitor/monitor_and_notify; required with user_groups for " + commonepm.PolicyTypeCommand + " monitor/monitor_and_notify. Optional for " + commonepm.PolicyTypeLeastPrivilege + ".",
				MarkdownDescription: "**Machine collections**: `\"*\"` to select all machines, or a set of machine collection IDs. Cannot use **\"*\"** with other IDs.<br> Required (with user_groups and/or applications) for " + commonepm.PolicyTypeElevation + "/" + commonepm.PolicyTypeFileAccess + " monitor/monitor_and_notify; required with user_groups for " + commonepm.PolicyTypeCommand + " monitor/monitor_and_notify. Optional for **" + commonepm.PolicyTypeLeastPrivilege + "**.",
				Validators: []validator.Set{
					commonepm.CollectionSetValidator{DisplayName: "Machine collection"},
				},
			},
			"applications": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Applications: \"*\" to select all applications, or a set of application collection IDs. Cannot use \"*\" with other IDs.<br> Required (with user_groups and/or machine_collections) for " + commonepm.PolicyTypeElevation + "/" + commonepm.PolicyTypeFileAccess + " monitor/monitor_and_notify. Not allowed for " + commonepm.PolicyTypeCommand + " policy type.",
				MarkdownDescription: "**Applications**: `\"*\"` to select all applications, or a set of application collection IDs. Cannot use **\"*\"** with other IDs.<br> Required (with user_groups and/or machine_collections) for " + commonepm.PolicyTypeElevation + "/" + commonepm.PolicyTypeFileAccess + " monitor/monitor_and_notify. Not allowed for **" + commonepm.PolicyTypeCommand + "** policy type.",
				Validators: []validator.Set{
					commonepm.CollectionSetValidator{DisplayName: "Application"},
				},
			},
			"day_filter": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Day filter. Set of days, each one of: " + commonepm.DayFilterDescription() + " (case-insensitive). Not allowed for " + commonepm.PolicyTypeLeastPrivilege + " policy type.",
				MarkdownDescription: "**Day filter**. Set of days, each one of: " + commonepm.DayFilterMarkdown() + " (case-insensitive). Not allowed for **" + commonepm.PolicyTypeLeastPrivilege + "** policy type.",
				Validators: []validator.Set{
					commonepm.DayFilterSetValidator{},
				},
			},
			"time_filter": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Policy time filter. Set of hour ranges as start-end (hours 0–23), e.g. \"9-12\". Ranges must not overlap. Not allowed for " + commonepm.PolicyTypeLeastPrivilege + " policy type.",
				MarkdownDescription: "**Time filter**. Set of hour ranges as **start-end** (hours **0–23**), e.g. `9-12`. Ranges must **not overlap**. Not allowed for **" + commonepm.PolicyTypeLeastPrivilege + "** policy type.",
				Validators: []validator.Set{
					commonepm.TimeFilterSetValidator{},
				},
			},
			"date_filter": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Policy date filter. Set of date ranges in ISO format YYYY-MM-DD:YYYY-MM-DD. Ranges must not overlap. Not allowed for " + commonepm.PolicyTypeLeastPrivilege + " policy type.",
				MarkdownDescription: "**Date filter**. Set of date ranges in **ISO format** `YYYY-MM-DD:YYYY-MM-DD`. Ranges must **not overlap**. Not allowed for **" + commonepm.PolicyTypeLeastPrivilege + "** policy type.",
				Validators: []validator.Set{
					commonepm.DateFilterSetValidator{},
				},
			},
		},
	}
}
