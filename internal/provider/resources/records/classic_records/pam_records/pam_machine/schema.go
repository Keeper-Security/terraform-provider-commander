// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pammachine

import (
	"context"

	commonpamrecords "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records"
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/classic_records/pam_records/pam_machine"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *PamMachineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Creates and manages PAM machine record with pam settings in your Keeper vault.\n\n" + "A PAM Machine record is a type of KeeperPAM resource that represents a workload, such as a Windows or Linux server.\n\n" + "For more information, see the [PAM Machine documentation](https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-machine).",
		MarkdownDescription: "Creates and manages **PAM machine record with pam settings** in your Keeper vault.\n\n" + "A PAM Machine record is a type of KeeperPAM resource that represents a workload, such as a Windows or Linux server.\n\n" + "For more information, see the [PAM Machine documentation](https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-resources/pam-machine).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         commonpammachine.IDDescription,
				MarkdownDescription: commonpammachine.IDMarkdownDescription,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Required:            true,
				Description:         commonpammachine.TitleDescription,
				MarkdownDescription: commonpammachine.TitleMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Title", 1, false),
				},
			},
			"hostname_or_ip": schema.SingleNestedAttribute{
				Required:            true,
				Description:         commonpammachine.HostnameOrIPDescription,
				MarkdownDescription: commonpammachine.HostnameOrIPMarkdownDescription,
				Attributes: map[string]schema.Attribute{
					"hostname": schema.StringAttribute{
						Required:            true,
						Description:         commonpammachine.HostNameDescription,
						MarkdownDescription: commonpammachine.HostNameMarkdownDescription,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Host Name", 1, false),
						},
					},
					"administrative_port": schema.Int32Attribute{
						Optional:            true,
						Description:         commonpammachine.PortDescription,
						MarkdownDescription: commonpammachine.PortMarkdownDescription,
						Validators: []validator.Int32{
							utils.Int32NonNegativeValidator("Administrative Port", true),
						},
					},
				},
			},
			"operating_system": schema.StringAttribute{
				Optional:            true,
				Description:         commonpammachine.OperatingSystemDescription,
				MarkdownDescription: commonpammachine.OperatingSystemMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Operating System", 1, true),
				},
			},
			"instance_name": schema.StringAttribute{
				Optional:            true,
				Description:         commonpammachine.InstanceNameDescription,
				MarkdownDescription: commonpammachine.InstanceNameMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Instance Name", 1, true),
				},
			},
			"instance_id": schema.StringAttribute{
				Optional:            true,
				Description:         commonpammachine.InstanceIdDescription,
				MarkdownDescription: commonpammachine.InstanceIdMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Instance ID", 1, true),
				},
			},
			"provider_group": schema.StringAttribute{
				Optional:            true,
				Description:         commonpammachine.ProviderGroupDescription,
				MarkdownDescription: commonpammachine.ProviderGroupMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Provider Group", 1, true),
				},
			},
			"provider_region": schema.StringAttribute{
				Optional:            true,
				Description:         commonpammachine.ProviderRegionDescription,
				MarkdownDescription: commonpammachine.ProviderRegionMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Provider Region", 1, true),
				},
			},
			"notes": schema.StringAttribute{
				Optional:            true,
				Description:         commonpammachine.NotesDescription,
				MarkdownDescription: commonpammachine.NotesMarkdownDescription,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Notes", 0, true),
				},
			},
			"folder": schema.StringAttribute{
				Optional:            true,
				Description:         commonpammachine.FolderDescription,
				MarkdownDescription: commonpammachine.FolderMarkdownDescription,
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
