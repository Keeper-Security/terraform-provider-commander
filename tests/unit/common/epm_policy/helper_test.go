// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package epmpolicy_test

import (
	"strings"
	"testing"

	commonepm "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/epm_policy"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/google/go-cmp/cmp"
)

func TestPolicyTypeFromAPI(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want string
	}{
		{"PrivilegeElevation", commonepm.PolicyTypeElevation},
		{"FileAccess", commonepm.PolicyTypeFileAccess},
		{"Command", commonepm.PolicyTypeCommand},
		{"LeastPrivilege", commonepm.PolicyTypeLeastPrivilege},
		{"privilege elevation", commonepm.PolicyTypeElevation},
		{"FILE_ACCESS", commonepm.PolicyTypeFileAccess},
		{"least_privilege", commonepm.PolicyTypeLeastPrivilege},
		{"Privilege Elevation", commonepm.PolicyTypeElevation},
	}
	for _, tc := range cases {
		got, err := commonepm.PolicyTypeFromAPI(tc.in)
		if err != nil || got != tc.want {
			t.Errorf("PolicyTypeFromAPI(%q) = %q, %v; want %q, nil", tc.in, got, err, tc.want)
		}
	}
	if _, err := commonepm.PolicyTypeFromAPI(""); err == nil {
		t.Error("empty string: want error")
	}
	if _, err := commonepm.PolicyTypeFromAPI("UnknownThing"); err == nil {
		t.Error("unknown: want error")
	}
}

func TestStatusFromAPI(t *testing.T) {
	t.Parallel()
	if got := commonepm.StatusFromAPI(""); got != "" {
		t.Errorf("empty: got %q", got)
	}
	if got := commonepm.StatusFromAPI("MonitorAndNotify"); got != commonepm.StatusMonitorAndNotify {
		t.Errorf("MonitorAndNotify: got %q", got)
	}
	if got := commonepm.StatusFromAPI("monitor_and_notify"); got != commonepm.StatusMonitorAndNotify {
		t.Errorf("monitor_and_notify: got %q", got)
	}
	if got := commonepm.StatusFromAPI("Enforce"); got != "enforce" {
		t.Errorf("Enforce: got %q", got)
	}
}

func TestControlsFromAPI(t *testing.T) {
	t.Parallel()
	if commonepm.ControlsFromAPI(nil) != nil {
		t.Fatal("nil input")
	}
	got := commonepm.ControlsFromAPI([]string{"NOTIFY", " notify ", "AUDIT", "NOTIFY"})
	if len(got) != 2 || got[0] != "audit" || got[1] != "notify" {
		t.Fatalf("got %#v", got)
	}
	if emptySkip := commonepm.ControlsFromAPI([]string{"", "  ", "audit"}); len(emptySkip) != 1 || emptySkip[0] != "audit" {
		t.Fatalf("got %#v", emptySkip)
	}
}

func TestDayNumbersToTerraformDayFilter(t *testing.T) {
	t.Parallel()
	if got, err := commonepm.DayNumbersToTerraformDayFilter(nil); err != nil || got != nil {
		t.Fatalf("nil: %v, %v", got, err)
	}
	got, err := commonepm.DayNumbersToTerraformDayFilter([]int{1, 1, 2})
	if err != nil || len(got) != 2 || got[0] != commonepm.DayMonday || got[1] != commonepm.DayTuesday {
		t.Fatalf("dedupe sort: %v, %v", got, err)
	}
	if _, err := commonepm.DayNumbersToTerraformDayFilter([]int{99}); err == nil {
		t.Fatal("want error for bad weekday")
	}
}

func TestDateSpansToTerraformDateFilter(t *testing.T) {
	t.Parallel()
	if got, err := commonepm.DateSpansToTerraformDateFilter(nil); err != nil || got != nil {
		t.Fatalf("nil: %v, %v", got, err)
	}
	got, err := commonepm.DateSpansToTerraformDateFilter([]utils.EpmPolicyDateSpan{
		{StartDate: "", EndDate: ""},
		{StartDate: "2025-01-01", EndDate: "2025-01-31"},
	})
	if err != nil || len(got) != 1 || got[0] != "2025-01-01:2025-01-31" {
		t.Fatalf("got %v, %v", got, err)
	}
	if _, err := commonepm.DateSpansToTerraformDateFilter([]utils.EpmPolicyDateSpan{
		{StartDate: "2025-01-01", EndDate: ""},
	}); err == nil {
		t.Fatal("incomplete: want error")
	}
}

func TestTimeSpansToTerraformTimeFilter(t *testing.T) {
	t.Parallel()
	if got, err := commonepm.TimeSpansToTerraformTimeFilter(nil); err != nil || got != nil {
		t.Fatalf("nil: %v, %v", got, err)
	}
	got, err := commonepm.TimeSpansToTerraformTimeFilter([]utils.EpmPolicyTimeSpan{
		{StartTime: "", EndTime: ""},
		{StartTime: "09:00:00", EndTime: "17:30:00"},
	})
	if err != nil || len(got) != 1 || got[0] != "9-17" {
		t.Fatalf("got %v, %v", got, err)
	}
	if _, err := commonepm.TimeSpansToTerraformTimeFilter([]utils.EpmPolicyTimeSpan{
		{StartTime: "09:00", EndTime: ""},
	}); err == nil {
		t.Fatal("incomplete: want error")
	}
	if _, err := commonepm.TimeSpansToTerraformTimeFilter([]utils.EpmPolicyTimeSpan{
		{StartTime: "12:00:00", EndTime: "09:00:00"},
	}); err == nil {
		t.Fatal("start after end: want error")
	}
}

func TestControlsFromPolicyView(t *testing.T) {
	t.Parallel()
	if commonepm.ControlsFromPolicyView(nil) != nil {
		t.Fatal("nil view")
	}
	v := &utils.EpmPolicyResponse{Actions: &utils.EpmPolicyActions{
		OnSuccess: &utils.EpmPolicyOnSuccess{Controls: []string{"audit"}},
	}}
	if got := commonepm.ControlsFromPolicyView(v); len(got) != 1 || got[0] != "audit" {
		t.Fatalf("got %#v", got)
	}
}

func TestRestoreStringSliceOrder(t *testing.T) {
	t.Parallel()
	if got := commonepm.RestoreStringSliceOrder(nil, []string{"a"}); got != nil {
		t.Fatal("nil api")
	}
	got := commonepm.RestoreStringSliceOrder([]string{"b", "a"}, []string{"a", "c"})
	if !cmp.Equal(got, []string{"a", "b"}) {
		t.Fatal(got)
	}
}

func TestRestoreDateFilterSliceOrder(t *testing.T) {
	t.Parallel()
	if got := commonepm.RestoreDateFilterSliceOrder(nil, []string{"2025-01-01:2025-01-02"}); got != nil {
		t.Fatal("nil api")
	}
	api := []string{"2025-01-01:2025-01-31", "2025-02-01:2025-02-28"}
	prior := []string{" 2025-02-01 : 2025-02-28 ", "bad"}
	got := commonepm.RestoreDateFilterSliceOrder(api, prior)
	if !cmp.Equal(got, []string{" 2025-02-01 : 2025-02-28 ", "2025-01-01:2025-01-31"}) {
		t.Fatal(got)
	}
}

func TestMapPolicyViewToAttributes(t *testing.T) {
	t.Parallel()
	if _, err := commonepm.MapPolicyViewToAttributes(nil, nil); err == nil {
		t.Fatal("nil view")
	}
	emptyID := &utils.EpmPolicyResponse{PolicyId: "  "}
	if _, err := commonepm.MapPolicyViewToAttributes(emptyID, nil); err == nil {
		t.Fatal("empty id")
	}
	badType := &utils.EpmPolicyResponse{PolicyId: "1", PolicyType: "???", Status: "enforce"}
	if _, err := commonepm.MapPolicyViewToAttributes(badType, nil); err == nil {
		t.Fatal("bad policy type")
	}
	noStatus := &utils.EpmPolicyResponse{PolicyId: "1", PolicyType: "Command", Status: "  "}
	if _, err := commonepm.MapPolicyViewToAttributes(noStatus, nil); err == nil {
		t.Fatal("empty status")
	}
	badDay := &utils.EpmPolicyResponse{PolicyId: "1", PolicyType: "Command", Status: "enforce", DayCheck: []int{8}}
	if _, err := commonepm.MapPolicyViewToAttributes(badDay, nil); err == nil {
		t.Fatal("bad day")
	}
	view := &utils.EpmPolicyResponse{
		PolicyId:     "42",
		PolicyName:   "N",
		PolicyType:   "Command",
		Status:       "Enforce",
		UserCheck:    []string{"u1"},
		MachineCheck: []string{"m1"},
		DayCheck:     []int{1},
		DateCheck:    []utils.EpmPolicyDateSpan{{StartDate: "2025-06-01", EndDate: "2025-06-30"}},
		TimeCheck:    []utils.EpmPolicyTimeSpan{{StartTime: "10:00:00", EndTime: "11:00:00"}},
		Actions:      &utils.EpmPolicyActions{OnSuccess: &utils.EpmPolicyOnSuccess{Controls: []string{"NOTIFY"}}},
	}
	prior := &commonepm.PolicyViewPriorSets{
		Control:    []string{"notify"},
		UserGroups: []string{"u1"},
		DayFilter:  []string{commonepm.DayMonday},
	}
	mapped, err := commonepm.MapPolicyViewToAttributes(view, prior)
	if err != nil {
		t.Fatal(err)
	}
	if mapped.ID != "42" || mapped.PolicyName != "N" || mapped.PolicyType != commonepm.PolicyTypeCommand || mapped.Status != "enforce" {
		t.Fatalf("%+v", mapped)
	}
	if len(mapped.Control) != 1 || mapped.Control[0] != "notify" {
		t.Fatal(mapped.Control)
	}
	if len(mapped.TimeFilter) != 1 || mapped.TimeFilter[0] != "10-11" {
		t.Fatalf("TimeFilter: %#v", mapped.TimeFilter)
	}

	noActions, err := commonepm.MapPolicyViewToAttributes(&utils.EpmPolicyResponse{
		PolicyId: "7", PolicyType: "FileAccess", Status: "off",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(noActions.Control) != 0 {
		t.Fatal(noActions.Control)
	}
}

func TestBuildPolicyViewCommand(t *testing.T) {
	t.Parallel()
	got := commonepm.BuildPolicyViewCommand("  abc  ")
	if !strings.Contains(got, "abc") || !strings.Contains(got, "view") {
		t.Fatal(got)
	}
}

func TestStringSliceToStringSet(t *testing.T) {
	t.Parallel()
	s, err := commonepm.StringSliceToStringSet(nil)
	if err != nil || !s.IsNull() {
		t.Fatalf("nil slice: %v, null=%v", err, s.IsNull())
	}
	s, err = commonepm.StringSliceToStringSet([]string{"a", "b"})
	if err != nil || s.IsNull() {
		t.Fatal(err)
	}
}
