// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *PamConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         "Use this data source to look up a PAM configuration by UID.",
		MarkdownDescription: "Use this data source to look up a **PAM configuration** by **UID**.",
		Attributes: map[string]dschema.Attribute{
			"pam_configuration": dschema.StringAttribute{
				Required:            true,
				Description:         "PAM configuration UID to read.",
				MarkdownDescription: "PAM configuration **UID** to read.",
			},
			"id": dschema.StringAttribute{
				Computed:            true,
				Description:         "The PAM configuration UID.",
				MarkdownDescription: "The PAM configuration **UID**.",
			},
			"environment": dschema.StringAttribute{
				Computed:            true,
				Description:         "PAM configuration type. One of: local, aws, azure, gcp, domain.",
				MarkdownDescription: "PAM configuration type.",
			},
			"title": dschema.StringAttribute{
				Computed:            true,
				Description:         "Title of the PAM configuration.",
				MarkdownDescription: "Title of the PAM configuration.",
			},
			"gateway": dschema.StringAttribute{
				Computed:            true,
				Description:         "The configured gateway UID or name.",
				MarkdownDescription: "The configured gateway `UID` or `name`.",
			},
			"application_folder": dschema.StringAttribute{
				Computed:            true,
				Description:         "The shared folder name or UID where the PAM Configuration data is stored.",
				MarkdownDescription: "The shared folder `name` or `UID` where the PAM Configuration data is stored.",
			},
			"schedule": dschema.StringAttribute{
				Computed:            true,
				Description:         "Rotation schedule using CRON syntax.",
				MarkdownDescription: "Rotation schedule using `CRON` syntax.",
			},
			"port_mapping": dschema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Alternative default port mappings.",
			},
			"connections": dschema.BoolAttribute{
				Computed:    true,
				Description: "Whether connections are enabled.",
			},
			"tunneling": dschema.BoolAttribute{
				Computed:    true,
				Description: "Whether tunneling is enabled.",
			},
			"rotation": dschema.BoolAttribute{
				Computed:    true,
				Description: "Whether rotation is enabled.",
			},
			"remote_browser_isolation": dschema.BoolAttribute{
				Computed:    true,
				Description: "Whether remote browser isolation is enabled.",
			},
			"connections_recording": dschema.BoolAttribute{
				Computed:    true,
				Description: "Whether connections recording is enabled.",
			},
			"typescript_recording": dschema.BoolAttribute{
				Computed:    true,
				Description: "Whether typescript recording is enabled.",
			},
			"ai_threat_detection": dschema.BoolAttribute{
				Computed:    true,
				Description: "Whether AI threat detection is enabled.",
			},
			"ai_terminate_session_on_detection": dschema.BoolAttribute{
				Computed:    true,
				Description: "Whether AI session termination on threat detection is enabled.",
			},
			"local_network": dschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Local network environment configuration.",
				Attributes: map[string]dschema.Attribute{
					"network_id": dschema.StringAttribute{
						Computed:    true,
						Description: "Unique ID for the network.",
					},
					"network_cidr": dschema.StringAttribute{
						Computed:    true,
						Description: "Subnet of the IP address.",
					},
				},
			},
			"aws": dschema.SingleNestedAttribute{
				Computed:    true,
				Description: "AWS environment configuration.",
				Attributes: map[string]dschema.Attribute{
					"aws_id": dschema.StringAttribute{
						Computed:    true,
						Description: "A unique id for the instance of AWS.",
					},
					"access_key_id": dschema.StringAttribute{
						Computed:    true,
						Description: "The Access key ID.",
					},
					"access_secret_key": dschema.StringAttribute{
						Computed:    true,
						Sensitive:   true,
						Description: "The secret key for the access key.",
					},
					"region_names": dschema.SetAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "AWS region names used for discovery.",
					},
				},
			},
			"azure": dschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Azure environment configuration.",
				Attributes: map[string]dschema.Attribute{
					"azure_id": dschema.StringAttribute{
						Computed:    true,
						Description: "A unique id for your instance of Azure.",
					},
					"client_id": dschema.StringAttribute{
						Computed:    true,
						Description: "The application/client id of the Azure application.",
					},
					"client_secret": dschema.StringAttribute{
						Computed:    true,
						Sensitive:   true,
						Description: "The client credentials secret for the Azure application.",
					},
					"subscription_id": dschema.StringAttribute{
						Computed:    true,
						Description: "The UUID of the subscription.",
					},
					"tenant_id": dschema.StringAttribute{
						Computed:    true,
						Description: "The UUID of the Azure Active Directory.",
					},
					"resource_groups": dschema.SetAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "A list of resource groups to be checked.",
					},
				},
			},
			"domain": dschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Domain environment configuration.",
				Attributes: map[string]dschema.Attribute{
					"domain_id": dschema.StringAttribute{
						Computed:    true,
						Description: "The FQDN domain used by the Domain Controller.",
					},
					"domain_hostname": dschema.StringAttribute{
						Computed:    true,
						Description: "Hostname for the domain controller.",
					},
					"domain_port": dschema.StringAttribute{
						Computed:    true,
						Description: "Port for the domain controller.",
					},
					"domain_use_ssl": dschema.BoolAttribute{
						Computed:    true,
						Description: "Whether LDAPS is used.",
					},
					"domain_scan_dc_cidr": dschema.BoolAttribute{
						Computed:    true,
						Description: "Whether CIDRs from the domain controller are scanned.",
					},
					"domain_network_cidr": dschema.StringAttribute{
						Computed:    true,
						Description: "Additional CIDRs to scan.",
					},
					"domain_admin": dschema.StringAttribute{
						Computed:    true,
						Sensitive:   true,
						Description: "Domain Administrative Credentials.",
					},
					"user_match": dschema.StringAttribute{
						Computed:    true,
						Description: "OU/DN filter (or regex) that limits which Active Directory or OpenLDAP users the Gateway discovers and imports during PAM Discovery.",
					},
				},
			},
			"gcp": dschema.SingleNestedAttribute{
				Computed:    true,
				Description: "GCP environment configuration.",
				Attributes: map[string]dschema.Attribute{
					"gcp_id": dschema.StringAttribute{
						Computed:    true,
						Description: "A unique id for the instance of Google Cloud.",
					},
					"service_account_key": dschema.StringAttribute{
						Computed:    true,
						Sensitive:   true,
						Description: "The service account key in JSON format.",
					},
					"google_admin_email": dschema.StringAttribute{
						Computed:    true,
						Description: "The email address for a Google Workspace administrator account.",
					},
					"gcp_region": dschema.StringAttribute{
						Computed:    true,
						Description: "GCP region names used for discovery.",
					},
				},
			},
		},
	}
}
