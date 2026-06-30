// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *PamConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         DescResource,
		MarkdownDescription: DescResourceMD,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description:         DescId,
				MarkdownDescription: DescIdMD,
			},
			"environment": schema.StringAttribute{
				Required:            true,
				Description:         DescEnvironment + " Allowed values: " + ValidEnvironmentsMarkdown() + ".",
				MarkdownDescription: DescEnvironmentMD + " Allowed values: " + ValidEnvironmentsMarkdown() + ".",
				Validators: []validator.String{
					environmentValidator{},
				},
			},
			"title": schema.StringAttribute{
				Required:            true,
				Description:         DescTitle,
				MarkdownDescription: DescTitleMD,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Title", 1, false),
				},
			},
			"gateway": schema.StringAttribute{
				Required:            true,
				Description:         DescGateway,
				MarkdownDescription: DescGatewayMD,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Gateway", 1, false),
				},
			},
			"application_folder": schema.StringAttribute{
				Required:            true,
				Description:         DescApplicationFolder,
				MarkdownDescription: DescApplicationFolderMD,
				Validators: []validator.String{
					utils.StringMinLengthValidator("Application Folder", 1, false),
				},
			},
			"schedule": schema.StringAttribute{
				Optional:            true,
				Description:         DescSchedule,
				MarkdownDescription: DescScheduleMD,
			},
			"port_mapping": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         DescPortMapping,
				MarkdownDescription: DescPortMappingMD,
				Validators: []validator.Set{
					utils.SetNotEmptyValidator("Port mapping"),
					utils.SetNoEmptyStringsValidator("Port mapping"),
				},
			},
			"connections": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Description:         DescConnections,
				MarkdownDescription: DescConnectionsMD,
			},
			"tunneling": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Description:         DescTunneling,
				MarkdownDescription: DescTunnelingMD,
			},
			"rotation": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Description:         DescRotation,
				MarkdownDescription: DescRotationMD,
			},
			"remote_browser_isolation": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Description:         DescRemoteBrowserIsolation,
				MarkdownDescription: DescRemoteBrowserIsolationMD,
			},
			"connections_recording": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Description:         DescConnectionsRecording,
				MarkdownDescription: DescConnectionsRecordingMD,
			},
			"typescript_recording": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Description:         DescTypescriptRecording,
				MarkdownDescription: DescTypescriptRecordingMD,
			},
			"ai_threat_detection": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Description:         DescAIThreatDetection,
				MarkdownDescription: DescAIThreatDetectionMD,
			},
			"ai_terminate_session_on_detection": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				Description:         DescAITerminateSessionOnDetection,
				MarkdownDescription: DescAITerminateSessionOnDetectionMD,
			},
		},
		Blocks: map[string]schema.Block{
			"local_network": schema.SingleNestedBlock{
				Description:         DescLocalNetwork,
				MarkdownDescription: DescLocalNetworkMD,
				Attributes: map[string]schema.Attribute{
					"network_id": schema.StringAttribute{
						Optional:            true,
						Description:         DescNetworkId,
						MarkdownDescription: DescNetworkIdMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Network ID", 1, true),
						},
					},
					"network_cidr": schema.StringAttribute{
						Optional:            true,
						Description:         DescNetworkCidr,
						MarkdownDescription: DescNetworkCidrMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Network CIDR", 1, true),
						},
					},
				},
			},
			"aws": schema.SingleNestedBlock{
				Description:         DescAws,
				MarkdownDescription: DescAwsMD,
				Attributes: map[string]schema.Attribute{
					"aws_id": schema.StringAttribute{
						Optional:            true,
						Description:         DescAwsId,
						MarkdownDescription: DescAwsIdMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("AWS ID", 1, true),
						},
					},
					"access_key_id": schema.StringAttribute{
						Optional:            true,
						Description:         DescAccessKeyId,
						MarkdownDescription: DescAccessKeyIdMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Access Key ID", 1, true),
						},
					},
					"access_secret_key": schema.StringAttribute{
						Optional:            true,
						Sensitive:           true,
						Description:         DescAccessSecretKey,
						MarkdownDescription: DescAccessSecretKeyMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Access Secret Key", 1, true),
						},
					},
					"region_names": schema.SetAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						Description:         DescRegionNames,
						MarkdownDescription: DescRegionNamesMD,
						Validators: []validator.Set{
							utils.SetNotEmptyValidator("Region names"),
							utils.SetNoEmptyStringsValidator("Region names"),
						},
					},
				},
			},
			"azure": schema.SingleNestedBlock{
				Description:         DescAzure,
				MarkdownDescription: DescAzureMD,
				Attributes: map[string]schema.Attribute{
					"azure_id": schema.StringAttribute{
						Optional:            true,
						Description:         DescAzureId,
						MarkdownDescription: DescAzureIdMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Azure ID", 1, true),
						},
					},
					"client_id": schema.StringAttribute{
						Optional:            true,
						Description:         DescClientId,
						MarkdownDescription: DescClientIdMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Client ID", 1, true),
						},
					},
					"client_secret": schema.StringAttribute{
						Optional:            true,
						Sensitive:           true,
						Description:         DescClientSecret,
						MarkdownDescription: DescClientSecretMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Client Secret", 1, true),
						},
					},
					"subscription_id": schema.StringAttribute{
						Optional:            true,
						Description:         DescSubscriptionId,
						MarkdownDescription: DescSubscriptionIdMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Subscription ID", 1, true),
						},
					},
					"tenant_id": schema.StringAttribute{
						Optional:            true,
						Description:         DescTenantId,
						MarkdownDescription: DescTenantIdMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Tenant ID", 1, true),
						},
					},
					"resource_groups": schema.SetAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						Description:         DescResourceGroups,
						MarkdownDescription: DescResourceGroupsMD,
						Validators: []validator.Set{
							utils.SetNotEmptyValidator("Resource groups"),
							utils.SetNoEmptyStringsValidator("Resource groups"),
						},
					},
				},
			},
			"domain": schema.SingleNestedBlock{
				Description:         DescDomain,
				MarkdownDescription: DescDomainMD,
				Attributes: map[string]schema.Attribute{
					"domain_id": schema.StringAttribute{
						Optional:            true,
						Description:         DescDomainId,
						MarkdownDescription: DescDomainIdMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Domain ID", 1, true),
						},
					},
					"domain_hostname": schema.StringAttribute{
						Optional:            true,
						Description:         DescDomainHostname,
						MarkdownDescription: DescDomainHostnameMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Domain Hostname", 1, true),
						},
					},
					"domain_port": schema.StringAttribute{
						Optional:            true,
						Description:         DescDomainPort,
						MarkdownDescription: DescDomainPortMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Domain Port", 1, true),
						},
					},
					"domain_use_ssl": schema.BoolAttribute{
						Optional:            true,
						Description:         DescDomainUseSsl,
						MarkdownDescription: DescDomainUseSslMD,
					},
					"domain_scan_dc_cidr": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
						Description:         DescDomainScanDcCidr,
						MarkdownDescription: DescDomainScanDcCidrMD,
					},
					"domain_network_cidr": schema.StringAttribute{
						Optional:            true,
						Description:         DescDomainNetworkCidr,
						MarkdownDescription: DescDomainNetworkCidrMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Domain Network CIDR", 1, true),
						},
					},
					"domain_admin": schema.StringAttribute{
						Optional:            true,
						Sensitive:           true,
						Description:         DescDomainAdmin,
						MarkdownDescription: DescDomainAdminMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Domain Admin", 1, true),
						},
					},
					"user_match": schema.StringAttribute{
						Optional:            true,
						Description:         DescUserMatch,
						MarkdownDescription: DescUserMatchMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("User Match", 1, true),
						},
					},
				},
			},
			"gcp": schema.SingleNestedBlock{
				Description:         DescGcp,
				MarkdownDescription: DescGcpMD,
				Attributes: map[string]schema.Attribute{
					"gcp_id": schema.StringAttribute{
						Optional:            true,
						Description:         DescGcpId,
						MarkdownDescription: DescGcpIdMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("GCP ID", 1, true),
						},
					},
					"service_account_key": schema.StringAttribute{
						Optional:            true,
						Sensitive:           true,
						Description:         DescServiceAccountKey,
						MarkdownDescription: DescServiceAccountKeyMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Service Account Key", 1, true),
							utils.JSONStringValidator("Service Account Key"),
						},
					},
					"google_admin_email": schema.StringAttribute{
						Optional:            true,
						Description:         DescGoogleAdminEmail,
						MarkdownDescription: DescGoogleAdminEmailMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("Google Admin Email", 1, true),
						},
					},
					"gcp_region": schema.StringAttribute{
						Optional:            true,
						Description:         DescGcpRegion,
						MarkdownDescription: DescGcpRegionMD,
						Validators: []validator.String{
							utils.StringMinLengthValidator("GCP Region", 1, true),
						},
					},
				},
			},
		},
	}
}
