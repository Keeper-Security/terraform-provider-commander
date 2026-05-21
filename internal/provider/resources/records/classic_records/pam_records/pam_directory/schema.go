// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamdirectory

import (
	"context"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records"
	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records/pam_directory"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamDirectoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages PAM directory record with pam settings in your Keeper vault.\n\n" +
			"A PAM Directory record is a type of KeeperPAM resource that represents an Active Directory or OpenLDAP service, either on-prem or hosted in the cloud..\n\n" +
			"For more information, see the https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-directory.",
		MarkdownDescription: "Creates and manages **PAM directory record with pam settings** in your Keeper vault.\n\n" +
			"A PAM Directory record is a type of KeeperPAM resource that represents an Active Directory or OpenLDAP service, either on-prem or hosted in the cloud.\n\n" +
			"For more information, see the [PAM Directory documentation](https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-directory).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         commonpamdirectory.IDDescription,
				MarkdownDescription: commonpamdirectory.IDMarkdownDescription,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:            true,
				Description:         commonpamdirectory.TitleDescription,
				MarkdownDescription: commonpamdirectory.TitleMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Title", 1, false),
				},
			},
			"hostname_or_ip": schema.SingleNestedAttribute{
				Required:            true,
				Description:         commonpamdirectory.HostnameOrIPDescription,
				MarkdownDescription: commonpamdirectory.HostnameOrIPMarkdownDescription,
				Attributes: map[string]schema.Attribute{
					"hostname": schema.StringAttribute{
						Required:            true,
						Description:         commonpamdirectory.HostNameDescription,
						MarkdownDescription: commonpamdirectory.HostNameMarkdownDescription,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Host Name", 1, false),
						},
					},
					"administrative_port": schema.Int32Attribute{
						Optional:            true,
						Description:         commonpamdirectory.PortDescription,
						MarkdownDescription: commonpamdirectory.PortMarkdownDescription,
						Validators: []validator.Int32{
							utils.Int32NonNegativeValidator("Administrative Port", true),
						},
					},
				},
			},
			"use_ssl": schema.BoolAttribute{
				Optional:            true,
				Description:         commonpamdirectory.UseSSLDescription,
				MarkdownDescription: commonpamdirectory.UseSSLMarkdownDescription,
			},
			"domain_name": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamdirectory.DomainNameDescription,
				MarkdownDescription: commonpamdirectory.DomainNameMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Domain Name", 1, true),
				},
			},
			"alternative_ips": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         commonpamdirectory.AlternativeIPsDescription,
				MarkdownDescription: commonpamdirectory.AlternativeIPsMarkdownDescription,
				Validators: []validator.Set{
					utils.SetNotEmptyValidator("Alternative IPs"),
					utils.SetNoEmptyStringsValidator("Alternative IPs"),
				},
			},
			"directory_id": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamdirectory.DirectoryIdDescription,
				MarkdownDescription: commonpamdirectory.DirectoryIdMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Directory ID", 1, true),
				},
			},
			"directory_type": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamdirectory.DirectoryTypeDescription,
				MarkdownDescription: commonpamdirectory.DirectoryTypeMarkdownDescription,
				Validators: []validator.String{
					utils.StringOneOfValidator("Directory Type", []string{"active_directory", "openldap"}, true),
				},
			},
			"user_match": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamdirectory.UserMatchDescription,
				MarkdownDescription: commonpamdirectory.UserMatchMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("User Match", 1, true),
				},
			},
			"provider_group": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamdirectory.ProviderGroupDescription,
				MarkdownDescription: commonpamdirectory.ProviderGroupMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Provider Group", 1, true),
				},
			},
			"provider_region": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamdirectory.ProviderRegionDescription,
				MarkdownDescription: commonpamdirectory.ProviderRegionMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Provider Region", 1, true),
				},
			},
			"notes": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamdirectory.NotesDescription,
				MarkdownDescription: commonpamdirectory.NotesMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Notes", 0, true),
				},
			},
			"folder": schema.StringAttribute{
				Optional:            true,
				Description:         commonpamdirectory.FolderDescription,
				MarkdownDescription: commonpamdirectory.FolderMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Folder", 1, true),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"pam_settings": commonpamrecords.CommonPamSettingsBlock(),
		},
	}
}
