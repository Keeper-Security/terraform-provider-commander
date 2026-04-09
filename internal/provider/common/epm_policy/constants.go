// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

// PolicyType values for EPM policy.
const (
	PolicyTypeElevation      = "elevation"
	PolicyTypeFileAccess     = "file_access"
	PolicyTypeCommand        = "command"
	PolicyTypeLeastPrivilege = "least_privilege"
)

// PolicyTypeValues is the set of allowed policy types (for validation).
var PolicyTypeValues = []string{
	PolicyTypeElevation, PolicyTypeFileAccess, PolicyTypeCommand, PolicyTypeLeastPrivilege,
}

// ApiPolicyTypeToTerraform maps Commander API policy type strings to Terraform policy_type values.
var ApiPolicyTypeToTerraformMap = map[string]string{
	"PrivilegeElevation": PolicyTypeElevation,
	"LeastPrivilege":     PolicyTypeLeastPrivilege,
	"FileAccess":         PolicyTypeFileAccess,
	"CommandLine":        PolicyTypeCommand,
}

// Status values for EPM policy.
const (
	StatusEnforce          = "enforce"
	StatusMonitor          = "monitor"
	StatusMonitorAndNotify = "monitor_and_notify"
	StatusOff              = "off"
)

// StatusValues is the set of allowed status values.
var StatusValues = []string{
	StatusEnforce, StatusMonitor, StatusMonitorAndNotify, StatusOff,
}

// AllowedStatusValuesForLeastPrivilege is the set of statuses allowed when policy type is least_privilege.
var AllowedStatusValuesForLeastPrivilege = []string{
	StatusEnforce, StatusOff,
}

// Control values for EPM policy.
const (
	ControlAllow    = "allow"
	ControlDeny     = "deny"
	ControlAudit    = "audit"
	ControlNotify   = "notify"
	ControlMfa      = "mfa"
	ControlJustify  = "justify"
	ControlApproval = "approval"
)

// ControlValues is the set of allowed control values.
var ControlValues = []string{
	ControlAllow, ControlDeny, ControlAudit, ControlNotify, ControlMfa, ControlJustify, ControlApproval,
}

// DayFilter values (case-insensitive in validation; stored capitalized for display).
const (
	DayMonday    = "Monday"
	DayTuesday   = "Tuesday"
	DayWednesday = "Wednesday"
	DayThursday  = "Thursday"
	DayFriday    = "Friday"
	DaySaturday  = "Saturday"
	DaySunday    = "Sunday"
)

// Weekday numbers (American convention: 0 = Sunday, 1 = Monday, …, 6 = Saturday).
// Use when converting to/from API which returns an array of weekday numbers.
const (
	WeekdaySunday    = 0
	WeekdayMonday    = 1
	WeekdayTuesday   = 2
	WeekdayWednesday = 3
	WeekdayThursday  = 4
	WeekdayFriday    = 5
	WeekdaySaturday  = 6
)

// DayFilterValues is the set of allowed day filter values.
var DayFilterValues = []string{
	DayMonday, DayTuesday, DayWednesday, DayThursday, DayFriday, DaySaturday, DaySunday,
}

// DayNameToWeekdayNumber maps full day names to weekday numbers (American convention).
// Use when sending day_filter to the API. Keys are the Day* constants (e.g. "Monday" -> 1).
var DayNameToWeekdayNumber = map[string]int{
	DaySunday: WeekdaySunday, DayMonday: WeekdayMonday, DayTuesday: WeekdayTuesday,
	DayWednesday: WeekdayWednesday, DayThursday: WeekdayThursday, DayFriday: WeekdayFriday,
	DaySaturday: WeekdaySaturday,
}

// WeekdayNumberToDayName maps weekday numbers to full day names (American convention).
// Use when reading from the API (array of weekday numbers) to get day names for state.
var WeekdayNumberToDayName = map[int]string{
	WeekdaySunday: DaySunday, WeekdayMonday: DayMonday, WeekdayTuesday: DayTuesday,
	WeekdayWednesday: DayWednesday, WeekdayThursday: DayThursday, WeekdayFriday: DayFriday,
	WeekdaySaturday: DaySaturday,
}

// TimeFilterFormat describes each time_filter value: start-end hours (0–23), e.g. "9-12".
const TimeFilterFormat = "start-end"

// DateFilterFormat is ISO date range "YYYY-MM-DD:YYYY-MM-DD".
const DateFilterDateFormat = "2006-01-02"

// Commander CLI base command for EPM policy operations.
const CmdEpmPolicyBase = "epm policy"

// Commander CLI commands for EPM policy operations (subcommands; use with CmdEpmPolicyBase or full command).
const (
	CmdEpmPolicyAdd    = CmdEpmPolicyBase + " add"
	CmdEpmPolicyEdit   = CmdEpmPolicyBase + " edit"
	CmdEpmPolicyView   = CmdEpmPolicyBase + " view"
	CmdEpmPolicyDelete = CmdEpmPolicyBase + " delete"
)

// Command flags for epm policy add.
const (
	FlagUserFilter             = "--user-filter"             // Policy user filter. User collection UID or *
	FlagMachineFilter          = "--machine-filter"          // Policy machine filter. Machine collection UID
	FlagAppFilter              = "--app-filter"              // Policy application filter. Application collection UID
	FlagDateFilter             = "--date-filter"             // Policy date filter. Date range in ISO format. YYYY-MM-DD:YYYY-MM-DD
	FlagTimeFilter             = "--time-filter"             // Policy time filter. Hours 0–23 as start-end (e.g. 9-12)
	FlagDayFilter              = "--day-filter"              // Policy day filter. Day of Week
	FlagRiskLevel              = "--risk-level"              // Policy risk level
	FlagPolicyType             = "--policy-type"             // Policy type
	FlagPolicyName             = "--policy-name"             // Policy name
	FlagControl                = "--control"                 // Control action
	FlagStatus                 = "--status"                  // Policy status
	FlagEnable                 = "--enable"                  // Policy enable
	FlagMessage                = "--message"                 // Policy message (only for monitor_and_notify status)
	FlagRequireAcknowledgement = "--require-acknowledgement" // Policy require acknowledgement (only for monitor_and_notify status)
)

// Flag for view command output format.
const (
	FlagFormat      = "--format"
	ValueFormatJSON = "json"
)

// Error summaries (first argument to AddError).
const (
	ErrSummaryReadFailed = "Read EPM Policy Failed"
)

// Error operation messages (second argument to ExecuteCommand; short description for logs).
const (
	ErrOpReadEpmPolicy = "Unable to read EPM policy"
)
