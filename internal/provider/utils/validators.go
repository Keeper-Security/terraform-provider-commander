// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ----- GENERIC: STRING MIN LENGTH --------------------------------
// StringMinLengthValidator validates that the string has at least MinLen characters after
// strings.TrimSpace (leading and trailing Unicode whitespace is ignored for the count; whitespace-only input is rejected).
// DisplayName is used in error messages (e.g. "Enterprise Node Name").
// AllowNull: if true, null values are skipped (for optional attributes); if false, null is treated as empty.
func StringMinLengthValidator(displayName string, minLen int, allowNull bool) stringMinLengthValidator {
	return stringMinLengthValidator{DisplayName: displayName, MinLen: minLen, AllowNull: allowNull}
}

type stringMinLengthValidator struct {
	DisplayName string
	MinLen      int
	AllowNull   bool
}

func (v stringMinLengthValidator) Description(ctx context.Context) string {
	return v.DisplayName + " must be at least " + strconv.Itoa(v.MinLen) + " character(s) long."
}

func (v stringMinLengthValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v stringMinLengthValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() {
		return
	}
	if v.AllowNull && req.ConfigValue.IsNull() {
		return
	}
	value := strings.TrimSpace(req.ConfigValue.ValueString())
	if len(value) < v.MinLen {
		resp.Diagnostics.AddError(
			"Invalid "+v.DisplayName,
			v.DisplayName+" must be at least "+strconv.Itoa(v.MinLen)+" character(s) long, without leading or trailing whitespace.",
		)
	}
}

// ----- GENERIC: SET NO EMPTY STRINGS --------------------------------
// SetNoEmptyStringsValidator validates that a set of strings contains no empty or whitespace-only strings.
// DisplayName is used in error messages (e.g. "Team", "User", "Role").
func SetNoEmptyStringsValidator(displayName string) setNoEmptyStringsValidator {
	return setNoEmptyStringsValidator{DisplayName: displayName}
}

type setNoEmptyStringsValidator struct {
	DisplayName string
}

func (v setNoEmptyStringsValidator) Description(ctx context.Context) string {
	return v.DisplayName + " set must not contain empty or whitespace-only strings."
}

func (v setNoEmptyStringsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// SetNotEmptyValidator validates that a set has at least one element.
func SetNotEmptyValidator(displayName string) setNotEmptyValidator {
	return setNotEmptyValidator{DisplayName: displayName}
}

type setNotEmptyValidator struct {
	DisplayName string
}

func (v setNotEmptyValidator) Description(ctx context.Context) string {
	return v.DisplayName + " set cannot be empty. At least one value is required."
}

func (v setNotEmptyValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v setNotEmptyValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if len(req.ConfigValue.Elements()) == 0 {
		resp.Diagnostics.AddError(
			"Empty "+v.DisplayName+" Set",
			v.DisplayName+" set cannot be empty. At least one value is required.",
		)
		return
	}
}

func (v setNoEmptyStringsValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for _, elem := range req.ConfigValue.Elements() {
		strValue, ok := elem.(types.String)
		if !ok {
			resp.Diagnostics.AddError(
				"Invalid "+v.DisplayName+" Type",
				fmt.Sprintf("Expected string, got: %T", elem),
			)
			continue
		}
		if strValue.IsUnknown() {
			continue
		}
		if strings.TrimSpace(strValue.ValueString()) == "" {
			resp.Diagnostics.AddError(
				"Empty "+v.DisplayName+" String",
				v.DisplayName+" set cannot contain empty strings. Each "+v.DisplayName+" must be a non-empty string, without leading or trailing whitespace.",
			)
		}
	}
}

// ----- CONVENIENCE: NODE (optional string) --------------------------------.
var NodeValidator = stringMinLengthValidator{DisplayName: "Node Name", MinLen: 1, AllowNull: true}

// ----- CONVENIENCE: MANAGED COMPANY (optional string) --------------------------------.
var ManagedCompanyValidator = stringMinLengthValidator{DisplayName: "Managed Company Name", MinLen: 1, AllowNull: true}

// ----- CONVENIENCE: TEAMS SET --------------------------------.
var TeamsValidator = setNoEmptyStringsValidator{DisplayName: "Team"}

// ----- CONVENIENCE: ROLES SET --------------------------------.
var RolesValidator = setNoEmptyStringsValidator{DisplayName: "Role"}

// ----- GENERIC: MAP KEYS MIN LENGTH --------------------------------
// MapKeysMinLengthValidator validates that all keys in a map have at least MinLen characters
// after strings.TrimSpace (whitespace-only keys are rejected).
// Used for map attributes where keys are identifiers (e.g. managing node names).
func MapKeysMinLengthValidator(displayName string, minLen int) mapKeysMinLengthValidator {
	return mapKeysMinLengthValidator{DisplayName: displayName, MinLen: minLen}
}

type mapKeysMinLengthValidator struct {
	DisplayName string
	MinLen      int
}

func (v mapKeysMinLengthValidator) Description(ctx context.Context) string {
	return "All " + v.DisplayName + " (map keys) must be at least " + strconv.Itoa(v.MinLen) + " character(s) long."
}

func (v mapKeysMinLengthValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v mapKeysMinLengthValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	for key := range req.ConfigValue.Elements() {
		if len(strings.TrimSpace(key)) < v.MinLen {
			resp.Diagnostics.AddAttributeError(
				req.Path.AtMapKey(key),
				"Invalid "+v.DisplayName,
				v.DisplayName+" (map key) must be at least "+strconv.Itoa(v.MinLen)+" character(s) long, without leading or trailing whitespace.",
			)
		}
	}
}
