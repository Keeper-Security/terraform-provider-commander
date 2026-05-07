// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ----- Policy type (string enum) -----

// PolicyTypeValidator validates that policy type is one of the PolicyTypeValues.
type PolicyTypeValidator struct{}

func (PolicyTypeValidator) Description(ctx context.Context) string {
	return "Policy type must be one of: " + PolicyTypeDescription() + "."
}

func (PolicyTypeValidator) MarkdownDescription(ctx context.Context) string {
	return "Policy type must be one of: " + PolicyTypeMarkdown() + "."
}

func (v PolicyTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	val := strings.TrimSpace(strings.ToLower(req.ConfigValue.ValueString()))
	for _, allowed := range PolicyTypeValues {
		if val == allowed {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid policy type",
		"Policy type must be one of: "+PolicyTypeDescription()+". Got: "+req.ConfigValue.ValueString(),
	)
}

// ----- Status (string enum) -----

// StatusValidator validates that status is one of the StatusValues.
type StatusValidator struct{}

func (StatusValidator) Description(ctx context.Context) string {
	return "Status must be one of: " + StatusDescription() + "."
}

func (StatusValidator) MarkdownDescription(ctx context.Context) string {
	return "Status must be one of: " + StatusMarkdown() + "."
}

func (v StatusValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	val := strings.TrimSpace(strings.ToLower(req.ConfigValue.ValueString()))
	for _, allowed := range StatusValues {
		if val == allowed {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid status",
		"Status must be one of: "+StatusDescription()+". Got: "+req.ConfigValue.ValueString(),
	)
}

// ValidateMonitorAndNotifyOnlyFields rejects message and require_policy_acknowledgement unless status is monitor_and_notify.
func ValidateMonitorAndNotifyOnlyFields(status string, message types.String, requireAck types.Bool, pathMessage, pathRequireAck path.Path, diags *diag.Diagnostics) {
	st := strings.TrimSpace(strings.ToLower(status))
	if st == StatusMonitorAndNotify {
		return
	}
	if !message.IsNull() && !message.IsUnknown() {
		diags.AddAttributeError(
			pathMessage,
			"Invalid message for policy status",
			"The message attribute is only allowed when status is "+StatusMonitorAndNotify+".",
		)
	}
	if !requireAck.IsNull() && !requireAck.IsUnknown() {
		diags.AddAttributeError(
			pathRequireAck,
			"Invalid require_policy_acknowledgement for policy status",
			"The require_policy_acknowledgement attribute is only allowed when status is "+StatusMonitorAndNotify+".",
		)
	}
}

// ----- Control (set of enum values) -----

// ControlSetValidator validates that each element in the control set is one of the ControlValues.
type ControlSetValidator struct{}

func (ControlSetValidator) Description(ctx context.Context) string {
	return "Each control value must be one of: " + ControlDescription() + "."
}

func (ControlSetValidator) MarkdownDescription(ctx context.Context) string {
	return "Each control value must be one of: " + ControlMarkdown() + "."
}

func (v ControlSetValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for _, el := range req.ConfigValue.Elements() {
		strVal, ok := el.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid control element",
				"Each control value must be a string. One of: "+ControlDescription()+".",
			)
			continue
		}
		if strVal.IsNull() || strVal.IsUnknown() {
			continue
		}
		val := strings.TrimSpace(strings.ToLower(strVal.ValueString()))
		if val == "" {
			continue
		}
		found := false
		for _, allowed := range ControlValues {
			if val == allowed {
				found = true
				break
			}
		}
		if !found {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtSetValue(strVal),
				"Invalid control value",
				"Control must be one of: "+ControlDescription()+". Got: "+strVal.ValueString(),
			)
		}
	}
}

// ----- Day filter (set of weekdays, case-insensitive) -----

// DayFilterSetValidator validates that each element in the day filter set is one of the weekdays (case-insensitive).
type DayFilterSetValidator struct{}

func (DayFilterSetValidator) Description(ctx context.Context) string {
	return "Each day filter value must be one of: " + DayFilterDescription() + " (case-insensitive)."
}

func (DayFilterSetValidator) MarkdownDescription(ctx context.Context) string {
	return "Each day filter value must be one of: " + DayFilterMarkdown() + " (case-insensitive)."
}

func (v DayFilterSetValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for _, el := range req.ConfigValue.Elements() {
		strVal, ok := el.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid day filter element",
				"Each day filter value must be one of: "+DayFilterDescription()+" (case-insensitive).",
			)
			continue
		}
		if strVal.IsNull() || strVal.IsUnknown() {
			continue
		}
		val := strings.TrimSpace(strings.ToLower(strVal.ValueString()))
		if val == "" {
			continue // optional; empty treated as unset
		}
		found := false
		for _, allowed := range DayFilterValues {
			if val == strings.ToLower(allowed) {
				found = true
				break
			}
		}
		if !found {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtSetValue(strVal),
				"Invalid day filter value",
				"Day filter must be one of: "+DayFilterDescription()+" (case-insensitive). Got: "+strVal.ValueString(),
			)
		}
	}
}

// ----- Time filter list: format start-end (hours 0–23) and no overlapping ranges -----

// timeHourRangeToHalfOpenMinutes parses "H-H" or "HH-HH" (hours 0–23 inclusive on each side)
// into half-open minutes-from-midnight [start, end) for overlap checks. end is exclusive.
func timeHourRangeToHalfOpenMinutes(s string) (startM, endM int, ok bool) {
	parts := strings.SplitN(strings.TrimSpace(s), "-", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}
	sh, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	eh, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	if sh < 0 || sh > 23 || eh < 0 || eh > 23 || sh > eh {
		return 0, 0, false
	}
	startM = sh * 60
	endM = (eh + 1) * 60
	return startM, endM, true
}

// TimeFilterSetValidator validates each element is start-end (hours 0–23) and that no two ranges overlap.
type TimeFilterSetValidator struct{}

func (TimeFilterSetValidator) Description(ctx context.Context) string {
	return "Each time filter must be a range of hours 0–23 in format start-end (e.g. 9-12). Ranges must not overlap."
}

func (TimeFilterSetValidator) MarkdownDescription(ctx context.Context) string {
	return "Each time filter must be a range of hours **0–23** in format **start-end** (e.g. `9-12`). Ranges must **not overlap**."
}

func (v TimeFilterSetValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	elements := req.ConfigValue.Elements()
	var ranges []struct{ start, end int }
	for _, elem := range elements {
		strVal, ok := elem.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid time filter element",
				"Each element must be a string in format start-end (hours 0–23, e.g. 9-12).",
			)
			continue
		}
		if strVal.IsNull() || strVal.IsUnknown() {
			continue
		}
		s := strings.TrimSpace(strVal.ValueString())
		if s == "" {
			continue
		}
		start, end, ok := timeHourRangeToHalfOpenMinutes(s)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtSetValue(strVal),
				"Invalid time filter format",
				"Time filter must be start-end with hours from 0 to 23 (e.g. 9-12). Got: "+s,
			)
			continue
		}
		ranges = append(ranges, struct{ start, end int }{start, end})
	}
	// Check overlaps: half-open [a,b) and [c,d) overlap if a < d && c < b
	for i := 0; i < len(ranges); i++ {
		for j := i + 1; j < len(ranges); j++ {
			a, b := ranges[i].start, ranges[i].end
			c, d := ranges[j].start, ranges[j].end
			if a < d && c < b {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Overlapping time filter ranges",
					"Time filter ranges must not overlap. Found overlap between entries.",
				)
				return
			}
		}
	}
}

// ----- Date filter list: format YYYY-MM-DD:YYYY-MM-DD and no overlapping ranges -----

func parseDateRange(s string) (start, end time.Time, ok bool) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, false
	}
	start, err1 := time.Parse(DateFilterDateFormat, strings.TrimSpace(parts[0]))
	end, err2 := time.Parse(DateFilterDateFormat, strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil {
		return time.Time{}, time.Time{}, false
	}
	if start.After(end) {
		return time.Time{}, time.Time{}, false
	}
	return start, end, true
}

// DateFilterSetValidator validates each element is YYYY-MM-DD:YYYY-MM-DD and that no two ranges overlap.
type DateFilterSetValidator struct{}

func (DateFilterSetValidator) Description(ctx context.Context) string {
	return "Each date filter must be in ISO format YYYY-MM-DD:YYYY-MM-DD. Ranges must not overlap."
}

func (DateFilterSetValidator) MarkdownDescription(ctx context.Context) string {
	return "Each date filter must be in **ISO format** `YYYY-MM-DD:YYYY-MM-DD`. Ranges must **not overlap**."
}

func (v DateFilterSetValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	elements := req.ConfigValue.Elements()
	type dateRange struct{ start, end time.Time }
	var ranges []dateRange
	for _, elem := range elements {
		strVal, ok := elem.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid date filter element",
				"Each element must be a string in format YYYY-MM-DD:YYYY-MM-DD.",
			)
			continue
		}
		if strVal.IsNull() || strVal.IsUnknown() {
			continue
		}
		s := strings.TrimSpace(strVal.ValueString())
		if s == "" {
			continue
		}
		start, end, ok := parseDateRange(s)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtSetValue(strVal),
				"Invalid date filter format",
				"Date filter must be in ISO format YYYY-MM-DD:YYYY-MM-DD (e.g. 2025-01-01:2025-01-31). Got: "+s,
			)
			continue
		}
		ranges = append(ranges, dateRange{start, end})
	}
	for i := 0; i < len(ranges); i++ {
		for j := i + 1; j < len(ranges); j++ {
			a, b := ranges[i].start, ranges[i].end
			c, d := ranges[j].start, ranges[j].end
			if a.Before(d) && c.Before(b) {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Overlapping date filter ranges",
					"Date filter ranges must not overlap. Found overlap between entries.",
				)
				return
			}
		}
	}
}

// CollectionWildcardAll is the single value that means "select all" (cannot be combined with specific IDs).
const CollectionWildcardAll = "*"

// ----- User groups / machine collections / applications: set of strings ("*" or IDs) -----

// CollectionSetValidator validates a set of strings: no empty strings; either "*" alone or collection IDs only (not both).
type CollectionSetValidator struct {
	DisplayName string
}

func (v CollectionSetValidator) Description(ctx context.Context) string {
	return v.DisplayName + " must be \"*\" to select all, or a set of collection IDs. Cannot use \"*\" together with specific IDs. No empty strings. Duplicates are avoided by using a set."
}

func (v CollectionSetValidator) MarkdownDescription(ctx context.Context) string {
	return v.DisplayName + " must be **\"*\"** to select all, or a set of collection IDs. Cannot use **\"*\"** together with specific IDs. No empty strings."
}

func (v CollectionSetValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	var hasWildcard bool
	var nonWildcardCount int
	for _, el := range req.ConfigValue.Elements() {
		strVal, ok := el.(types.String)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid "+v.DisplayName+" element",
				"Each element must be a string (\"*\" or a collection ID).",
			)
			continue
		}
		if strVal.IsNull() || strVal.IsUnknown() {
			continue
		}
		val := strings.TrimSpace(strVal.ValueString())
		if val == "" {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Empty "+v.DisplayName+" value",
				v.DisplayName+" set cannot contain empty strings. Use \"*\" or a valid ID.",
			)
			continue
		}
		if val == CollectionWildcardAll {
			hasWildcard = true
		} else {
			nonWildcardCount++
		}
	}
	if hasWildcard && nonWildcardCount > 0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Cannot mix \"*\" with specific IDs",
			"When using \"*\" to select all, do not add other "+v.DisplayName+" IDs. Use either \"*\" alone or a set of specific IDs, not both.",
		)
	}
}

// ----- Helpers for policy-type allowed fields (used by resource config validator) -----

// IsPolicyTypeLeastPrivilege returns true if policyType is least_privilege (case-insensitive).
func IsPolicyTypeLeastPrivilege(policyType string) bool {
	return strings.TrimSpace(strings.ToLower(policyType)) == PolicyTypeLeastPrivilege
}

// IsPolicyTypeCommand returns true if policyType is command (case-insensitive).
func IsPolicyTypeCommand(policyType string) bool {
	return strings.TrimSpace(strings.ToLower(policyType)) == PolicyTypeCommand
}

// IsListEmptyOrNull returns true if the list is null, unknown, or has no elements (or only null/empty).
func IsListEmptyOrNull(list types.List) bool {
	if list.IsNull() || list.IsUnknown() {
		return true
	}
	for _, el := range list.Elements() {
		if s, ok := el.(types.String); ok && !s.IsNull() && strings.TrimSpace(s.ValueString()) != "" {
			return false
		}
	}
	return true
}

// IsStringEmptyOrNull returns true if the string is null, unknown, or empty/whitespace.
func IsStringEmptyOrNull(s types.String) bool {
	if s.IsNull() || s.IsUnknown() {
		return true
	}
	return strings.TrimSpace(s.ValueString()) == ""
}

// IsSetEmptyOrNull returns true if the set is null, unknown, or has no non-empty elements.
func IsSetEmptyOrNull(set types.Set) bool {
	if set.IsNull() || set.IsUnknown() {
		return true
	}
	for _, el := range set.Elements() {
		if s, ok := el.(types.String); ok && !s.IsNull() && strings.TrimSpace(s.ValueString()) != "" {
			return false
		}
	}
	return true
}

// IsSetPresent returns true if the attribute is explicitly set in config (not null, not unknown).
// Use when a field must not be set at all: e.g. applications = [] should still be rejected for command policy.
func IsSetPresent(set types.Set) bool {
	return !set.IsNull() && !set.IsUnknown()
}

// HasAtLeastOneCollection returns true if at least one of the given sets has at least one non-empty element.
func HasAtLeastOneCollection(userGroups, machineCollections, applications types.Set) bool {
	return !IsSetEmptyOrNull(userGroups) || !IsSetEmptyOrNull(machineCollections) || !IsSetEmptyOrNull(applications)
}

// HasAtLeastOneMachineAndUser returns true if both machine_collections and user_groups have at least one element (for command type).
func HasAtLeastOneMachineAndUser(userGroups, machineCollections types.Set) bool {
	return !IsSetEmptyOrNull(userGroups) && !IsSetEmptyOrNull(machineCollections)
}

// HasAllThreeCollections returns true if user_groups, machine_collections, and applications each have at least one element (elevation + monitor).
func HasAllThreeCollections(userGroups, machineCollections, applications types.Set) bool {
	return !IsSetEmptyOrNull(userGroups) && !IsSetEmptyOrNull(machineCollections) && !IsSetEmptyOrNull(applications)
}

// IsStatusAllowedForLeastPrivilege returns true if status is enforce or off.
func IsStatusAllowedForLeastPrivilege(status string) bool {
	st := strings.TrimSpace(strings.ToLower(status))
	for _, allowed := range AllowedStatusValuesForLeastPrivilege {
		if st == allowed {
			return true
		}
	}
	return false
}

// addNonEmptySetIfPresent reports an error when the set is explicitly configured (including [] or only blank strings)
// but has no non-empty element. Omitting the attribute (null) is allowed.
func addNonEmptySetIfPresent(p path.Path, attrDisplay string, set types.Set, diags *diag.Diagnostics) {
	if set.IsNull() || set.IsUnknown() {
		return
	}
	if !IsSetEmptyOrNull(set) {
		return
	}
	diags.AddAttributeError(
		p,
		attrDisplay+" cannot be empty",
		attrDisplay+" cannot be an empty list ([]) and cannot contain only blank values. Omit this argument or include at least one non-empty value.",
	)
}

// validateSixEpmPolicyOptionalSetsNonEmpty enforces non-empty values when user_groups, machine_collections,
// applications, day_filter, time_filter, or date_filter are set. policyType must be normalized (lowercase).
// For command, applications is never validated here (must not be set at all for command policy).
func validateSixEpmPolicyOptionalSetsNonEmpty(
	policyType string,
	dayFilter, userGroups, machineCollections, applications, timeFilter, dateFilter types.Set,
	pathDayFilter, pathUserGroups, pathMachineCollections, pathApplications, pathTimeFilter, pathDateFilter path.Path,
	diags *diag.Diagnostics,
) {
	addNonEmptySetIfPresent(pathDayFilter, "day_filter", dayFilter, diags)
	addNonEmptySetIfPresent(pathTimeFilter, "time_filter", timeFilter, diags)
	addNonEmptySetIfPresent(pathDateFilter, "date_filter", dateFilter, diags)
	addNonEmptySetIfPresent(pathUserGroups, "user_groups", userGroups, diags)
	addNonEmptySetIfPresent(pathMachineCollections, "machine_collections", machineCollections, diags)
	if policyType != PolicyTypeCommand {
		addNonEmptySetIfPresent(pathApplications, "applications", applications, diags)
	}
}

// ValidatePolicyTypeAllowedFields validates allowed/required fields based on policy type and status.
// Default required: policy_name, policy_type, status. Default optional: day_filter, date_filter, time_filter.
// For elevation and file_access with status enforce: control plus user_groups, machine_collections, and applications
// are required. For command with status enforce: control plus user_groups and machine_collections (applications
// are never allowed for command). day_filter, time_filter, and date_filter remain optional where allowed.
// Policy-type and status combinations enforce additional required fields and disallowed fields.
func ValidatePolicyTypeAllowedFields(
	policyType, status string,
	control, dayFilter types.Set,
	userGroups, machineCollections, applications, timeFilter, dateFilter types.Set,
	pathStatus, pathControl, pathDayFilter, pathUserGroups, pathMachineCollections, pathApplications, pathTimeFilter, pathDateFilter path.Path,
	diags *diag.Diagnostics,
) {
	pt := strings.TrimSpace(strings.ToLower(policyType))
	st := strings.TrimSpace(strings.ToLower(status))

	switch pt {
	case PolicyTypeLeastPrivilege:
		// Only enforce and off allowed for least_privilege.
		if !IsStatusAllowedForLeastPrivilege(st) {
			diags.AddAttributeError(
				pathStatus,
				"Invalid status for least_privilege policy",
				"For policy type "+PolicyTypeLeastPrivilege+", status must be "+StatusDescriptionForLeastPrivilege()+".",
			)
		}
		// Only policy name, type, status required; machine_collections optional. Other fields not allowed.
		// Use IsSetPresent so empty lists (e.g. user_groups = []) count as set and are rejected.
		var disallowedSet []string
		if IsSetPresent(control) {
			disallowedSet = append(disallowedSet, "control")
		}
		if IsSetPresent(dayFilter) {
			disallowedSet = append(disallowedSet, "day_filter")
		}
		if IsSetPresent(userGroups) {
			disallowedSet = append(disallowedSet, "user_groups")
		}
		if IsSetPresent(applications) {
			disallowedSet = append(disallowedSet, "applications")
		}
		if IsSetPresent(timeFilter) {
			disallowedSet = append(disallowedSet, "time_filter")
		}
		if IsSetPresent(dateFilter) {
			disallowedSet = append(disallowedSet, "date_filter")
		}
		if len(disallowedSet) > 0 {
			diags.AddAttributeError(
				pathControl,
				"Fields not allowed for least_privilege policy",
				"For policy type least_privilege, only policy name, policy type, status are allowed and machine_collections is required. The following must not be set (including empty lists []): "+strings.Join(disallowedSet, ", ")+".",
			)
		}
		if IsSetEmptyOrNull(machineCollections) {
			diags.AddAttributeError(
				pathMachineCollections,
				"machine_collections required for least_privilege policy",
				"For policy type least_privilege, at least one machine_collections value is required.",
			)
		}
		return
	}

	// command: applications never allowed (including applications = [])
	if pt == PolicyTypeCommand {
		if IsSetPresent(applications) {
			diags.AddAttributeError(
				pathApplications,
				"Field not allowed for command policy",
				"For policy type command, \"applications\" must not be set (including an empty list).",
			)
		}
	}

	// elevation
	if pt == PolicyTypeElevation {
		switch st {
		case StatusEnforce:
			if IsSetEmptyOrNull(control) {
				diags.AddAttributeError(
					pathControl,
					"Control required for elevation policy with status enforce",
					"For policy type elevation with status enforce, at least one control value is required.",
				)
			}
			if !HasAllThreeCollections(userGroups, machineCollections, applications) {
				diags.AddAttributeError(
					pathUserGroups,
					"Machine, application, and user collection required",
					"For policy type elevation with status enforce, at least one machine collection, one application collection, and one user collection are required.",
				)
			}
		case StatusMonitor, StatusMonitorAndNotify:
			// At least one machine, application, and user collection required; control, day_filter, date_filter, time_filter optional.
			if !HasAllThreeCollections(userGroups, machineCollections, applications) {
				diags.AddAttributeError(
					pathUserGroups,
					"Machine, application, and user collection required",
					"For policy type elevation with status "+st+", at least one machine collection, one application, and one user collection are required to save this policy.",
				)
			}
		case StatusOff:
			// Default required only; machine, application, user collection, control optional. Nothing to validate.
		}
		validateSixEpmPolicyOptionalSetsNonEmpty(pt, dayFilter, userGroups, machineCollections, applications, timeFilter, dateFilter,
			pathDayFilter, pathUserGroups, pathMachineCollections, pathApplications, pathTimeFilter, pathDateFilter, diags)
		return
	}

	// file_access
	if pt == PolicyTypeFileAccess {
		switch st {
		case StatusEnforce:
			if IsSetEmptyOrNull(control) {
				diags.AddAttributeError(
					pathControl,
					"Control required for file_access policy with status enforce",
					"For policy type file_access with status enforce, at least one control value is required.",
				)
			}
			if !HasAllThreeCollections(userGroups, machineCollections, applications) {
				diags.AddAttributeError(
					pathUserGroups,
					"Machine, application, and user collection required",
					"For policy type file_access with status enforce, at least one machine collection, one application collection, and one user collection are required.",
				)
			}
		case StatusMonitor, StatusMonitorAndNotify:
			if !HasAtLeastOneCollection(userGroups, machineCollections, applications) {
				diags.AddAttributeError(
					pathUserGroups,
					"Machine, application, and user collection required",
					"For policy type file_access with status "+st+", at least one machine collection, one application, and one user collection are required to save this policy.",
				)
			}
		case StatusOff:
			// Default required only; machine, application, user collection optional.
		}
		validateSixEpmPolicyOptionalSetsNonEmpty(pt, dayFilter, userGroups, machineCollections, applications, timeFilter, dateFilter,
			pathDayFilter, pathUserGroups, pathMachineCollections, pathApplications, pathTimeFilter, pathDateFilter, diags)
		return
	}

	// command
	if pt == PolicyTypeCommand {
		switch st {
		case StatusEnforce:
			if IsSetEmptyOrNull(control) {
				diags.AddAttributeError(
					pathControl,
					"Control required for command policy with status enforce",
					"For policy type command with status enforce, at least one control value is required.",
				)
			}
			if !HasAtLeastOneMachineAndUser(userGroups, machineCollections) {
				diags.AddAttributeError(
					pathUserGroups,
					"Machine and user collections required",
					"For policy type command with status enforce, at least one user_groups and at least one machine_collections are required.",
				)
			}
		case StatusMonitor, StatusMonitorAndNotify:
			if !HasAtLeastOneMachineAndUser(userGroups, machineCollections) {
				diags.AddAttributeError(
					pathUserGroups,
					"Machine and user collections required",
					"For policy type command with status "+st+", at least one user_groups and at least one machine_collections are required.",
				)
			}
		case StatusOff:
			// Default required only; control, user collection, machine optional.
		}
		validateSixEpmPolicyOptionalSetsNonEmpty(pt, dayFilter, userGroups, machineCollections, applications, timeFilter, dateFilter,
			pathDayFilter, pathUserGroups, pathMachineCollections, pathApplications, pathTimeFilter, pathDateFilter, diags)
	}
}
