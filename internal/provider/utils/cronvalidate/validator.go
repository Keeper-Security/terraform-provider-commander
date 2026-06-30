// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package cronvalidate

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// RotationCronString implements validator.String for attributes that hold a
// Keeper rotation-style cron value. Use [RotationCronString.Name] for the
// Terraform attribute name in diagnostics (e.g. "schedule_cron").
type RotationCronString struct {
	Name string
}

func (v RotationCronString) diagLabel() string { return AttributeLabel(v.Name) }

func (RotationCronString) Description(_ context.Context) string {
	return "must be a valid Keeper Quartz cron expression (6 or 7 fields) with at least a 1-hour interval between executions"
}

func (RotationCronString) MarkdownDescription(_ context.Context) string {
	return "Must be a valid **Keeper Quartz cron expression** (6 or 7 fields per the [Keeper cron spec](https://docs.keeper.io/keeperpam/privileged-access-manager/references/cron-spec)). Rotation schedules require at least a **1-hour interval** between executions."
}

func (v RotationCronString) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	s := strings.TrimSpace(req.ConfigValue.ValueString())
	l := v.diagLabel()
	if s == "" {
		resp.Diagnostics.AddAttributeError(req.Path, fmt.Sprintf("Invalid %s", l),
			fmt.Sprintf("%s cannot be empty; omit the attribute if you do not want a cron schedule.", l))
		return
	}
	if err := ValidateKeeperRotationCron(s); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, fmt.Sprintf("Invalid %s", l), err.Error())
	}
}
