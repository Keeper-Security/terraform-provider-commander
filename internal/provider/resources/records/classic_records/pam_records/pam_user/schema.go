// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/cronvalidate"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *PamUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         IDDescription,
				MarkdownDescription: IDMarkdownDescription,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:            true,
				Description:         TitleDescription,
				MarkdownDescription: TitleMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Title", 1, false),
				},
			},
			"login": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         LoginDescription,
				MarkdownDescription: LoginMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Login", 1, true),
				},
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Description:         PasswordDescription,
				MarkdownDescription: PasswordMarkdownDescription,
			},
			"folder": schema.StringAttribute{
				Optional:            true,
				Description:         FolderDescription,
				MarkdownDescription: FolderMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Folder", 1, true),
				},
			},
			"notes": schema.StringAttribute{
				Optional:            true,
				Description:         NotesDescription,
				MarkdownDescription: NotesMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Notes", 0, true),
				},
			},
			"distinguished_name": schema.StringAttribute{
				Optional:            true,
				Description:         DistinguishedNameDescription,
				MarkdownDescription: DistinguishedNameMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Distinguished Name", 1, true),
				},
			},
			"private_pem_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         PrivatePEMKeyDescription,
				MarkdownDescription: PrivatePEMKeyMarkdownDescription,
			},
			"connect_database": schema.StringAttribute{
				Optional:            true,
				Description:         ConnectDatabaseDescription,
				MarkdownDescription: ConnectDatabaseMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Connect Database", 1, true),
				},
			},
			"managed": schema.BoolAttribute{
				Optional:            true,
				Description:         ManagedDescription,
				MarkdownDescription: ManagedMarkdownDescription,
			},
			"rotation_settings": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         RotationSettingsDescription,
				MarkdownDescription: RotationSettingsMarkdownDescription,
				Validators: []validator.Object{
					RotationScheduleCombinationValidator(),
				},
				Attributes: map[string]schema.Attribute{
					"rotation_profile": schema.StringAttribute{
						Optional:            true,
						Description:         RotProfileDescription,
						MarkdownDescription: RotProfileMarkdownDescription,
					},
					"configuration": schema.StringAttribute{
						Optional:            true,
						Description:         RotConfigDescription,
						MarkdownDescription: RotConfigMarkdownDescription,
						Validators: []validator.String{
							utils.StringMinLengthValidator("PAM Configuration UID", 1, true),
						},
					},
					"iam_aad_config": schema.StringAttribute{
						Optional:            true,
						Description:         RotIamAadConfigDescription,
						MarkdownDescription: RotIamAadConfigMarkdownDescription,
						Validators: []validator.String{
							utils.StringMinLengthValidator("IAM/AAD Config UID", 1, true),
						},
					},
					"resource": schema.StringAttribute{
						Optional:            true,
						Description:         RotResourceDescription,
						MarkdownDescription: RotResourceMarkdownDescription,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Resource UID", 1, true),
						},
					},
					"admin_user": schema.StringAttribute{
						Optional:            true,
						Description:         RotAdminUserDescription,
						MarkdownDescription: RotAdminUserMarkdownDescription,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Admin User UID", 1, true),
						},
					},
					"enabled": schema.BoolAttribute{
						Optional:            true,
						Description:         RotEnabledDescription,
						MarkdownDescription: RotEnabledMarkdownDescription,
					},
					"schedule_cron": schema.StringAttribute{
						Optional:            true,
						Description:         RotScheduleCronDescription,
						MarkdownDescription: RotScheduleCronMarkdownDescription,
						Validators: []validator.String{
							cronvalidate.RotationCronString{Name: "schedule_cron"},
						},
					},
					"schedule_json": schema.StringAttribute{
						Optional:            true,
						Description:         RotScheduleJSONDescription,
						MarkdownDescription: RotScheduleJSONMarkdownDescription,
					},
					"on_demand": schema.BoolAttribute{
						Optional:            true,
						Description:         RotOnDemandDescription,
						MarkdownDescription: RotOnDemandMarkdownDescription,
					},
					"schedule_config": schema.BoolAttribute{
						Optional:            true,
						Description:         RotScheduleConfigDescription,
						MarkdownDescription: RotScheduleConfigMarkdownDescription,
					},
					"complexity": schema.StringAttribute{
						Optional:            true,
						Description:         RotComplexityDescription,
						MarkdownDescription: RotComplexityMarkdownDescription,
						Validators: []validator.String{
							RotationPasswordComplexityValidator(),
						},
					},
				},
			},
		},
	}
}
