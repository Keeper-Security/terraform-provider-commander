// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package sharedfolder

import (
	"context"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Expiration format: "never" | ISO date (yyyy-MM-dd) or datetime (yyyy-MM-dd HH:mm:ss) | relative period (e.g. 30d, 1y, 6mo, 24h, 90days).
var (
	expirationNever    = regexp.MustCompile(`(?i)^never$`)
	expirationRelative = regexp.MustCompile(`^\d+\s*(y|mo|d|h|mi|min|years?|months?|days?|hours?|minutes?)$`)
	// ISO date only (yyyy-MM-dd) or date + optional time (yyyy-MM-dd HH:mm or yyyy-MM-dd HH:mm:ss).
	expirationISO = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}(\s+\d{1,2}:\d{2}(:\d{2})?)?$`)
)

// expirationValidator validates the expiration string for user permission entries.
type expirationValidator struct{}

func (expirationValidator) Description(ctx context.Context) string {
	return ExpirationDoc
}

func (expirationValidator) MarkdownDescription(ctx context.Context) string {
	return ExpirationDoc
}

func (v expirationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	s := req.ConfigValue.ValueString()
	if s == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			ErrMsgInvalidExpiration,
			ErrMsgExpirationEmpty,
		)
		return
	}
	if expirationNever.MatchString(s) {
		return
	}
	if expirationRelative.MatchString(s) {
		return
	}
	if expirationISO.MatchString(s) {
		if _, err := time.Parse(TimeLayoutDate, s); err == nil {
			return
		}
		if _, err := time.Parse(TimeLayoutDateTime, s); err == nil {
			return
		}
		if _, err := time.Parse(TimeLayoutDateTimeShort, s); err == nil {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		ErrMsgInvalidExpiration,
		ExpirationDoc,
	)
}

// ExpirationValidator returns a validator for the expiration attribute.
func ExpirationValidator() expirationValidator {
	return expirationValidator{}
}

// Default user_permissions when null: manage_users = false, manage_records = false.
var userPermissionsDefaultAttrTypes = map[string]attr.Type{
	AttrManageUsers:   types.BoolType,
	AttrManageRecords: types.BoolType,
}

type userPermissionsDefaultPlanModifier struct{}

func (userPermissionsDefaultPlanModifier) Description(ctx context.Context) string {
	return DescUserPermissionsDefault
}

func (userPermissionsDefaultPlanModifier) MarkdownDescription(ctx context.Context) string {
	return DescUserPermissionsDefaultMD
}

func (userPermissionsDefaultPlanModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if !req.ConfigValue.IsNull() {
		return
	}
	defaultVal := types.ObjectValueMust(userPermissionsDefaultAttrTypes, map[string]attr.Value{
		AttrManageUsers:   types.BoolValue(false),
		AttrManageRecords: types.BoolValue(false),
	})
	resp.PlanValue = defaultVal
}

// Default record_permissions when null: can_share = false, can_edit = false.
var recordPermissionsDefaultAttrTypes = map[string]attr.Type{
	AttrCanShare: types.BoolType,
	AttrCanEdit:  types.BoolType,
}

type recordPermissionsDefaultPlanModifier struct{}

func (recordPermissionsDefaultPlanModifier) Description(ctx context.Context) string {
	return DescRecordPermissionsDefault
}

func (recordPermissionsDefaultPlanModifier) MarkdownDescription(ctx context.Context) string {
	return DescRecordPermissionsDefaultMD
}

func (recordPermissionsDefaultPlanModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if !req.ConfigValue.IsNull() {
		return
	}
	defaultVal := types.ObjectValueMust(recordPermissionsDefaultAttrTypes, map[string]attr.Value{
		AttrCanShare: types.BoolValue(false),
		AttrCanEdit:  types.BoolValue(false),
	})
	resp.PlanValue = defaultVal
}
