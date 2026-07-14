// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import "github.com/hashicorp/terraform-plugin-framework/types"

// PamLocalNetworkModel maps to `pam config new|edit` flags for environment=local.
type PamLocalNetworkModel struct {
	NetworkId   types.String `tfsdk:"network_id"`
	NetworkCidr types.String `tfsdk:"network_cidr"`
}

// PamAwsModel maps to AWS-specific CLI flags.
type PamAwsModel struct {
	AwsId           types.String `tfsdk:"aws_id"`
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	AccessSecretKey types.String `tfsdk:"access_secret_key"`
	RegionNames     types.Set    `tfsdk:"region_names"`
}

// PamAzureModel maps to Azure-specific CLI flags.
type PamAzureModel struct {
	AzureId        types.String `tfsdk:"azure_id"`
	ClientId       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
	TenantId       types.String `tfsdk:"tenant_id"`
	ResourceGroups types.Set    `tfsdk:"resource_groups"`
}

// PamDomainModel maps to domain-specific CLI flags.
type PamDomainModel struct {
	DomainId          types.String `tfsdk:"domain_id"`
	DomainHostname    types.String `tfsdk:"domain_hostname"`
	DomainPort        types.String `tfsdk:"domain_port"`
	DomainUseSsl      types.Bool   `tfsdk:"domain_use_ssl"`
	DomainScanDcCidr  types.Bool   `tfsdk:"domain_scan_dc_cidr"`
	DomainNetworkCidr types.String `tfsdk:"domain_network_cidr"`
	DomainAdmin       types.String `tfsdk:"domain_admin"`
	UserMatch         types.String `tfsdk:"user_match"`
}

// PamGcpModel maps to GCP-specific CLI flags.
type PamGcpModel struct {
	GcpId             types.String `tfsdk:"gcp_id"`
	ServiceAccountKey types.String `tfsdk:"service_account_key"`
	GoogleAdminEmail  types.String `tfsdk:"google_admin_email"`
	GcpRegion         types.String `tfsdk:"gcp_region"`
}

// PamConfigurationResourceModel is the Terraform state for commander_pam_configuration.
type PamConfigurationResourceModel struct {
	Id types.String `tfsdk:"id"`

	Environment       types.String `tfsdk:"environment"`
	Title             types.String `tfsdk:"title"`
	Gateway           types.String `tfsdk:"gateway"`
	ApplicationFolder types.String `tfsdk:"application_folder"`
	Schedule          types.String `tfsdk:"schedule"`
	PortMapping       types.Set    `tfsdk:"port_mapping"`

	Connections                   types.Bool `tfsdk:"connections"`
	Tunneling                     types.Bool `tfsdk:"tunneling"`
	Rotation                      types.Bool `tfsdk:"rotation"`
	RemoteBrowserIsolation        types.Bool `tfsdk:"remote_browser_isolation"`
	ConnectionsRecording          types.Bool `tfsdk:"connections_recording"`
	TypescriptRecording           types.Bool `tfsdk:"typescript_recording"`
	AIThreatDetection             types.Bool `tfsdk:"ai_threat_detection"`
	AITerminateSessionOnDetection types.Bool `tfsdk:"ai_terminate_session_on_detection"`

	LocalNetwork *PamLocalNetworkModel `tfsdk:"local_network"`
	Aws          *PamAwsModel          `tfsdk:"aws"`
	Azure        *PamAzureModel        `tfsdk:"azure"`
	Domain       *PamDomainModel       `tfsdk:"domain"`
	Gcp          *PamGcpModel          `tfsdk:"gcp"`
}
