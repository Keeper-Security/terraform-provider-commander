// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils/cronvalidate"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SharedAttributes returns the PAM User resource attribute map shared between
// classic and new resources. The `rotation_settings` block is included inline
// because it is an Optional nested attribute (not a Block); callers add any
// share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
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
			Required:            true,
			Description:         LoginDescription,
			MarkdownDescription: LoginMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Login", 1, false),
			},
		},
		"password": schema.StringAttribute{
			Optional:            true,
			Sensitive:           true,
			Description:         PasswordDescription,
			MarkdownDescription: PasswordMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Password", 1, true),
			},
		},
		"folder_location": schema.StringAttribute{
			Optional:            true,
			Description:         FolderDescription,
			MarkdownDescription: FolderMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Folder Location", 1, true),
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
			Validators: []validator.String{
				utils.StringMinLengthValidator("Public Key", 1, true),
			},
		},
		"public_key": schema.StringAttribute{
			Optional:            true,
			Sensitive:           true,
			Description:         PublicKeyDescription,
			MarkdownDescription: PublicKeyMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Public Key", 1, true),
			},
		},
		"private_key_passphrase": schema.StringAttribute{
			Optional:            true,
			Sensitive:           true,
			Description:         PrivateKeyPassphraseDescription,
			MarkdownDescription: PrivateKeyPassphraseMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Private Key Passphrase", 1, true),
			},
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
				RotationProfileRequirementsValidator(),
				RotationScheduleCombinationValidator(),
			},
			Attributes: map[string]schema.Attribute{
				"rotation_profile": schema.StringAttribute{
					Optional:            true,
					Description:         RotProfileDescription,
					MarkdownDescription: RotProfileMarkdownDescription,
					Validators: []validator.String{
						utils.StringOneOfValidator("Rotation Profile", []string{RotProfileGeneral, RotProfileIAMUser, RotProfileScriptsOnly, RotProfileSaaS}, true),
					},
				},
				"configuration": schema.StringAttribute{
					Optional:            true,
					Description:         RotConfigDescription,
					MarkdownDescription: RotConfigMarkdownDescription,
					Validators: []validator.String{
						utils.StringMinLengthValidator("PAM Configuration UID", 1, true),
					},
				},
				"saas_config": schema.StringAttribute{
					Optional:            true,
					Description:         RotSaaSConfigDescription,
					MarkdownDescription: RotSaaSConfigMarkdownDescription,
					Validators: []validator.String{
						utils.StringMinLengthValidator("SaaS Config UID", 1, true),
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
					Validators: []validator.String{
						RotationScheduleJSONValidator(),
					},
				},
				"on_demand": schema.BoolAttribute{
					Optional:            true,
					Description:         RotOnDemandDescription,
					MarkdownDescription: RotOnDemandMarkdownDescription,
				},
				"use_default_rotation_schedule": schema.BoolAttribute{
					Optional:            true,
					Description:         RotUseDefaultRotationScheduleDescription,
					MarkdownDescription: RotUseDefaultRotationScheduleMarkdownDescription,
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
	}
}

// SharedDataSourceAttributes returns computed PAM User data source attributes
// shared between classic and new data sources. Callers add the lookup key
// (e.g. pam_user) and share-extension attributes separately.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"id": dschema.StringAttribute{
			Computed:            true,
			Description:         IDDescription,
			MarkdownDescription: IDMarkdownDescription,
		},
		"title": dschema.StringAttribute{
			Computed:            true,
			Description:         TitleDescription,
			MarkdownDescription: TitleMarkdownDescription,
		},
		"login": dschema.StringAttribute{
			Computed:            true,
			Description:         LoginDescription,
			MarkdownDescription: LoginMarkdownDescription,
		},
		"password": dschema.StringAttribute{
			Computed:            true,
			Sensitive:           true,
			Description:         PasswordDescription,
			MarkdownDescription: PasswordMarkdownDescription,
		},
		"folder_location": dschema.StringAttribute{
			Computed:            true,
			Description:         FolderDescription,
			MarkdownDescription: FolderMarkdownDescription,
		},
		"notes": dschema.StringAttribute{
			Computed:            true,
			Description:         NotesDescription,
			MarkdownDescription: NotesMarkdownDescription,
		},
		"distinguished_name": dschema.StringAttribute{
			Computed:            true,
			Description:         DistinguishedNameDescription,
			MarkdownDescription: DistinguishedNameMarkdownDescription,
		},
		"private_pem_key": dschema.StringAttribute{
			Computed:            true,
			Sensitive:           true,
			Description:         PrivatePEMKeyDescription,
			MarkdownDescription: PrivatePEMKeyMarkdownDescription,
		},
		"public_key": dschema.StringAttribute{
			Computed:            true,
			Sensitive:           true,
			Description:         PublicKeyDescription,
			MarkdownDescription: PublicKeyMarkdownDescription,
		},
		"private_key_passphrase": dschema.StringAttribute{
			Computed:            true,
			Sensitive:           true,
			Description:         PrivateKeyPassphraseDescription,
			MarkdownDescription: PrivateKeyPassphraseMarkdownDescription,
		},
		"connect_database": dschema.StringAttribute{
			Computed:            true,
			Description:         ConnectDatabaseDescription,
			MarkdownDescription: ConnectDatabaseMarkdownDescription,
		},
		"managed": dschema.BoolAttribute{
			Computed:            true,
			Description:         ManagedDescription,
			MarkdownDescription: ManagedMarkdownDescription,
		},
		"rotation_settings": rotationSettingsDataSourceAttribute(),
	}
}

func rotationSettingsDataSourceAttribute() dschema.SingleNestedAttribute {
	return dschema.SingleNestedAttribute{
		Computed:            true,
		Description:         RotationSettingsDescription,
		MarkdownDescription: RotationSettingsMarkdownDescription,
		Attributes: map[string]dschema.Attribute{
			"rotation_profile": dschema.StringAttribute{
				Computed:            true,
				Description:         RotProfileDescription,
				MarkdownDescription: RotProfileMarkdownDescription,
			},
			"configuration": dschema.StringAttribute{
				Computed:            true,
				Description:         RotConfigDescription,
				MarkdownDescription: RotConfigMarkdownDescription,
			},
			"resource": dschema.StringAttribute{
				Computed:            true,
				Description:         RotResourceDescription,
				MarkdownDescription: RotResourceMarkdownDescription,
			},
			"saas_config": dschema.StringAttribute{
				Computed:            true,
				Description:         RotSaaSConfigDescription,
				MarkdownDescription: RotSaaSConfigMarkdownDescription,
			},
			"enabled": dschema.BoolAttribute{
				Computed:            true,
				Description:         RotEnabledDescription,
				MarkdownDescription: RotEnabledMarkdownDescription,
			},
			"schedule_cron": dschema.StringAttribute{
				Computed:            true,
				Description:         RotScheduleCronDescription,
				MarkdownDescription: RotScheduleCronMarkdownDescription,
			},
			"schedule_json": dschema.StringAttribute{
				Computed:            true,
				Description:         RotScheduleJSONDescription,
				MarkdownDescription: RotScheduleJSONMarkdownDescription,
			},
			"on_demand": dschema.BoolAttribute{
				Computed:            true,
				Description:         RotOnDemandDescription,
				MarkdownDescription: RotOnDemandMarkdownDescription,
			},
			"use_default_rotation_schedule": dschema.BoolAttribute{
				Computed:            true,
				Description:         RotUseDefaultRotationScheduleDescription,
				MarkdownDescription: RotUseDefaultRotationScheduleMarkdownDescription,
			},
			"complexity": dschema.StringAttribute{
				Computed:            true,
				Description:         RotComplexityDescription,
				MarkdownDescription: RotComplexityMarkdownDescription,
			},
		},
	}
}
