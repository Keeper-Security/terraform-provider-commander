// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SharedAttributes returns the PAM Directory resource attribute map shared
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
		"use_ssl": schema.BoolAttribute{
			Optional:            true,
			Description:         UseSSLDescription,
			MarkdownDescription: UseSSLMarkdownDescription,
		},
		"domain_name": schema.StringAttribute{
			Optional:            true,
			Description:         DomainNameDescription,
			MarkdownDescription: DomainNameMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Domain Name", 1, true),
			},
		},
		"alternative_ips": schema.SetAttribute{
			Optional:            true,
			ElementType:         types.StringType,
			Description:         AlternativeIPsDescription,
			MarkdownDescription: AlternativeIPsMarkdownDescription,
			Validators: []validator.Set{
				utils.SetNotEmptyValidator("Alternative IPs"),
				utils.SetNoEmptyStringsValidator("Alternative IPs"),
			},
		},
		"directory_id": schema.StringAttribute{
			Optional:            true,
			Description:         DirectoryIdDescription,
			MarkdownDescription: DirectoryIdMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Directory ID", 1, true),
			},
		},
		"directory_type": schema.StringAttribute{
			Optional:            true,
			Description:         DirectoryTypeDescription,
			MarkdownDescription: DirectoryTypeMarkdownDescription,
			Validators: []validator.String{
				utils.StringOneOfValidator("Directory Type", []string{"active_directory", "openldap"}, true),
			},
		},
		"user_match": schema.StringAttribute{
			Optional:            true,
			Description:         UserMatchDescription,
			MarkdownDescription: UserMatchMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("User Match", 1, true),
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
