// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy_test

import (
	"testing"

	commonepm "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/epm_policy"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func mustSetVals(t *testing.T, vals ...string) types.Set {
	t.Helper()
	elems := make([]attr.Value, len(vals))
	for i, v := range vals {
		elems[i] = types.StringValue(v)
	}
	s, d := types.SetValue(types.StringType, elems)
	if d.HasError() {
		t.Fatal(d)
	}
	return s
}

type pathBundle struct {
	status, control, day, user, machine, apps, time, date path.Path
}

func policyPaths() pathBundle {
	return pathBundle{
		status:  path.Root("status"),
		control: path.Root("control"),
		day:     path.Root("day_filter"),
		user:    path.Root("user_groups"),
		machine: path.Root("machine_collections"),
		apps:    path.Root("applications"),
		time:    path.Root("time_filter"),
		date:    path.Root("date_filter"),
	}
}

func TestValidatePolicyTypeAllowedFields_LeastPrivilege(t *testing.T) {
	t.Parallel()
	basePaths := policyPaths()
	var diags diag.Diagnostics
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeLeastPrivilege, "monitor",
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		basePaths.status, basePaths.control, basePaths.day, basePaths.user, basePaths.machine, basePaths.apps, basePaths.time, basePaths.date,
		&diags,
	)
	if !diags.HasError() {
		t.Fatal("least_privilege + monitor: want status error")
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeLeastPrivilege, commonepm.StatusEnforce,
		mustSetVals(t, "audit"), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		basePaths.status, basePaths.control, basePaths.day, basePaths.user, basePaths.machine, basePaths.apps, basePaths.time, basePaths.date,
		&diags,
	)
	if !diags.HasError() {
		t.Fatal("least_privilege + disallowed control: want error")
	}

	diags = nil
	emptyMc, _ := types.SetValue(types.StringType, []attr.Value{})
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeLeastPrivilege, commonepm.StatusEnforce,
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), emptyMc, types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		basePaths.status, basePaths.control, basePaths.day, basePaths.user, basePaths.machine, basePaths.apps, basePaths.time, basePaths.date,
		&diags,
	)
	if !diags.HasError() {
		t.Fatal("empty machine_collections: want error")
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeLeastPrivilege, commonepm.StatusEnforce,
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		basePaths.status, basePaths.control, basePaths.day, basePaths.user, basePaths.machine, basePaths.apps, basePaths.time, basePaths.date,
		&diags,
	)
	if !diags.HasError() {
		t.Fatal("least_privilege enforce with null machine_collections: want error")
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeLeastPrivilege, commonepm.StatusEnforce,
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), mustSetVals(t, "*"), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		basePaths.status, basePaths.control, basePaths.day, basePaths.user, basePaths.machine, basePaths.apps, basePaths.time, basePaths.date,
		&diags,
	)
	if diags.HasError() {
		t.Fatal("least_privilege enforce with machine_collections: want no error", diags)
	}
}

func TestValidatePolicyTypeAllowedFields_Command(t *testing.T) {
	t.Parallel()
	p := policyPaths()
	var diags diag.Diagnostics

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeCommand, commonepm.StatusEnforce,
		mustSetVals(t, "audit"), types.SetNull(types.StringType),
		mustSetVals(t, "u"), mustSetVals(t, "m"), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if diags.HasError() {
		t.Fatal("command enforce valid", diags)
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeCommand, commonepm.StatusEnforce,
		mustSetVals(t, "allow"), types.SetNull(types.StringType),
		mustSetVals(t, "u"), mustSetVals(t, "m"), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if diags.HasError() {
		t.Fatal("command enforce with allow control: want no error", diags)
	}

	diags = nil
	appsPresent, _ := types.SetValue(types.StringType, []attr.Value{})
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeCommand, commonepm.StatusEnforce,
		mustSetVals(t, "audit"), types.SetNull(types.StringType),
		mustSetVals(t, "u"), mustSetVals(t, "m"), appsPresent,
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if !diags.HasError() {
		t.Fatal("command + applications set (empty): want error")
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeCommand, commonepm.StatusEnforce,
		mustSetVals(t, "audit"), types.SetNull(types.StringType),
		mustSetVals(t, "u"), types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if !diags.HasError() {
		t.Fatal("command enforce missing machine: want error")
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeCommand, commonepm.StatusMonitor,
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		mustSetVals(t, "u"), types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if !diags.HasError() {
		t.Fatal("command monitor missing machine: want error")
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeCommand, commonepm.StatusOff,
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if diags.HasError() {
		t.Fatal("command off: no extra errors", diags)
	}

	diags = nil
	emptyDay, _ := types.SetValue(types.StringType, []attr.Value{})
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeCommand, commonepm.StatusOff,
		types.SetNull(types.StringType), emptyDay,
		types.SetNull(types.StringType), types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if !diags.HasError() {
		t.Fatal("empty day_filter: want error (validateSix)")
	}
}

func TestValidatePolicyTypeAllowedFields_Elevation(t *testing.T) {
	t.Parallel()
	p := policyPaths()
	var diags diag.Diagnostics

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeElevation, commonepm.StatusEnforce,
		mustSetVals(t, "audit"), types.SetNull(types.StringType),
		mustSetVals(t, "u"), mustSetVals(t, "m"), mustSetVals(t, "a"),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if diags.HasError() {
		t.Fatal("elevation enforce valid", diags)
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeElevation, commonepm.StatusEnforce,
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		mustSetVals(t, "u"), mustSetVals(t, "m"), mustSetVals(t, "a"),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if !diags.HasError() {
		t.Fatal("elevation enforce missing control: want error")
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeElevation, commonepm.StatusEnforce,
		mustSetVals(t, "audit"), types.SetNull(types.StringType),
		mustSetVals(t, "u"), mustSetVals(t, "m"), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if !diags.HasError() {
		t.Fatal("elevation enforce missing apps: want error")
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeElevation, commonepm.StatusMonitor,
		mustSetVals(t, "audit"), types.SetNull(types.StringType),
		mustSetVals(t, "u"), mustSetVals(t, "m"), mustSetVals(t, "a"),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if diags.HasError() {
		t.Fatal("elevation monitor valid", diags)
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeElevation, commonepm.StatusMonitor,
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		mustSetVals(t, "u"), types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if !diags.HasError() {
		t.Fatal("elevation monitor incomplete collections: want error")
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeElevation, commonepm.StatusOff,
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if diags.HasError() {
		t.Fatal("elevation off", diags)
	}
}

func TestValidatePolicyTypeAllowedFields_FileAccess(t *testing.T) {
	t.Parallel()
	p := policyPaths()
	var diags diag.Diagnostics

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeFileAccess, commonepm.StatusEnforce,
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		mustSetVals(t, "u"), mustSetVals(t, "m"), mustSetVals(t, "a"),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if !diags.HasError() {
		t.Fatal("file_access enforce without control: want error")
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeFileAccess, commonepm.StatusEnforce,
		mustSetVals(t, "audit"), types.SetNull(types.StringType),
		mustSetVals(t, "u"), mustSetVals(t, "m"), mustSetVals(t, "a"),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if diags.HasError() {
		t.Fatal("file_access enforce valid", diags)
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeFileAccess, commonepm.StatusMonitorAndNotify,
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if !diags.HasError() {
		t.Fatal("file_access monitor with no collections: want error")
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeFileAccess, commonepm.StatusMonitorAndNotify,
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		mustSetVals(t, "u"), types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if diags.HasError() {
		t.Fatal("file_access monitor with user only should pass (OR)", diags)
	}

	diags = nil
	commonepm.ValidatePolicyTypeAllowedFields(
		commonepm.PolicyTypeFileAccess, commonepm.StatusOff,
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if diags.HasError() {
		t.Fatal("file_access off", diags)
	}
}

func TestValidatePolicyTypeAllowedFields_UnknownPolicyType(t *testing.T) {
	t.Parallel()
	p := policyPaths()
	var diags diag.Diagnostics
	commonepm.ValidatePolicyTypeAllowedFields(
		"custom_future_type", commonepm.StatusEnforce,
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType), types.SetNull(types.StringType),
		types.SetNull(types.StringType), types.SetNull(types.StringType),
		p.status, p.control, p.day, p.user, p.machine, p.apps, p.time, p.date,
		&diags,
	)
	if diags.HasError() {
		t.Fatal("unknown policy type: validator currently silent", diags)
	}
}
