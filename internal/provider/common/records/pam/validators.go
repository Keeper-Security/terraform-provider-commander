// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamrecords

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MachineDirectoryProtocols lists the connection protocols allowed on PAM
// Machine and PAM Directory records.
var MachineDirectoryProtocols = []string{
	ConnectionProtocolKubernetes,
	ConnectionProtocolRdp,
	ConnectionProtocolSsh,
	ConnectionProtocolTelnet,
	ConnectionProtocolVnc,
}

// DatabaseProtocols lists the connection protocols allowed on PAM Database
// records.
var DatabaseProtocols = []string{
	ConnectionProtocolMysql,
	ConnectionProtocolPostgreSql,
	ConnectionProtocolSqlServer,
	ConnectionProtocolMariaDb,
	ConnectionProtocolOracle,
}

// protocolToAttributeKey maps each protocol constant to the tfsdk attribute
// name used in CommonPamSettingsConnectionResourceModel.
var protocolToAttributeKey = map[string]string{
	ConnectionProtocolKubernetes: "kubernetes",
	ConnectionProtocolMysql:      "mysql",
	ConnectionProtocolPostgreSql: "postgresql",
	ConnectionProtocolRdp:        "rdp",
	ConnectionProtocolSqlServer:  "sql_server",
	ConnectionProtocolSsh:        "ssh",
	ConnectionProtocolTelnet:     "telnet",
	ConnectionProtocolVnc:        "vnc",
	ConnectionProtocolMariaDb:    "mariadb",
	ConnectionProtocolOracle:     "oracle",
}

// ---------------------------------------------------------------------------
// Block-level "required when present" validators
// ---------------------------------------------------------------------------

// pamSettingsRequiredFieldsValidator enforces that "configuration" is provided
// when the pam_settings block is explicitly written in HCL.
type pamSettingsRequiredFieldsValidator struct{}

func PamSettingsRequiredFieldsValidator() pamSettingsRequiredFieldsValidator {
	return pamSettingsRequiredFieldsValidator{}
}

func (v pamSettingsRequiredFieldsValidator) Description(_ context.Context) string {
	return "configuration is required when pam_settings block is provided."
}

func (v pamSettingsRequiredFieldsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v pamSettingsRequiredFieldsValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attrs := req.ConfigValue.Attributes()

	configAttr, ok := attrs["configuration"]
	if !ok || utils.IsExplicitlyMissing(configAttr) {
		resp.Diagnostics.AddAttributeError(
			req.Path.AtName("configuration"),
			"Missing Required PAM Settings Attribute",
			"configuration is required when pam_settings block is provided.",
		)
	}
}

// connectionRequiredFieldsValidator enforces that "enable" and "protocol" are
// provided when the connection block is explicitly written in HCL.
type connectionRequiredFieldsValidator struct{}

func ConnectionRequiredFieldsValidator() connectionRequiredFieldsValidator {
	return connectionRequiredFieldsValidator{}
}

func (v connectionRequiredFieldsValidator) Description(_ context.Context) string {
	return "enable and protocol are required when connection block is provided."
}

func (v connectionRequiredFieldsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v connectionRequiredFieldsValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attrs := req.ConfigValue.Attributes()

	for _, field := range []string{"enable", "protocol"} {
		attr, ok := attrs[field]
		if !ok || attr.IsNull() || attr.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName(field),
				"Missing Required Connection Attribute",
				fmt.Sprintf("%s is required when connection block is provided.", field),
			)
		}
	}
}

// tunnelRequiredFieldsValidator enforces that "enable" is provided
// when the tunnel block is explicitly written in HCL.
type tunnelRequiredFieldsValidator struct{}

func TunnelRequiredFieldsValidator() tunnelRequiredFieldsValidator {
	return tunnelRequiredFieldsValidator{}
}

func (v tunnelRequiredFieldsValidator) Description(_ context.Context) string {
	return "enable is required when tunnel block is provided."
}

func (v tunnelRequiredFieldsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v tunnelRequiredFieldsValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attrs := req.ConfigValue.Attributes()

	enableAttr, ok := attrs["enable"]
	if !ok || enableAttr.IsNull() || enableAttr.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			req.Path.AtName("enable"),
			"Missing Required Tunnel Attribute",
			"enable is required when tunnel block is provided.",
		)
	}
}

// ---------------------------------------------------------------------------
// Protocol string validator
// ---------------------------------------------------------------------------

type connectionProtocolValidator struct {
	Allowed []string
}

// ConnectionProtocolValidator validates that connection.protocol is one of
// the protocols allowed for the calling PAM record type. Pass
// MachineDirectoryProtocols for pam_machine / pam_directory and
// DatabaseProtocols for pam_database.
func ConnectionProtocolValidator(allowed []string) connectionProtocolValidator {
	return connectionProtocolValidator{Allowed: allowed}
}

func (v connectionProtocolValidator) Description(_ context.Context) string {
	return "Protocol must be one of: " + strings.Join(v.Allowed, ", ") + "."
}

func (v connectionProtocolValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v connectionProtocolValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	val := req.ConfigValue.ValueString()
	for _, allowed := range v.Allowed {
		if val == allowed {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Connection Protocol",
		fmt.Sprintf("Protocol %q is not supported for this PAM record type. Must be one of: %s.", val, strings.Join(v.Allowed, ", ")),
	)
}

// ---------------------------------------------------------------------------
// Connection object validator – only the block matching "protocol" may be set
// ---------------------------------------------------------------------------

type connectionProtocolBlockValidator struct{}

func ConnectionProtocolBlockValidator() connectionProtocolBlockValidator {
	return connectionProtocolBlockValidator{}
}

func (v connectionProtocolBlockValidator) Description(_ context.Context) string {
	return "Only the protocol-specific block matching the selected protocol may be set."
}

func (v connectionProtocolBlockValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v connectionProtocolBlockValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attrs := req.ConfigValue.Attributes()

	protocolAttr, ok := attrs["protocol"]
	if !ok {
		return
	}
	protocolVal, ok := protocolAttr.(types.String)
	if !ok || protocolVal.IsNull() || protocolVal.IsUnknown() {
		return
	}

	selectedKey, known := protocolToAttributeKey[protocolVal.ValueString()]
	if !known {
		return
	}

	for _, attrKey := range protocolToAttributeKey {
		attr, exists := attrs[attrKey]
		if !exists {
			continue
		}
		if attrKey == selectedKey {
			continue
		}
		if !attr.IsNull() && !attr.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName(attrKey),
				"Invalid Connection Configuration",
				fmt.Sprintf("The %q block is not allowed when protocol is %q. Only the %q block may be set.", attrKey, protocolVal.ValueString(), selectedKey),
			)
		}
	}
}

// connectionFieldsRequireEnabledValidator requires connection_port and
// launch_credential when enable is true.
type connectionFieldsRequireEnabledValidator struct{}

func ConnectionFieldsRequireEnabledValidator() connectionFieldsRequireEnabledValidator {
	return connectionFieldsRequireEnabledValidator{}
}

func (v connectionFieldsRequireEnabledValidator) Description(_ context.Context) string {
	return "connection_port and launch_credential are required when enable is true."
}

func (v connectionFieldsRequireEnabledValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v connectionFieldsRequireEnabledValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attrs := req.ConfigValue.Attributes()

	enableAttr, ok := attrs["enable"]
	if !ok {
		return
	}
	enableVal, ok := enableAttr.(types.Bool)
	if !ok || enableVal.IsNull() || enableVal.IsUnknown() {
		return
	}

	if !enableVal.ValueBool() {
		return
	}

	requiredFields := []string{"connection_port", "launch_credential"}
	for _, field := range requiredFields {
		attr, exists := attrs[field]
		if !exists || utils.IsExplicitlyMissing(attr) {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName(field),
				"Missing Required Connection Attribute",
				fmt.Sprintf("%s is required when connection enable is true.", field),
			)
		}
	}
}

// tunnelFieldsRequireEnabledValidator rejects tunnel sub-fields
// (remote_target_port, re_use_port, use_specified_local_port, local_port)
// when "enabled" is false or null.
type tunnelFieldsRequireEnabledValidator struct{}

func TunnelFieldsRequireEnabledValidator() tunnelFieldsRequireEnabledValidator {
	return tunnelFieldsRequireEnabledValidator{}
}

func (v tunnelFieldsRequireEnabledValidator) Description(_ context.Context) string {
	return "Tunnel sub-fields (remote_target_port, re_use_port, use_specified_local_port, local_port) are only allowed when enabled is true."
}

func (v tunnelFieldsRequireEnabledValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v tunnelFieldsRequireEnabledValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attrs := req.ConfigValue.Attributes()

	enabledAttr, ok := attrs["enabled"]
	if !ok {
		return
	}
	enabledVal, ok := enabledAttr.(types.Bool)
	if !ok || enabledVal.IsNull() || enabledVal.IsUnknown() {
		return
	}

	if enabledVal.ValueBool() {
		return
	}

	conditionalFields := []string{
		"remote_target_port",
		"re_use_port",
		"use_specified_local_port",
		"local_port",
	}

	for _, field := range conditionalFields {
		attr, exists := attrs[field]
		if !exists {
			continue
		}
		if !attr.IsNull() && !attr.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtName(field),
				"Invalid Tunnel Configuration",
				field+" is only allowed when tunnel enabled is true.",
			)
		}
	}
}

// tunnelLocalPortRequiredValidator requires local_port when use_specified_local_port is true.
type tunnelLocalPortRequiredValidator struct{}

func TunnelLocalPortRequiredValidator() tunnelLocalPortRequiredValidator {
	return tunnelLocalPortRequiredValidator{}
}

func (v tunnelLocalPortRequiredValidator) Description(_ context.Context) string {
	return "local_port is required when use_specified_local_port is true."
}

func (v tunnelLocalPortRequiredValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v tunnelLocalPortRequiredValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attrs := req.ConfigValue.Attributes()

	useSpecifiedAttr, ok := attrs["use_specified_local_port"]
	if !ok {
		return
	}
	useSpecifiedVal, ok := useSpecifiedAttr.(types.Bool)
	if !ok || useSpecifiedVal.IsNull() || useSpecifiedVal.IsUnknown() || !useSpecifiedVal.ValueBool() {
		return
	}

	localPortAttr, exists := attrs["local_port"]
	if !exists || localPortAttr.IsNull() || localPortAttr.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			req.Path.AtName("local_port"),
			"Missing Required Tunnel Attribute",
			"local_port is required when use_specified_local_port is true.",
		)
	}
}

// sftpUserUidRequiredValidator requires sftp_user_uid when enable_sftp is true.
type sftpUserUidRequiredValidator struct{}

func SftpUserUidRequiredValidator() sftpUserUidRequiredValidator {
	return sftpUserUidRequiredValidator{}
}

func (v sftpUserUidRequiredValidator) Description(_ context.Context) string {
	return "sftp_user_uid is required when enable_sftp is true."
}

func (v sftpUserUidRequiredValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v sftpUserUidRequiredValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attrs := req.ConfigValue.Attributes()

	enableAttr, ok := attrs["enable_sftp"]
	if !ok {
		return
	}
	enableVal, ok := enableAttr.(types.Bool)
	if !ok || enableVal.IsNull() || enableVal.IsUnknown() {
		return
	}

	if !enableVal.ValueBool() {
		return
	}

	userUidAttr, exists := attrs["sftp_user_uid"]
	if !exists || userUidAttr.IsNull() || userUidAttr.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			req.Path.AtName("sftp_user_uid"),
			"Missing Required SFTP Attribute",
			"sftp_user_uid is required when enable_sftp is true.",
		)
	}
}

// ---------------------------------------------------------------------------
// Locale validator – accepts $LANG or valid POSIX locale strings
// ---------------------------------------------------------------------------

var posixLocalePattern = regexp.MustCompile(
	`^[a-z]{2,3}(_[A-Z]{2,3})?(\.[A-Za-z0-9_-]+)?(@[A-Za-z0-9]+)?$`,
)

type localeValidator struct{}

func LocaleValidator() localeValidator {
	return localeValidator{}
}

func (v localeValidator) Description(_ context.Context) string {
	return `Must be "$LANG" (use client's locale) or a valid POSIX locale string (e.g. "en_US.UTF-8", "fr_FR.UTF-8").`
}

func (v localeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v localeValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	val := req.ConfigValue.ValueString()
	if val == "$LANG" {
		return
	}
	if posixLocalePattern.MatchString(val) {
		return
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Locale",
		fmt.Sprintf(
			"Locale must be \"$LANG\" (default, uses client locale) or a valid POSIX locale string "+
				"(e.g. \"en_US.UTF-8\", \"fr_FR.UTF-8\", \"ja_JP.UTF-8\"). Got: %q", val,
		),
	)
}

// ---------------------------------------------------------------------------
// Timezone validator – accepts $TZ or valid IANA timezone identifiers
// ---------------------------------------------------------------------------

var ianaTimezonePattern = regexp.MustCompile(
	`^[A-Za-z][A-Za-z0-9_+-]+(/[A-Za-z][A-Za-z0-9_+-]+){0,2}$`,
)

type timezoneValidator struct{}

func TimezoneValidator() timezoneValidator {
	return timezoneValidator{}
}

func (v timezoneValidator) Description(_ context.Context) string {
	return `Must be "$TZ" (use client's timezone) or a valid IANA timezone (e.g. "America/New_York", "Europe/London", "UTC").`
}

func (v timezoneValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v timezoneValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	val := req.ConfigValue.ValueString()
	if val == "$TZ" {
		return
	}
	if ianaTimezonePattern.MatchString(val) {
		return
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid Timezone",
		fmt.Sprintf(
			"Timezone must be \"$TZ\" (default, uses client timezone) or a valid IANA timezone "+
				"(e.g. \"America/New_York\", \"Europe/London\", \"UTC\"). "+
				"See https://en.wikipedia.org/wiki/List_of_tz_database_time_zones for a complete list. Got: %q", val,
		),
	)
}

// ---------------------------------------------------------------------------
// Color scheme validator – accepts named presets or Guacamole terminal syntax
// ---------------------------------------------------------------------------

var namedColorSchemes = []string{"black-white", "gray-black", "green-black", "white-black"}

// guacColorEntryPattern matches a single Guacamole color entry like
// "background: rgb:FF/FF/FF;" or "color12: rgb:00/3D/FC;".
var guacColorEntryPattern = regexp.MustCompile(
	`^(background|foreground|color(?:1[0-5]|[0-9])):\s*rgb:[0-9A-Fa-f]{2}/[0-9A-Fa-f]{2}/[0-9A-Fa-f]{2};$`,
)

type colorSchemeValidator struct{}

func ColorSchemeValidator() colorSchemeValidator {
	return colorSchemeValidator{}
}

func (v colorSchemeValidator) Description(_ context.Context) string {
	return "Must be one of the named schemes (" + strings.Join(namedColorSchemes, ", ") +
		") or a valid Guacamole terminal color scheme string (e.g. \"background: rgb:00/00/00;\\nforeground: rgb:FF/FF/FF;\")."
}

func (v colorSchemeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v colorSchemeValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	val := req.ConfigValue.ValueString()

	for _, named := range namedColorSchemes {
		if val == named {
			return
		}
	}

	lines := strings.Split(val, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !guacColorEntryPattern.MatchString(trimmed) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Color Scheme",
				fmt.Sprintf(
					"Color scheme must be one of the named presets (%s) or a valid Guacamole terminal color scheme.\n"+
						"Invalid line: %q\n"+
						"Example of valid Guacamole color scheme:\n"+
						"  \"background: rgb:00/3D/FC;\\nforeground: rgb:74/1A/1A;\\ncolor0: rgb:00/00/00;\\ncolor1: rgb:99/3E/3E;\\n"+
						"color2: rgb:3E/99/3E;\\ncolor3: rgb:99/99/3E;\\ncolor4: rgb:3E/3E/99;\\ncolor5: rgb:99/3E/99;\\n"+
						"color6: rgb:3E/99/99;\\ncolor7: rgb:99/99/99;\\ncolor8: rgb:3E/3E/3E;\\ncolor9: rgb:FF/67/67;\\n"+
						"color10: rgb:67/FF/67;\\ncolor11: rgb:FF/FF/67;\\ncolor12: rgb:67/67/FF;\\ncolor13: rgb:FF/67/FF;\\n"+
						"color14: rgb:67/FF/FF;\\ncolor15: rgb:FF/FF/FF;\"",
					strings.Join(namedColorSchemes, ", "), trimmed,
				),
			)
			return
		}
	}
}

// ---------------------------------------------------------------------------
// allow_supply_host vs hostname_or_ip config validation
// ---------------------------------------------------------------------------

// IsHostnameOrIPSet reports whether the hostname_or_ip block is present in config.
func IsHostnameOrIPSet(m *HostnameOrIPModel) bool {
	if m == nil {
		return false
	}
	return isHostnameSet(m) ||
		(!m.AdministrativePort.IsNull() && !m.AdministrativePort.IsUnknown())
}

func isHostnameSet(m *HostnameOrIPModel) bool {
	if m == nil {
		return false
	}
	return !m.HostName.IsNull() && !m.HostName.IsUnknown() && strings.TrimSpace(m.HostName.ValueString()) != ""
}

func allowSupplyHostEnabled(v types.Bool) bool {
	return !v.IsNull() && !v.IsUnknown() && v.ValueBool()
}

// ValidateAllowSupplyHostHostnameConstraints enforces mutual exclusion between
// pam_settings.allow_supply_host and hostname_or_ip. When allow_supply_host is
// false or unset, hostname_or_ip.hostname is required.
func ValidateAllowSupplyHostHostnameConstraints(allowSupplyHost types.Bool, hostnameOrIP *HostnameOrIPModel) diag.Diagnostics {
	if allowSupplyHostEnabled(allowSupplyHost) {
		if !IsHostnameOrIPSet(hostnameOrIP) {
			return nil
		}

		var diags diag.Diagnostics
		diags.AddAttributeError(
			path.Root("hostname_or_ip"),
			"Invalid Hostname Configuration",
			"hostname_or_ip must not be set when pam_settings.allow_supply_host is true.",
		)
		diags.AddAttributeError(
			path.Root("pam_settings").AtName("allow_supply_host"),
			"Invalid PAM Settings Configuration",
			"allow_supply_host cannot be true when hostname_or_ip is set.",
		)
		return diags
	}

	if isHostnameSet(hostnameOrIP) {
		return nil
	}

	var diags diag.Diagnostics
	diags.AddAttributeError(
		path.Root("hostname_or_ip").AtName("hostname"),
		"Missing Required Attribute",
		"hostname is required when pam_settings.allow_supply_host is not true.",
	)
	return diags
}

// AllowSupplyHostFromMachineDirectoryPamSettings extracts allow_supply_host
// from a machine/directory pam_settings block.
func AllowSupplyHostFromMachineDirectoryPamSettings(settings *MachineDirectoryPamSettingsFieldResourceModel) types.Bool {
	if settings == nil {
		return types.BoolNull()
	}
	return settings.AllowSupplyHost
}

// AllowSupplyHostFromDatabasePamSettings extracts allow_supply_host from a
// database pam_settings block.
func AllowSupplyHostFromDatabasePamSettings(settings *DatabasePamSettingsFieldResourceModel) types.Bool {
	if settings == nil {
		return types.BoolNull()
	}
	return settings.AllowSupplyHost
}

// AllowSupplyHostHostnameConfigValidatorDescription returns the shared
// description for allow_supply_host vs hostname_or_ip config validators.
func AllowSupplyHostHostnameConfigValidatorDescription(_ context.Context) string {
	return "When pam_settings.allow_supply_host is true, hostname_or_ip must not be set; otherwise hostname_or_ip.hostname is required."
}

type allowSupplyHostHostnameConfigValidator[T any] struct {
	extract func(T) (types.Bool, *HostnameOrIPModel)
}

// NewAllowSupplyHostHostnameConfigValidator returns a resource config validator
// that rejects hostname_or_ip when pam_settings.allow_supply_host is true.
func NewAllowSupplyHostHostnameConfigValidator[T any](extract func(T) (types.Bool, *HostnameOrIPModel)) resource.ConfigValidator {
	return allowSupplyHostHostnameConfigValidator[T]{extract: extract}
}

func (v allowSupplyHostHostnameConfigValidator[T]) Description(ctx context.Context) string {
	return AllowSupplyHostHostnameConfigValidatorDescription(ctx)
}

func (v allowSupplyHostHostnameConfigValidator[T]) MarkdownDescription(ctx context.Context) string {
	return AllowSupplyHostHostnameConfigValidatorDescription(ctx)
}

func (v allowSupplyHostHostnameConfigValidator[T]) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config T
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	allowSupplyHost, hostnameOrIP := v.extract(config)
	resp.Diagnostics.Append(ValidateAllowSupplyHostHostnameConstraints(allowSupplyHost, hostnameOrIP)...)
}
