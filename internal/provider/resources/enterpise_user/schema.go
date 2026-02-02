package enterpiseuser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *EnterpriseUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The ID of the user.",
				MarkdownDescription: "The ID of the user.",
			},
			"email": schema.StringAttribute{
				Required:            true,
				Description:         "Email address.",
				MarkdownDescription: "Email address.",
				Validators: []validator.String{
					emailValidator{},
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "Name.",
				MarkdownDescription: "Name.",
				Validators: []validator.String{
					nameValidator{},
				},
			},
			"job_title": schema.StringAttribute{
				Optional:            true,
				Description:         "Job title.",
				MarkdownDescription: "Job title.",
				Validators: []validator.String{
					jobTitleValidator{},
				},
			},
			"roles": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					utils.RolesValidator{},
				},
			},
			"teams": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					utils.TeamsValidator{},
				},
			},
			"node": schema.StringAttribute{
				Required:            true,
				Description:         "Managing Node name or ID.",
				MarkdownDescription: "Managing Node name or ID.",
				Validators: []validator.String{
					utils.NodeValidator{},
				},
			},

			"managed_company": schema.StringAttribute{
				Optional:            true,
				Description:         "Managed Company name or ID.",
				MarkdownDescription: "Managed Company name or ID.",
				Validators: []validator.String{
					utils.ManagedCompanyValidator{},
				},
			},
		},
	}
}
