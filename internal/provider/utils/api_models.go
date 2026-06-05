// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NodeInfo represents a node from the enterprise-info API response.
type EnterpriseNodeResponse struct {
	NodeId         int    `json:"node_id"`
	Name           string `json:"name"`
	ParentNodeName string `json:"parent_node"`
	ParentNodeId   int    `json:"parent_id"`
	Isolated       bool   `json:"isolated"`
}

type ManagedCompanyResponse struct {
	CompanyId   int      `json:"company_id"`
	CompanyName string   `json:"company_name"`
	Node        string   `json:"node"`
	NodeName    string   `json:"node_name"`
	Plan        string   `json:"plan"`
	Storage     string   `json:"storage"`
	Addons      []string `json:"addons"`
	Allocated   int      `json:"allocated"`
}

// EnterpriseTeamResponse represents the team information from the read API response.
type EnterpriseTeamResponse struct {
	TeamUid   string   `json:"team_uid"`
	Name      string   `json:"name"`
	Restricts string   `json:"restricts"`
	Node      string   `json:"node"`
	Users     []string `json:"users"`
	Roles     []string `json:"roles"`
}

// ManagedNodePermission represents one entry in the managed_nodes_permissions array from the API.
type ManagedNodePermission struct {
	NodeName   string   `json:"node_name"`
	NodeId     int64    `json:"node_id"`
	Cascade    bool     `json:"cascade"`
	Privileges []string `json:"privileges"`
}

// EnterpriseRoleResponse represents a role from the enterprise-info API response.
type EnterpriseRoleResponse struct {
	RoleId                  int                     `json:"role_id"`
	Name                    string                  `json:"name"`
	Node                    string                  `json:"node"`
	VisibleBelow            bool                    `json:"visible_below"`
	DefaultRole             bool                    `json:"default_role"`
	Admin                   bool                    `json:"admin"`
	Users                   []string                `json:"users"`
	Teams                   []string                `json:"teams"`
	Enforcements            map[string]string       `json:"enforcements"`
	ManagedNodesPermissions []ManagedNodePermission `json:"managed_nodes_permissions"`
}

// EnterpriseUserResponse represents a user from the API response.
type EnterpriseUserResponse struct {
	UserId   int      `json:"user_id"`
	Email    string   `json:"email"`
	Status   string   `json:"status"`
	Name     string   `json:"name"`
	JobTitle string   `json:"job_title"`
	Roles    []string `json:"roles"`
	Teams    []string `json:"teams"`
	Node     string   `json:"node"`
}

// EnterpriseScimResponse represents the SCIM configuration from the read API response.
type EnterpriseScimResponse struct {
	ScimID            int    `json:"scim_id"`
	ScimURL           string `json:"scim_url"`
	NodeID            int    `json:"node_id"`
	NodeName          string `json:"node_name"`
	Status            string `json:"status"`
	Prefix            string `json:"prefix"`
	UniqueGroups      bool   `json:"unique_groups"`
	ProvisioningToken string `json:"provisioning_token"`
}

// EpmPolicyCreateResponse represents the response from the EPM policy add command.
type EpmPolicyCreateResponse struct {
	PolicyID string `json:"policy_id"`
}

// EpmPolicyResponse is the JSON from `epm policy view <id> --format json` (apiResp.Data).
type EpmPolicyResponse struct {
	PolicyName                   string              `json:"PolicyName"`
	PolicyType                   string              `json:"PolicyType"`
	PolicyId                     string              `json:"PolicyId"`
	Status                       string              `json:"Status"`
	Actions                      *EpmPolicyActions   `json:"Actions"`
	UserCheck                    []string            `json:"UserCheck"`
	MachineCheck                 []string            `json:"MachineCheck"`
	ApplicationCheck             []string            `json:"ApplicationCheck"`
	DayCheck                     []int               `json:"DayCheck"`
	DateCheck                    []EpmPolicyDateSpan `json:"DateCheck"`
	TimeCheck                    []EpmPolicyTimeSpan `json:"TimeCheck"`
	Message                      string              `json:"NotificationMessage"`
	RequirePolicyAcknowledgement bool                `json:"NotificationRequiresAcknowledge"`
}

// EpmPolicyActions mirrors API "Actions" (controls live under OnSuccess).
type EpmPolicyActions struct {
	OnSuccess *EpmPolicyOnSuccess `json:"OnSuccess"`
}

// EpmPolicyOnSuccess holds success-path controls from the API.
type EpmPolicyOnSuccess struct {
	Controls []string `json:"Controls"`
}

// EpmPolicyTimeSpan is one API time range (24h clock strings).
type EpmPolicyTimeSpan struct {
	StartTime string `json:"StartTime"`
	EndTime   string `json:"EndTime"`
}

// EpmPolicyDateSpan holds normalized YYYY-MM-DD start/end after JSON unmarshal.
type EpmPolicyDateSpan struct {
	StartDate string
	EndDate   string
}

const epmPolicyDateLayout = "2006-01-02"

// UnmarshalJSON accepts StartDate/EndDate as ISO strings or numeric epoch milliseconds.
func (d *EpmPolicyDateSpan) UnmarshalJSON(b []byte) error {
	var aux struct {
		StartDate json.RawMessage `json:"StartDate"`
		EndDate   json.RawMessage `json:"EndDate"`
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	var err error
	d.StartDate, err = normalizeEpmPolicyDateField(aux.StartDate)
	if err != nil {
		return fmt.Errorf("StartDate: %w", err)
	}
	d.EndDate, err = normalizeEpmPolicyDateField(aux.EndDate)
	if err != nil {
		return fmt.Errorf("EndDate: %w", err)
	}
	return nil
}

func normalizeEpmPolicyDateField(raw json.RawMessage) (string, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || string(raw) == "null" {
		return "", nil
	}
	// ISO date string
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return "", err
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return "", nil
		}
		// Epoch milliseconds sent as a JSON string
		if isDecimalDigits(s) {
			if ms, err := strconv.ParseInt(s, 10, 64); err == nil {
				return time.UnixMilli(ms).UTC().Format(epmPolicyDateLayout), nil
			}
		}
		t, err := time.Parse(epmPolicyDateLayout, s)
		if err != nil {
			return "", fmt.Errorf("invalid date string %q: %w", s, err)
		}
		return t.Format(epmPolicyDateLayout), nil
	}
	// JSON number: epoch milliseconds (e.g. 1777228200000). Prefer json.Number for Int64.
	var num json.Number
	if err := json.Unmarshal(raw, &num); err == nil && num != "" {
		if ms, err := num.Int64(); err == nil {
			return time.UnixMilli(ms).UTC().Format(epmPolicyDateLayout), nil
		}
		f, err := num.Float64()
		if err != nil {
			return "", fmt.Errorf("invalid numeric date field: %w", err)
		}
		return time.UnixMilli(int64(f)).UTC().Format(epmPolicyDateLayout), nil
	}
	var n float64
	if err := json.Unmarshal(raw, &n); err != nil {
		return "", err
	}
	return time.UnixMilli(int64(n)).UTC().Format(epmPolicyDateLayout), nil
}

func isDecimalDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// PamConfigListSharedFolder is the shared_folder object from pam config list --format json.
type PamConfigListSharedFolder struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

// PamConfigListResponse is the data payload from pam config list --config ID --format json.
type PamConfigListResponse struct {
	UID                            string                            `json:"uid"`
	Name                           string                            `json:"name"`
	ConfigType                     string                            `json:"config_type"`
	SharedFolder                   PamConfigListSharedFolder         `json:"shared_folder"`
	GatewayUID                     string                            `json:"gateway_uid"`
	GatewayName                    string                            `json:"gateway_name"`
	AllowedSettings                *PamConfigAllowedSettingsResponse `json:"allowed_settings,omitempty"`
	Fields                         *PamConfigFieldsResponse          `json:"fields,omitempty"`
	DomainAdministrativeCredential string                            `json:"domain_administrative_credential,omitempty"`
}

// PamConfigFieldsResponse maps the "fields" object from pam config list response.
// Each key holds a string slice; only the keys relevant to the environment will be populated.
type PamConfigFieldsResponse struct {
	DefaultSchedule []string `json:"Default Schedule,omitempty"`
	PortMapping     []string `json:"portMapping,omitempty"`

	// Local Network
	NetworkId   []string `json:"networkId,omitempty"`
	NetworkCIDR []string `json:"networkCIDR,omitempty"`

	// AWS
	AwsId           []string `json:"awsId,omitempty"`
	AccessKeyId     []string `json:"accessKeyId,omitempty"`
	AccessSecretKey []string `json:"accessSecretKey,omitempty"`
	RegionNames     []string `json:"regionNames,omitempty"`

	// Azure
	AzureId        []string `json:"azureId,omitempty"`
	ClientId       []string `json:"clientId,omitempty"`
	ClientSecret   []string `json:"clientSecret,omitempty"`
	SubscriptionId []string `json:"subscriptionId,omitempty"`
	TenantId       []string `json:"tenantId,omitempty"`
	ResourceGroups []string `json:"resourceGroups,omitempty"`

	// Domain
	PamHostname []string `json:"pamHostname,omitempty"`
	PamDomainId []string `json:"pamDomainId,omitempty"`
	UseSSL      []string `json:"useSSL,omitempty"`
	ScanDCCIDR  []string `json:"scanDCCIDR,omitempty"`

	// GCP
	PamGcpId             []string `json:"pamGcpId,omitempty"`
	PamGoogleAdminEmail  []string `json:"pamGoogleAdminEmail,omitempty"`
	PamServiceAccountKey []string `json:"pamServiceAccountKey,omitempty"`
	PamGcpRegionName     []string `json:"pamGcpRegionName,omitempty"`
}

// PamConfigAllowedSettingsResponse maps the allowed_settings object from the API response.
type PamConfigAllowedSettingsResponse struct {
	Connections                   bool `json:"connections"`
	Tunneling                     bool `json:"tunneling"`
	Rotation                      bool `json:"rotation"`
	RemoteBrowserIsolation        bool `json:"remote_browser_isolation"`
	ConnectionsRecording          bool `json:"connections_recording"`
	TypescriptRecording           bool `json:"typescript_recording"`
	AIThreatDetection             bool `json:"ai_threat_detection"`
	AITerminateSessionOnDetection bool `json:"ai_terminate_session_on_detection"`
}

// VaultRecordGetResponse is the JSON payload from `get <record_uid> --format json` for a Keeper vault record.
type VaultRecordGetResponse struct {
	RecordUID                    string                                `json:"record_uid"`
	Type                         string                                `json:"type"`
	Title                        string                                `json:"title"`
	Notes                        string                                `json:"notes"`
	Fields                       []VaultRecordFieldResponse            `json:"fields"`
	Custom                       []VaultRecordFieldResponse            `json:"custom"`
	PamSettingsEnabled           *PamSettingsEnabledResponse           `json:"pamSettingsEnabled,omitempty"`
	DagDebug                     *DagDebugResponse                     `json:"dagDebug,omitempty"`
	AssociatedCredentials        *AssociatedCredentialsResponse        `json:"associatedCredentials,omitempty"`
	FolderLocation               *FolderLocationResponse               `json:"folder,omitempty"`
	PamConfigurationUID          string                                `json:"pam_configuration_uid,omitempty"`
	ConfigurationAllowedSettings *ConfigurationAllowedSettingsResponse `json:"configuration_allowed_settings,omitempty"`
	UserPermissions              []UserPermissionEntry                 `json:"user_permissions,omitempty"`
}

// UserPermissionEntry is one element of the API response's user_permissions
// array. The same JSON array key (`user_permissions`) is returned in two
// different shapes depending on the record style:
//
//   - NSF (nsf-get / nsf-share-*) entries carry {accessor, role}; consumed
//     by new_share.MapResponseToModel.
//   - Classic (get / share-record) entries carry {username, shareable,
//     editable}; consumed by classic_share.MapResponseToModel.
//
// All fields use `omitempty` so each helper sees zero values for the
// irrelevant shape and naturally filters them out.
// new_share.UserPermissionEntry and classic_share.UserPermissionEntry both
// alias this type.
type UserPermissionEntry struct {
	Accessor  string `json:"accessor,omitempty"`
	Role      string `json:"role,omitempty"`
	Username  string `json:"username,omitempty"`
	Shareable bool   `json:"shareable,omitempty"`
	Editable  bool   `json:"editable,omitempty"`
}

// ConfigurationAllowedSettingsResponse maps the configuration_allowed_settings object from the API response.
type ConfigurationAllowedSettingsResponse struct {
	ConnectionsRecording   bool `json:"connections_recording"`
	RemoteBrowserIsolation bool `json:"remote_browser_isolation"`
}

/*
FolderLocationResponse is struct of folder field that is returned by the API response for get <record/folder> --format json.
Which consists of uid and path where record/folder is located.

To use this struct in type of api response of folder and record response like:

folder type FolderLocationResponse  'json:"folder"'.
*/
type FolderLocationResponse struct {
	UID  string `json:"uid"`
	Path string `json:"path"`
}

type PamSettingsEnabledResponse struct {
	Connections            *bool `json:"connections"`
	Tunneling              *bool `json:"tunneling"`
	Rotation               *bool `json:"rotation"`
	SessionRecording       *bool `json:"sessionRecording"`
	TypescriptRecording    *bool `json:"typescriptRecording"`
	RemoteBrowserIsolation *bool `json:"remoteBrowserIsolation"`
}

type DagDebugResponse struct {
	VertexContent *DagDebugVertexContentResponse `json:"vertex_content,omitempty"`
	AllEdges      []DagDebugEdgeResponse         `json:"all_edges,omitempty"`
}

type DagDebugEdgeResponse struct {
	Type    string `json:"type"`
	HeadUID string `json:"head_uid"`
}

type DagDebugVertexContentResponse struct {
	AllowedSettings     *DagDebugAllowedSettingsResponse `json:"allowedSettings,omitempty"`
	RotateOnTermination bool                             `json:"rotateOnTermination"`
}

type DagDebugAllowedSettingsResponse struct {
	PortForwards        bool `json:"portForwards"`
	Rotation            bool `json:"rotation"`
	Connections         bool `json:"connections"`
	SessionRecording    bool `json:"sessionRecording"`
	TypescriptRecording bool `json:"typescriptRecording"`
	AiEnabled           bool `json:"aiEnabled"`
	AiSessionTerminate  bool `json:"aiSessionTerminate"`
}

type AssociatedCredentialsResponse struct {
	AdminCredential  *string `json:"adminCredential"`
	LaunchCredential *string `json:"launchCredential"`
}

// VaultRecordFieldResponse is one typed field entry inside a vault record (get --format json).
type VaultRecordFieldResponse struct {
	Type     string          `json:"type"`
	Label    string          `json:"label"`
	Value    json.RawMessage `json:"value"`
	Required bool            `json:"required"`
}

// PamSettingsFieldValueResponse is one element of the pamSettings field value array from `get <uid> --format json`.
// The fields array contains an entry with type "pamSettings" whose value is a JSON array;
// the first element carries the settings object.
type PamSettingsFieldValueResponse struct {
	AllowSupplyHost bool                            `json:"allowSupplyHost"`
	PortForward     *PamSettingsPortForwardResponse `json:"portForward,omitempty"`
	Connection      json.RawMessage                 `json:"connection,omitempty"`
}

// PamSettingsConnectionBaseResponse is used for a first-pass unmarshal to
// peek at the "protocol" field before unmarshaling into the correct
// per-protocol struct.
type PamSettingsConnectionBaseResponse struct {
	Protocol string `json:"protocol"`
	Port     string `json:"port"`
}

// KubernetesConnectionResponse holds the kubernetes-specific fields
// returned inside pamSettings.value[0].connection when protocol is "kubernetes".
type KubernetesConnectionResponse struct {
	Protocol             string `json:"protocol"`
	Port                 string `json:"port"`
	RecordingIncludeKeys *bool  `json:"recordingIncludeKeys"`
	AllowSupplyUser      *bool  `json:"allowSupplyUser"`
	UseSSL               *bool  `json:"useSSL"`
	IgnoreCert           *bool  `json:"ignoreCert"`
	CaCert               string `json:"caCert"`
	ClientCert           string `json:"clientCert"`
	ClientKey            string `json:"clientKey"`
	Namespace            string `json:"namespace"`
	Pod                  string `json:"pod"`
	Container            string `json:"container"`
	Command              string `json:"command"`
	ColorScheme          string `json:"colorScheme"`
	FontName             string `json:"fontName"`
	FontSize             string `json:"fontSize"`
	Scrollback           int    `json:"scrollback"`
	ReadOnly             *bool  `json:"readOnly"`
	Backspace            string `json:"backspace"`
}

// RdpConnectionResponse holds the fields returned inside
// pamSettings.value[0].connection when protocol is "rdp".
type RdpConnectionResponse struct {
	Protocol                 string        `json:"protocol"`
	Port                     string        `json:"port"`
	AllowSupplyUser          *bool         `json:"allowSupplyUser"`
	RecordingIncludeKeys     *bool         `json:"recordingIncludeKeys"`
	EnableFullWindowDrag     *bool         `json:"enableFullWindowDrag"`
	EnableWallpaper          *bool         `json:"enableWallpaper"`
	IgnoreCert               *bool         `json:"ignoreCert"`
	NormalizeClipboard       string        `json:"normalizeClipboard"`
	Security                 string        `json:"security"`
	EnableTheming            *bool         `json:"enableTheming"`
	EnableFontSmoothing      *bool         `json:"enableFontSmoothing"`
	EnableDesktopComposition *bool         `json:"enableDesktopComposition"`
	EnableMenuAnimations     *bool         `json:"enableMenuAnimations"`
	DisableBitmapCaching     *bool         `json:"disableBitmapCaching"`
	DisableOffscreenCaching  *bool         `json:"disableOffscreenCaching"`
	DisableGlyphCaching      *bool         `json:"disableGlyphCaching"`
	LoadBalanceInfo          string        `json:"loadBalanceInfo"`
	PreconnectionId          string        `json:"preconnectionId"`
	PreconnectionBlob        string        `json:"preconnectionBlob"`
	Sftp                     *SftpResponse `json:"sftp"`
	ConsoleAudio             *bool         `json:"consoleAudio"`
	DisableAudio             *bool         `json:"disableAudio"`
	EnableAudioInput         *bool         `json:"enableAudioInput"`
	EnablePrinting           *bool         `json:"enablePrinting"`
	RedirectedPrinterName    string        `json:"redirectedPrinterName"`
	RemoteApp                string        `json:"remoteApp"`
	RemoteAppDir             string        `json:"remoteAppDir"`
	RemoteAppArgs            string        `json:"remoteAppArgs"`
	ForceLossless            *bool         `json:"forceLossless"`
	ReadOnly                 *bool         `json:"readOnly"`
	Dpi                      int           `json:"dpi"`
	Height                   int           `json:"height"`
	Width                    int           `json:"width"`
	EnableTouch              *bool         `json:"enableTouch"`
	Console                  *bool         `json:"console"`
	Timezone                 string        `json:"timezone"`
	ClientName               string        `json:"clientName"`
	InitialProgram           string        `json:"initialProgram"`
	DisableAuth              *bool         `json:"disableAuth"`
	ResizeMethod             string        `json:"resizeMethod"`
	ColorDepth               int           `json:"colorDepth"`
	ServerLayout             string        `json:"serverLayout"`
	DisableCopy              *bool         `json:"disableCopy"`
	DisablePaste             *bool         `json:"disablePaste"`
}

// SftpResponse is the shared SFTP nested block used by RDP and VNC.
type SftpResponse struct {
	EnableSftp              *bool  `json:"enableSftp"`
	SftpResourceUid         string `json:"sftpResourceUid"`
	SftpUserUid             string `json:"sftpUserUid"`
	SftpDirectory           string `json:"sftpDirectory"`
	SftpServerAliveInterval int    `json:"sftpServerAliveInterval"`
}

// SshConnectionResponse holds the fields returned inside
// pamSettings.value[0].connection when protocol is "ssh".
type SshConnectionResponse struct {
	Protocol             string           `json:"protocol"`
	Port                 string           `json:"port"`
	AllowSupplyUser      *bool            `json:"allowSupplyUser"`
	RecordingIncludeKeys *bool            `json:"recordingIncludeKeys"`
	ReadOnly             *bool            `json:"readOnly"`
	DisableCopy          *bool            `json:"disableCopy"`
	DisablePaste         *bool            `json:"disablePaste"`
	ColorScheme          string           `json:"colorScheme"`
	FontName             string           `json:"fontName"`
	FontSize             string           `json:"fontSize"`
	Scrollback           int              `json:"scrollback"`
	HostKey              string           `json:"hostKey"`
	Command              string           `json:"command"`
	Locale               string           `json:"locale"`
	Timezone             string           `json:"timezone"`
	ServerAliveInterval  int              `json:"serverAliveInterval"`
	Backspace            string           `json:"backspace"`
	TerminalType         string           `json:"terminalType"`
	Sftp                 *SshSftpResponse `json:"sftp"`
}

type SshSftpResponse struct {
	EnableSftp *bool `json:"enableSftp"`
}

// TelnetConnectionResponse holds the fields returned inside
// pamSettings.value[0].connection when protocol is "telnet".
type TelnetConnectionResponse struct {
	Protocol             string `json:"protocol"`
	Port                 string `json:"port"`
	AllowSupplyUser      *bool  `json:"allowSupplyUser"`
	RecordingIncludeKeys *bool  `json:"recordingIncludeKeys"`
	ReadOnly             *bool  `json:"readOnly"`
	DisableCopy          *bool  `json:"disableCopy"`
	DisablePaste         *bool  `json:"disablePaste"`
	ColorScheme          string `json:"colorScheme"`
	FontName             string `json:"fontName"`
	FontSize             string `json:"fontSize"`
	Scrollback           int    `json:"scrollback"`
	UsernameRegex        string `json:"usernameRegex"`
	PasswordRegex        string `json:"passwordRegex"`
	LoginSuccessRegex    string `json:"loginSuccessRegex"`
	LoginFailureRegex    string `json:"loginFailureRegex"`
	Backspace            string `json:"backspace"`
	TerminalType         string `json:"terminalType"`
}

// VncConnectionResponse holds the fields returned inside
// pamSettings.value[0].connection when protocol is "vnc".
type VncConnectionResponse struct {
	Protocol             string        `json:"protocol"`
	Port                 string        `json:"port"`
	AllowSupplyUser      *bool         `json:"allowSupplyUser"`
	RecordingIncludeKeys *bool         `json:"recordingIncludeKeys"`
	ReadOnly             *bool         `json:"readOnly"`
	SwapRedBlue          *bool         `json:"swapRedBlue"`
	ForceLossless        *bool         `json:"forceLossless"`
	EnableAudio          *bool         `json:"enableAudio"`
	AudioServername      string        `json:"audioServername"`
	DestHost             string        `json:"destHost"`
	DestPort             string        `json:"destPort"`
	DisableCopy          *bool         `json:"disableCopy"`
	DisablePaste         *bool         `json:"disablePaste"`
	ClipboardEncoding    string        `json:"clipboardEncoding"`
	Cursor               string        `json:"cursor"`
	ColorDepth           string        `json:"colorDepth"`
	Sftp                 *SftpResponse `json:"sftp"`
}

// DatabaseConnectionResponse holds the fields returned inside
// pamSettings.value[0].connection for mysql, postgresql, and sql-server protocols.
type DatabaseConnectionResponse struct {
	Protocol             string `json:"protocol"`
	Port                 string `json:"port"`
	AllowSupplyUser      *bool  `json:"allowSupplyUser"`
	RecordingIncludeKeys *bool  `json:"recordingIncludeKeys"`
	ReadOnly             *bool  `json:"readOnly"`
	DisableCopy          *bool  `json:"disableCopy"`
	DisablePaste         *bool  `json:"disablePaste"`
	DisableCsvExport     *bool  `json:"disableCsvExport"`
	DisableCsvImport     *bool  `json:"disableCsvImport"`
	Database             string `json:"database"`
	ColorScheme          string `json:"colorScheme"`
	FontName             string `json:"fontName"`
	FontSize             string `json:"fontSize"`
	Scrollback           int    `json:"scrollback"`
}

type PamSettingsPortForwardResponse struct {
	Port                  string `json:"port"`
	ReusePort             *bool  `json:"reusePort"`
	UseSpecifiedLocalPort *bool  `json:"useSpecifiedLocalPort"`
	LocalPort             string `json:"localPort"`
}

// PamHostnameFieldValue is one element of the pamHostname field value array from `get <uid> --format json`.
type PamRemoteBrowserHostnameFieldResponse struct {
	HostName           string `json:"hostName"`
	AdministrativePort string `json:"port"`
}

// PamRemoteBrowserSettingsFieldConnectionResponse is the API `connection` object inside pamRemoteBrowserSettings.
// Note: currently we are not getting --connections-recording, --remote-browser-isolation, --configuration data from the API.
type PamRemoteBrowserSettingsFieldConnectionResponse struct {
	ConfigurationUID           string `json:"configurationUid,omitempty"`
	HttpCredentialsUID         string `json:"httpCredentialsUid,omitempty"`
	AutofillConfiguration      string `json:"autofillConfiguration,omitempty"`
	RecordingIncludeKeys       bool   `json:"recordingIncludeKeys"`
	RecordingScreens           bool   `json:"recordingScreens"`
	RemoteBrowserIsolation     bool   `json:"remoteBrowserIsolation"`
	AllowUrlManipulation       bool   `json:"allowUrlManipulation"`
	IgnoreInitialSslCert       bool   `json:"ignoreInitialSslCert"`
	AllowedURLPatterns         string `json:"allowedUrlPatterns"`
	AllowedResourceURLPatterns string `json:"allowedResourceUrlPatterns"`
	DisableCopy                bool   `json:"disableCopy"`
	DisablePaste               bool   `json:"disablePaste"`
	DisableAudio               bool   `json:"disableAudio"`
	AudioChannels              int    `json:"audioChannels"`
	AudioBps                   int    `json:"audioBps"`
	AudioSampleRate            int    `json:"audioSampleRate"`
}

// PamRemoteBrowserSettingsFieldResponse is one element of the pamRemoteBrowserSettings field value array.
type PamRemoteBrowserSettingsFieldResponse struct {
	Connection PamRemoteBrowserSettingsFieldConnectionResponse `json:"connection"`
}

// VaultRecordField is one typed field entry inside a vault record (get --format json).
type VaultRecordField struct {
	Type     string          `json:"type"`
	Label    string          `json:"label"`
	Value    json.RawMessage `json:"value"`
	Required bool            `json:"required"`
}

// PamRemoteBrowserSettingsFieldConnection is the API `connection` object inside pamRemoteBrowserSettings.
// Note: currently we are not getting --connections-recording, --remote-browser-isolation, --configuration data from the API.
type PamRemoteBrowserSettingsFieldConnection struct {
	ConfigurationUID           string `json:"configurationUid,omitempty"`
	HttpCredentialsUID         string `json:"httpCredentialsUid,omitempty"`
	AutofillConfiguration      string `json:"autofillConfiguration,omitempty"`
	RecordingIncludeKeys       bool   `json:"recordingIncludeKeys"`
	RecordingScreens           bool   `json:"recordingScreens"`
	RemoteBrowserIsolation     bool   `json:"remoteBrowserIsolation"`
	AllowUrlManipulation       bool   `json:"allowUrlManipulation"`
	IgnoreInitialSslCert       bool   `json:"ignoreInitialSslCert"`
	AllowedURLPatterns         string `json:"allowedUrlPatterns"`
	AllowedResourceURLPatterns string `json:"allowedResourceUrlPatterns"`
	DisableCopy                bool   `json:"disableCopy"`
	DisablePaste               bool   `json:"disablePaste"`
	DisableAudio               bool   `json:"disableAudio"`
	AudioChannels              int    `json:"audioChannels"`
	AudioBps                   int    `json:"audioBps"`
	AudioSampleRate            int    `json:"audioSampleRate"`
}

// PamRemoteBrowserSettingsFieldEntry is one element of the pamRemoteBrowserSettings field value array.
type PamRemoteBrowserSettingsFieldEntry struct {
	Connection PamRemoteBrowserSettingsFieldConnection `json:"connection"`
}
