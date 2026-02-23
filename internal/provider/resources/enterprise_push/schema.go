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
		Description: "One-time action resource that pushes JSON file content to user vaults via the enterprise-push Commander CLI. Write-only: the API does not support read or delete. Adding email/team triggers Update and pushes only to newly added targets; removing email/team only updates state. File path, file content, and managed_company changes trigger replace.",
		MarkdownDescription: "One-time action resource that pushes JSON file content to user vaults via the **enterprise-push** Commander CLI.\n\n" +
			"**Write-only resource.** The API does not support read or delete. **email** and **team** are updatable: adding emails/teams triggers **Update** and pushes only to **newly added** targets; removing emails/teams only updates state (no push). Changes to **file_path**, file content, or **managed_company** trigger replace (destroy + create) and a full push.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Deterministic ID: sha256(content + sorted emails + sorted teams + managed_company). Same inputs → same ID (no re-push); any change → new ID → replace → push again.",
				MarkdownDescription: "Deterministic ID: `sha256(content + sorted emails + sorted teams + managed_company)`. Same inputs → same ID (no re-push); any change → new ID → replace → push again.",
			},
			"file_path": schema.StringAttribute{
				Required:            true,
				Description:         "Path to the JSON file on the machine where Terraform is running. The file content is read and sent as filedata to the enterprise-push command.",
				MarkdownDescription: "Path to the JSON file on the machine where Terraform is running. The file content is read and sent as **filedata** to the enterprise-push command.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("File path", 1, false),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_content_sha256": schema.StringAttribute{
				Computed:            true,
				Description:         "SHA256 hash of the file content. Used to detect file-only changes so Terraform replaces and re-pushes when the file is edited.",
				MarkdownDescription: "SHA256 hash of the file content. Used to detect file-only changes so Terraform replaces and re-pushes when the file is edited.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Optional set of email addresses to push records to. If provided, must contain at least one value. Changes trigger Update: push runs only to newly added emails.",
				MarkdownDescription: "Optional set of email addresses to push records to. If provided, must contain at least one value. Changes trigger **Update**: push runs only to **newly added** emails.",
				Validators: []validator.Set{
					utils.SetNotEmptyValidator("Email"),
					utils.SetNoEmptyStringsValidator("Email"),
				},
			},
			"team": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Optional set of team names or IDs to push records to. If provided, must contain at least one value. Changes trigger Update: push runs only to newly added teams.",
				MarkdownDescription: "Optional set of team names or IDs to push records to. If provided, must contain at least one value. Changes trigger **Update**: push runs only to **newly added** teams.",
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
