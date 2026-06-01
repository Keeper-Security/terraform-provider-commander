// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SharedAttributes returns the PAM Machine resource attribute map shared
// between classic and new resources. Callers add the `pam_settings` block and
// any share-extension attribute separately.
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
		"hostname_or_ip": schema.SingleNestedAttribute{
			Required:            true,
			Description:         HostnameOrIPDescription,
			MarkdownDescription: HostnameOrIPMarkdownDescription,
			Attributes: map[string]schema.Attribute{
				"hostname": schema.StringAttribute{
					Required:            true,
					Description:         HostNameDescription,
					MarkdownDescription: HostNameMarkdownDescription,
					Validators: []validator.String{
						utils.StringMinLengthValidator("Host Name", 1, false),
					},
				},
				"administrative_port": schema.Int32Attribute{
					Optional:            true,
					Description:         PortDescription,
					MarkdownDescription: PortMarkdownDescription,
					Validators: []validator.Int32{
						utils.Int32NonNegativeValidator("Administrative Port", true),
					},
				},
			},
		},
		"operating_system": schema.StringAttribute{
			Optional:            true,
			Description:         OperatingSystemDescription,
			MarkdownDescription: OperatingSystemMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Operating System", 1, true),
			},
		},
		"instance_name": schema.StringAttribute{
			Optional:            true,
			Description:         InstanceNameDescription,
			MarkdownDescription: InstanceNameMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Instance Name", 1, true),
			},
		},
		"instance_id": schema.StringAttribute{
			Optional:            true,
			Description:         InstanceIdDescription,
			MarkdownDescription: InstanceIdMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Instance ID", 1, true),
			},
		},
		"provider_group": schema.StringAttribute{
			Optional:            true,
			Description:         ProviderGroupDescription,
			MarkdownDescription: ProviderGroupMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Provider Group", 1, true),
			},
		},
		"provider_region": schema.StringAttribute{
			Optional:            true,
			Description:         ProviderRegionDescription,
			MarkdownDescription: ProviderRegionMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Provider Region", 1, true),
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
		"folder": schema.StringAttribute{
			Optional:            true,
			Description:         FolderDescription,
			MarkdownDescription: FolderMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Folder", 1, true),
			},
		},
	}
}
