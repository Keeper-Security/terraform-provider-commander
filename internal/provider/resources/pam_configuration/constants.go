// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import "strings"

// Commander CLI: `pam config new` / `pam config edit` and flags.
const (
	CmdPamConfig = "pam config"
	CmdPamNew    = "new"
	CmdPamEdit   = "edit"
	CmdPamRemove = "remove"
	CmdPamList   = "list"
)

const (
	FlagPamListConfig = "--config"
	FlagFormat        = "--format"
	FormatJSON        = "json"
)

const (
	FlagEnvironment                   = "--environment"
	FlagTitle                         = "--title"
	FlagGateway                       = "--gateway"
	FlagSharedFolder                  = "--shared-folder"
	FlagSchedule                      = "--schedule"
	FlagPortMapping                   = "--port-mapping"
	FlagConnections                   = "--connections"
	FlagTunneling                     = "--tunneling"
	FlagRotation                      = "--rotation"
	FlagRemoteBrowserIsolation        = "--remote-browser-isolation"
	FlagConnectionsRecording          = "--connections-recording"
	FlagTypescriptRecording           = "--typescript-recording"
	FlagAIThreatDetection             = "--ai-threat-detection"
	FlagAITerminateSessionOnDetection = "--ai-terminate-session-on-detection"

	FlagNetworkId   = "--network-id"
	FlagNetworkCidr = "--network-cidr"

	FlagAwsId           = "--aws-id"
	FlagAccessKeyId     = "--access-key-id"
	FlagAccessSecretKey = "--access-secret-key"
	FlagRegionName      = "--region-name"

	FlagAzureId        = "--azure-id"
	FlagClientId       = "--client-id"
	FlagClientSecret   = "--client-secret"
	FlagSubscriptionId = "--subscription_id"
	FlagTenantId       = "--tenant-id"
	FlagResourceGroup  = "--resource-group"

	FlagDomainId          = "--domain-id"
	FlagDomainHostname    = "--domain-hostname"
	FlagDomainPort        = "--domain-port"
	FlagDomainUseSsl      = "--domain-use-ssl"
	FlagDomainScanDcCidr  = "--domain-scan-dc-cidr"
	FlagDomainNetworkCidr = "--domain-network-cidr"
	FlagDomainAdmin       = "--domain-admin"

	FlagGcpId             = "--gcp-id"
	FlagServiceAccountKey = "--service-account-key"
	FlagGoogleAdminEmail  = "--google-admin-email"
	FlagGcpRegion         = "--gcp-region"
)

// CLI literals for on/off permission flags and domain booleans.
const (
	ValueOn    = "on"
	ValueOff   = "off"
	ValueTrue  = "true"
	ValueFalse = "false"
)

// API response keys when create returns JSON in `data` (try in order in helper).
const (
	KeyPamConfigurationUID = "pam_configuration_uid"
	KeyPamConfigUID        = "pam_config_uid"
	KeyConfigurationUID    = "configuration_uid"
	KeyUID                 = "uid"
	KeyId                  = "id"
)

const (
	ErrSummaryCreateFailed   = "Create PAM Configuration Failed"
	ErrSummaryReadFailed     = "Read PAM Configuration Failed"
	ErrSummaryUpdateFailed   = "Update PAM Configuration Failed"
	ErrSummaryDeleteFailed   = "Delete PAM Configuration Failed"
	ErrSummarySyncDownFailed = "Sync Down Failed"
	ErrSummaryInvalidConfig  = "Invalid PAM Configuration"
	ErrSummaryNotImplemented = "PAM Configuration Not Implemented"
)

const (
	ErrOpCreatePamConfig = "Unable to create PAM configuration"
	ErrOpGetPamConfig    = "Unable to get PAM configuration"
	ErrOpListPamConfig   = "Unable to list PAM configuration"
	ErrOpEditPamConfig   = "Unable to update PAM configuration"
	ErrOpDeletePamConfig = "Unable to delete PAM configuration"
)

const (
	EnvLocal  = "local"
	EnvAWS    = "aws"
	EnvAzure  = "azure"
	EnvGCP    = "gcp"
	EnvDomain = "domain"
)

// ValidEnvironments lists allowed `environment` values (CLI: --environment).
var ValidEnvironments = []string{EnvLocal, EnvAWS, EnvAzure, EnvGCP, EnvDomain}

const (
	DescResource = "Manages a Keeper PAM configurations\n\n" +
		"In Keeper, the PAM Configuration contains essential information of your target infrastructure, settings and associated Keeper Gateway. We recommend setting up one PAM Configuration for each Gateway and network being managed.\n\n" +
		"For more information, see https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-configuration."
	DescResourceMD = "Manages a Keeper PAM configurations\n\n" +
		"In Keeper, the PAM Configuration contains essential information of your target infrastructure, settings and associated Keeper Gateway. We recommend setting up one PAM Configuration for each Gateway and network being managed.\n\n" +
		"For more information, see [Keeper PAM Configuration documentation](https://docs.keeper.io/en/keeperpam/privileged-access-manager/getting-started/pam-configuration)."

	DescId   = "The PAM configuration UID assigned by Keeper after create."
	DescIdMD = "The PAM configuration UID assigned by Keeper after create."

	DescEnvironment   = "PAM configuration type. One of: local, aws, azure, gcp, domain."
	DescEnvironmentMD = "PAM configuration type."

	DescTitle   = "Title of the PAM configuration."
	DescTitleMD = "Title of the PAM configuration."

	DescGateway   = "The configured gateway UID or name."
	DescGatewayMD = "The configured gateway `UID` or `name`."

	DescApplicationFolder   = "The shared folder name or UID where the PAM Configuration data will be stored"
	DescApplicationFolderMD = "The shared folder `name` or `UID` where the PAM Configuration data will be stored"

	DescSchedule   = "Specify frequency of Rotation Schedule using CRON syntax."
	DescScheduleMD = "Specify frequency of Rotation Schedule using `CRON` syntax."

	DescPortMapping   = "Define alternative default ports. Multiple values allowed."
	DescPortMappingMD = "Define `alternative default ports`. Multiple values allowed."

	DescConnections   = "If enabled, allow connections on resources managed by this PAM configuration ."
	DescConnectionsMD = "If `enabled`, allow connections on resources managed by this PAM configuration ."

	DescTunneling   = "If enabled, allow tunnels on resources managed by this PAM configuration."
	DescTunnelingMD = "If `enabled`, allow tunnels on resources managed by this PAM configuration"

	DescRotation   = "If enabled, allow rotations on privileged user users managed by this PAM configuration."
	DescRotationMD = "If `enabled`, allow rotations on privileged user users managed by this PAM configuration."

	DescRemoteBrowserIsolation   = "If enabled, allow RBI sessions on resources managed by this PAM configuration."
	DescRemoteBrowserIsolationMD = "If `enabled`, allow RBI sessions on resources managed by this PAM configuration."

	DescConnectionsRecording   = "If enabled, visual playback sessions will be recorded for all connections and RBI sessions"
	DescConnectionsRecordingMD = "If `enabled`, visual playback sessions will be recorded for all connections and RBI sessions"

	DescTypescriptRecording   = "If enabled, text input and output logs will be logged for all connections and RBI sessions "
	DescTypescriptRecordingMD = "If `enabled`, text input and output logs will be logged for all connections and RBI sessions "

	DescAIThreatDetection   = "If enabled, AI threat detection will be enabled "
	DescAIThreatDetectionMD = "If `enabled`, AI threat detection will be enabled."

	DescAITerminateSessionOnDetection   = "If enabled, AI session termination on threat detection will be enabled."
	DescAITerminateSessionOnDetectionMD = "If `enabled`, AI session termination on threat detection will be enabled."

	DescLocalNetwork   = "Use this block if you are using the local as environment type for your PAM configuration."
	DescLocalNetworkMD = "Use this block if you are using the `local` as environment type for your PAM configuration."

	DescAws   = "Use this block if you are using the aws as environment type for your PAM configuration."
	DescAwsMD = "Use this block if you are using the `aws` as environment type for your PAM configuration."

	DescAzure   = "Use this block if you are using the azure as environment type for your PAM configuration."
	DescAzureMD = "Use this block if you are using the `azure` as environment type for your PAM configuration."

	DescDomain   = "Use this block if you are using the `domain` as environment type for your PAM configuration."
	DescDomainMD = "Use this block if you are using the `domain` as environment type for your PAM configuration."

	DescGcp   = "Use this block if you are using the `gcp` as environment type for your PAM configuration."
	DescGcpMD = "Use this block if you are using the `gcp` as environment type for your PAM configuration."

	/** Local Network */
	DescNetworkId   = "Unique ID for the network. This is for the user's reference, Ex: My Network"
	DescNetworkIdMD = "`Unique ID` for the network. This is for the user's reference, Ex: `My Network`"

	DescNetworkCidr   = "Subnet of the IP address. Ex: 192.168.0.15/24"
	DescNetworkCidrMD = "`Subnet` of the IP address, Ex: `192.168.0.15/24`. Refer to [this](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing) for more info"

	/** AWS */

	DescAwsId   = "A unique id for the instance of AWS."
	DescAwsIdMD = "A `Unique ID` for the instance of AWS."

	DescAccessKeyId   = "From an IAM user account, the Access key ID from the desired Access key."
	DescAccessKeyIdMD = "From an IAM user account, the `Access key ID` from the desired Access key."

	DescAccessSecretKey   = "The secret key for the access key."
	DescAccessSecretKeyMD = "The `secret key` for the access key."

	DescRegionNames   = "AWS region names used for discovery. Separate newline per region. Ex: us-east-2"
	DescRegionNamesMD = "AWS region names used for discovery. Separate newline per region. Ex: `us-east-2`"

	/** Azure */

	DescAzureId   = "A unique id for your instance of Azure."
	DescAzureIdMD = "A `Unique ID` for your instance of Azure."

	DescClientId   = "The application/client id (UUID) of the Azure application"
	DescClientIdMD = "The `application/client id` (UUID) of the Azure application"

	DescClientSecret   = "The client credentials secret for the Azure application"
	DescClientSecretMD = "The `client credentials secret` for the Azure application"

	DescSubscriptionId   = "The UUID of the subscription (i.e. Pay-As-You-GO)."
	DescSubscriptionIdMD = "The `UUID` of the subscription (i.e. Pay-As-You-GO)."

	DescTenantId   = "The UUID of the Azure Active Directory"
	DescTenantIdMD = "The `UUID` of the Azure Active Directory"

	DescResourceGroups   = "A list of resource groups to be checked. If left blank, all resource groups will be checked. "
	DescResourceGroupsMD = "A list of `resource groups` to be checked. If left blank, all resource groups will be checked. "

	/** Domain */

	DescDomainId   = "The FQDN domain used by the Domain Controller. For example, EXAMPLE.COM and not EXAMPLE."
	DescDomainIdMD = "The `FQDN` domain used by the Domain Controller. For example, `EXAMPLE.COM` and not `EXAMPLE`."

	DescDomainHostname   = "Hostname for the domain controller."
	DescDomainHostnameMD = "Hostname for the `domain controller`."

	DescDomainPort   = "Port for the domain controller."
	DescDomainPortMD = "Port for the `domain controller`."

	DescDomainUseSsl   = "Provide true if using LDAPS (default 636), Provide false if using LDAP (default 389)."
	DescDomainUseSslMD = "Provide `true` if using `LDAPS` (default 636), Provide `false` if using `LDAP` (default 389)."

	DescDomainScanDcCidr   = "Scan the CIDRs from the domain controller. Default to False"
	DescDomainScanDcCidrMD = "Scan the CIDRs from the domain controller. Default to `False`"

	DescDomainNetworkCidr   = "Scan additional CIDRs from the field."
	DescDomainNetworkCidrMD = "`Scan additional CIDRs` from the field."

	DescDomainAdmin   = "Credentials of a domain administrator or an account with equivalent privileges, required to perform full discovery and access all domain resources."
	DescDomainAdminMD = "Credentials of a `domain administrator` or an account with equivalent privileges, required to perform full discovery and access all domain resources."

	/** GCP */

	DescGcpId   = "A unique id for the instance of Google Cloud. This is for the user's reference. Example: GCP-US-CENTRAL1"
	DescGcpIdMD = "A `unique id` for the instance of Google Cloud. This is for the user's reference. Example: `GCP-US-CENTRAL1`"

	DescServiceAccountKey   = "The service account key in JSON format."
	DescServiceAccountKeyMD = "The *service account key* in `JSON` format."

	DescGoogleAdminEmail   = "The email address for a Google Workspace administrator account that can be used to manage passwords for GCP Principals. Omit if no such account exists, or if the environment does not require Principal rotation."
	DescGoogleAdminEmailMD = "The `email address` for a Google Workspace administrator account that can be *used to manage passwords for GCP Principals*. Omit if **no such account exists**, or if the **environment does not require Principal rotation**."

	DescGcpRegion   = "GCP region names used for discovery."
	DescGcpRegionMD = "GCP region names used for discovery."
)

// ValidEnvironmentsMarkdown formats environments for schema markdown.
func ValidEnvironmentsMarkdown() string {
	b := strings.Builder{}
	for i, e := range ValidEnvironments {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("`")
		b.WriteString(e)
		b.WriteString("`")
	}
	return b.String()
}
