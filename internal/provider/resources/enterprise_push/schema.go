// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisepush

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterprisePushResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         ResourceDescription,
		MarkdownDescription: ResourceMarkdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         "Deterministic ID: sha256(content + sorted emails + sorted teams + managed_company). Same inputs → same ID (no re-push); any change → new ID → replace → push again.",
				MarkdownDescription: "Deterministic ID: `sha256(content + sorted emails + sorted teams + managed_company)`. Same inputs → same ID (no re-push); any change → new ID → replace → push again.",
			},
			"file_path": schema.StringAttribute{
				Required:            true,
				Description:         "Path to the file with template records. File must be JSON format. File must be located on the machine where Terraform is running.",
				MarkdownDescription: "Path to the file with template records. File must be JSON format. File must be located on the machine where Terraform is running.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("File path", 1, false),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_content_sha256": schema.StringAttribute{
				Computed:            true,
				Description:         "Used to detect file-only changes so Terraform replaces and re-pushes when the file is edited.",
				MarkdownDescription: "Used to detect file-only changes so Terraform replaces and re-pushes when the file is edited.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Users to assign records to.",
				MarkdownDescription: "Users to assign records to.",
				Validators: []validator.Set{
					utils.SetNotEmptyValidator("Email"),
					utils.SetNoEmptyStringsValidator("Email"),
				},
			},
			"team": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Teams to assign records to.",
				MarkdownDescription: "Teams to assign records to.",
				Validators: []validator.Set{
					utils.SetNotEmptyValidator("Team"),
					utils.SetNoEmptyStringsValidator("Team"),
				},
			},
			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         utils.EnterpriseManagedCompanySchemaAttributeDescription,
				MarkdownDescription: utils.EnterpriseManagedCompanySchemaAttributeMarkdownDescription,
				Validators: []validator.String{
					utils.ManagedCompanyValidator,
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}
