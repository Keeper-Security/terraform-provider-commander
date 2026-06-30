// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
		"folder_location": schema.StringAttribute{
			Optional:            true,
			Description:         FolderDescription,
			MarkdownDescription: FolderMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Folder Location", 1, true),
			},
		},
	}
}

// SharedDataSourceAttributes returns computed PAM Machine data source attributes
// shared between classic and new data sources. Callers add the lookup key
// (e.g. pam_machine) and share-extension attributes separately.
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
		"hostname_or_ip": hostnameOrIPDataSourceAttribute(),
		"operating_system": dschema.StringAttribute{
			Computed:            true,
			Description:         OperatingSystemDescription,
			MarkdownDescription: OperatingSystemMarkdownDescription,
		},
		"instance_name": dschema.StringAttribute{
			Computed:            true,
			Description:         InstanceNameDescription,
			MarkdownDescription: InstanceNameMarkdownDescription,
		},
		"instance_id": dschema.StringAttribute{
			Computed:            true,
			Description:         InstanceIdDescription,
			MarkdownDescription: InstanceIdMarkdownDescription,
		},
		"provider_group": dschema.StringAttribute{
			Computed:            true,
			Description:         ProviderGroupDescription,
			MarkdownDescription: ProviderGroupMarkdownDescription,
		},
		"provider_region": dschema.StringAttribute{
			Computed:            true,
			Description:         ProviderRegionDescription,
			MarkdownDescription: ProviderRegionMarkdownDescription,
		},
		"notes": dschema.StringAttribute{
			Computed:            true,
			Description:         NotesDescription,
			MarkdownDescription: NotesMarkdownDescription,
		},
		"folder_location": dschema.StringAttribute{
			Computed:            true,
			Description:         FolderDescription,
			MarkdownDescription: FolderMarkdownDescription,
		},
		"pam_settings": commonpamrecords.CommonPamSettingsDataSourceAttribute(commonpamrecords.MachineDirectoryProtocols),
	}
}

func hostnameOrIPDataSourceAttribute() dschema.SingleNestedAttribute {
	return dschema.SingleNestedAttribute{
		Computed:            true,
		Description:         HostnameOrIPDescription,
		MarkdownDescription: HostnameOrIPMarkdownDescription,
		Attributes: map[string]dschema.Attribute{
			"hostname": dschema.StringAttribute{
				Computed:            true,
				Description:         HostNameDescription,
				MarkdownDescription: HostNameMarkdownDescription,
			},
			"administrative_port": dschema.Int32Attribute{
				Computed:            true,
				Description:         PortDescription,
				MarkdownDescription: PortMarkdownDescription,
			},
		},
	}
}
