// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PolicyTypeFromAPI maps Commander policy type strings to Terraform policy_type values.
func PolicyTypeFromAPI(api string) (string, error) {
	key := strings.TrimSpace(api)
	if key == "" {
		return "", fmt.Errorf("empty PolicyType in API response")
	}
	// Exact common API enums
	switch key {
	case "PrivilegeElevation":
		return PolicyTypeElevation, nil
	case "FileAccess":
		return PolicyTypeFileAccess, nil
	case "Command":
		return PolicyTypeCommand, nil
	case "LeastPrivilege":
		return PolicyTypeLeastPrivilege, nil
	}
	lower := strings.ToLower(strings.ReplaceAll(key, " ", ""))
	switch lower {
	case "privilegeelevation", "elevation", "privilege_elevation":
		return PolicyTypeElevation, nil
	case "fileaccess", "file_access":
		return PolicyTypeFileAccess, nil
	case "command":
		return PolicyTypeCommand, nil
	case "leastprivilege", "least_privilege":
		return PolicyTypeLeastPrivilege, nil
	}
	return "", fmt.Errorf("unknown API policy type %q", api)
}

// StatusFromAPI maps API status to Terraform status values.
func StatusFromAPI(api string) string {
	t := strings.TrimSpace(api)
	if t == "" {
		return ""
	}
	lower := strings.ToLower(t)
	switch lower {
	case "monitorandnotify", "monitor_and_notify":
		return StatusMonitorAndNotify
	default:
		return lower
	}
}

// ControlsFromAPI maps Actions.OnSuccess.Controls (e.g. NOTIFY) to Terraform control values (notify).
func ControlsFromAPI(controls []string) []string {
	if len(controls) == 0 {
		return nil
	}
	out := make([]string, 0, len(controls))
	seen := make(map[string]struct{})
	for _, c := range controls {
		v := strings.ToLower(strings.TrimSpace(c))
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

// DayNumbersToTerraformDayFilter converts DayCheck weekday numbers to Terraform day names (sorted Mon→Sun order).
func DayNumbersToTerraformDayFilter(nums []int) ([]string, error) {
	if len(nums) == 0 {
		return nil, nil
	}
	seen := make(map[int]struct{})
	var names []string
	for _, n := range nums {
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		name, ok := WeekdayNumberToDayName[n]
		if !ok {
			return nil, fmt.Errorf("unknown weekday number %d in DayCheck", n)
		}
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		ai := DayNameToWeekdayNumber[names[i]]
		aj := DayNameToWeekdayNumber[names[j]]
		return ai < aj
	})
	return names, nil
}

// DateSpansToTerraformDateFilter converts DateCheck to "YYYY-MM-DD:YYYY-MM-DD" strings (sorted).
func DateSpansToTerraformDateFilter(spans []utils.EpmPolicyDateSpan) ([]string, error) {
	if len(spans) == 0 {
		return nil, nil
	}
	out := make([]string, 0, len(spans))
	for _, s := range spans {
		if strings.TrimSpace(s.StartDate) == "" && strings.TrimSpace(s.EndDate) == "" {
			continue
		}
		if s.StartDate == "" || s.EndDate == "" {
			return nil, fmt.Errorf("incomplete date range: StartDate=%q EndDate=%q", s.StartDate, s.EndDate)
		}
		out = append(out, s.StartDate+":"+s.EndDate)
	}
	sort.Strings(out)
	return out, nil
}

// canonicalDateFilterRange parses a Terraform date_filter element "YYYY-MM-DD:YYYY-MM-DD" into a canonical key.
func canonicalDateFilterRange(r string) (string, error) {
	r = strings.TrimSpace(r)
	if r == "" {
		return "", fmt.Errorf("empty date range")
	}
	idx := strings.Index(r, ":")
	if idx < 0 {
		return "", fmt.Errorf("date range must contain ':'")
	}
	startStr := strings.TrimSpace(r[:idx])
	endStr := strings.TrimSpace(r[idx+1:])
	startT, err := time.Parse(DateFilterDateFormat, startStr)
	if err != nil {
		return "", fmt.Errorf("start date %q: %w", startStr, err)
	}
	endT, err := time.Parse(DateFilterDateFormat, endStr)
	if err != nil {
		return "", fmt.Errorf("end date %q: %w", endStr, err)
	}
	return startT.UTC().Format(DateFilterDateFormat) + ":" + endT.UTC().Format(DateFilterDateFormat), nil
}

// RestoreDateFilterSliceOrder matches API-derived ranges to prior state by calendar dates (not only exact string).
// When DateCheck uses epoch ms, API values become canonical "YYYY-MM-DD:YYYY-MM-DD"; prior config strings that
// denote the same range are kept verbatim ("original input format").
func RestoreDateFilterSliceOrder(apiRanges []string, prior []string) []string {
	if len(apiRanges) == 0 {
		return nil
	}
	available := make(map[string]int, len(apiRanges))
	for _, r := range apiRanges {
		available[r]++
	}
	var out []string
	for _, p := range prior {
		canon, err := canonicalDateFilterRange(p)
		if err != nil {
			continue
		}
		if available[canon] > 0 {
			out = append(out, p)
			available[canon]--
		}
	}
	var rest []string
	for r, n := range available {
		for i := 0; i < n; i++ {
			rest = append(rest, r)
		}
	}
	sort.Strings(rest)
	return append(out, rest...)
}

// clockToHHMM strips seconds from API times like "09:00:00" to Terraform "09:00".
func clockToHHMM(clock string) string {
	clock = strings.TrimSpace(clock)
	parts := strings.Split(clock, ":")
	if len(parts) >= 2 {
		return parts[0] + ":" + parts[1]
	}
	return clock
}

// TimeSpansToTerraformTimeFilter converts TimeCheck to "HH:MM-HH:MM" strings (sorted).
func TimeSpansToTerraformTimeFilter(spans []utils.EpmPolicyTimeSpan) ([]string, error) {
	if len(spans) == 0 {
		return nil, nil
	}
	out := make([]string, 0, len(spans))
	for _, s := range spans {
		start := clockToHHMM(s.StartTime)
		end := clockToHHMM(s.EndTime)
		if start == "" && end == "" {
			continue
		}
		if start == "" || end == "" {
			return nil, fmt.Errorf("incomplete time range: StartTime=%q EndTime=%q", s.StartTime, s.EndTime)
		}
		out = append(out, start+"-"+end)
	}
	sort.Strings(out)
	return out, nil
}

// ControlsFromPolicyView returns control strings from Actions.OnSuccess.Controls.
func ControlsFromPolicyView(view *utils.EpmPolicyResponse) []string {
	if view == nil || view.Actions == nil || view.Actions.OnSuccess == nil {
		return nil
	}
	return ControlsFromAPI(view.Actions.OnSuccess.Controls)
}

// PolicyViewPriorSets holds prior config/state string order per set attribute (optional).
// Used by MapPolicyViewToAttributes so refresh keeps the same element order when values match the API.
// A future data source can pass nil to use API/default ordering only.
type PolicyViewPriorSets struct {
	Control            []string
	UserGroups         []string
	MachineCollections []string
	Applications       []string
	DayFilter          []string
	DateFilter         []string
	TimeFilter         []string
}

// PolicyViewMappedAttributes is the API view mapped to Terraform string values (resource and data source).
type PolicyViewMappedAttributes struct {
	ID                 string
	PolicyName         string
	PolicyType         string
	Status             string
	Control            []string
	UserGroups         []string
	MachineCollections []string
	Applications       []string
	DayFilter          []string
	DateFilter         []string
	TimeFilter         []string
}

// RestoreStringSliceOrder builds a slice from apiVals, preserving order from prior for values that still exist
// in the API, then appends any remaining API values in API order.
func RestoreStringSliceOrder(apiVals []string, prior []string) []string {
	if len(apiVals) == 0 {
		return nil
	}
	remaining := make(map[string]struct{}, len(apiVals))
	for _, v := range apiVals {
		remaining[v] = struct{}{}
	}
	var ordered []string
	for _, v := range prior {
		if _, ok := remaining[v]; ok {
			ordered = append(ordered, v)
			delete(remaining, v)
		}
	}
	for _, v := range apiVals {
		if _, ok := remaining[v]; ok {
			ordered = append(ordered, v)
			delete(remaining, v)
		}
	}
	return ordered
}

// MapPolicyViewToAttributes maps `epm policy view --format json` into Terraform-oriented strings.
// prior may be nil (e.g. data source with no previous config).
func MapPolicyViewToAttributes(view *utils.EpmPolicyResponse, prior *PolicyViewPriorSets) (*PolicyViewMappedAttributes, error) {
	if view == nil {
		return nil, fmt.Errorf("nil policy view response")
	}
	id := strings.TrimSpace(view.PolicyId)
	if id == "" {
		return nil, fmt.Errorf("empty PolicyId in API response")
	}

	policyType, err := PolicyTypeFromAPI(view.PolicyType)
	if err != nil {
		return nil, err
	}

	status := StatusFromAPI(view.Status)
	if status == "" {
		return nil, fmt.Errorf("empty Status in API response")
	}

	dayNames, err := DayNumbersToTerraformDayFilter(view.DayCheck)
	if err != nil {
		return nil, err
	}
	dateRanges, err := DateSpansToTerraformDateFilter(view.DateCheck)
	if err != nil {
		return nil, err
	}
	timeRanges, err := TimeSpansToTerraformTimeFilter(view.TimeCheck)
	if err != nil {
		return nil, err
	}

	controlStrs := ControlsFromPolicyView(view)

	var p PolicyViewPriorSets
	if prior != nil {
		p = *prior
	}

	return &PolicyViewMappedAttributes{
		ID:                 id,
		PolicyName:         strings.TrimSpace(view.PolicyName),
		PolicyType:         policyType,
		Status:             status,
		Control:            RestoreStringSliceOrder(controlStrs, p.Control),
		UserGroups:         RestoreStringSliceOrder(view.UserCheck, p.UserGroups),
		MachineCollections: RestoreStringSliceOrder(view.MachineCheck, p.MachineCollections),
		Applications:       RestoreStringSliceOrder(view.ApplicationCheck, p.Applications),
		DayFilter:          RestoreStringSliceOrder(dayNames, p.DayFilter),
		DateFilter:         RestoreDateFilterSliceOrder(dateRanges, p.DateFilter),
		TimeFilter:         RestoreStringSliceOrder(timeRanges, p.TimeFilter),
	}, nil
}

// BuildPolicyViewCommand builds `epm policy view <id> --format json` for the given policy ID.
func BuildPolicyViewCommand(policyID string) string {
	policyID = strings.TrimSpace(policyID)
	return strings.Join([]string{
		CmdEpmPolicyView,
		policyID,
		FlagFormat,
		ValueFormatJSON,
	}, " ")
}

// stringSliceToStringSet builds a Terraform string set from a slice (e.g. after API read mapping).
func StringSliceToStringSet(values []string) (types.Set, error) {
	if len(values) == 0 {
		return types.SetNull(types.StringType), nil
	}
	elems := make([]attr.Value, len(values))
	for i, v := range values {
		elems[i] = types.StringValue(v)
	}
	s, diags := types.SetValue(types.StringType, elems)
	if diags.HasError() {
		return types.SetNull(types.StringType), fmt.Errorf("%s", diags.Errors()[0].Summary())
	}
	return s, nil
}
