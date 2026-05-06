// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"encoding/json"
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

// ----- GENERIC: INT32 NON-NEGATIVE --------------------------------
// Int32NonNegativeValidator validates that an int32 value is >= 0.
// DisplayName is used in error messages (e.g. "Remote Target Port").
// AllowNull: if true, null values are skipped (for optional attributes); if false, null is an error.
func Int32NonNegativeValidator(displayName string, allowNull bool) int32NonNegativeValidator {
	return int32NonNegativeValidator{DisplayName: displayName, AllowNull: allowNull}
}

type int32NonNegativeValidator struct {
	DisplayName string
	AllowNull   bool
}

func (v int32NonNegativeValidator) Description(_ context.Context) string {
	return v.DisplayName + " must be a non-negative integer (>= 0)."
}

func (v int32NonNegativeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v int32NonNegativeValidator) ValidateInt32(ctx context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	if req.ConfigValue.IsUnknown() {
		return
	}
	if v.AllowNull && req.ConfigValue.IsNull() {
		return
	}
	if req.ConfigValue.ValueInt32() < 0 {
		resp.Diagnostics.AddError(
			"Invalid "+v.DisplayName,
			fmt.Sprintf("%s must be a non-negative integer (>= 0), got: %d.", v.DisplayName, req.ConfigValue.ValueInt32()),
		)
	}
}

// ----- GENERIC: STRING ONE-OF --------------------------------
// StringOneOfValidator validates that a string value is one of the allowed values.
// DisplayName is used in error messages. AllowNull: if true, null values are skipped.
func StringOneOfValidator(displayName string, allowed []string, allowNull bool) stringOneOfValidator {
	return stringOneOfValidator{DisplayName: displayName, Allowed: allowed, AllowNull: allowNull}
}

type stringOneOfValidator struct {
	DisplayName string
	Allowed     []string
	AllowNull   bool
}

func (v stringOneOfValidator) Description(_ context.Context) string {
	return v.DisplayName + " must be one of: " + strings.Join(v.Allowed, ", ") + "."
}

func (v stringOneOfValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v stringOneOfValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() {
		return
	}
	if v.AllowNull && req.ConfigValue.IsNull() {
		return
	}
	val := req.ConfigValue.ValueString()
	for _, a := range v.Allowed {
		if val == a {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid "+v.DisplayName,
		fmt.Sprintf("%s %q is not supported. Must be one of: %s.", v.DisplayName, val, strings.Join(v.Allowed, ", ")),
	)
}

// ----- GENERIC: INT32 ONE-OF --------------------------------
// Int32OneOfValidator validates that an int32 value is one of the allowed values.
func Int32OneOfValidator(displayName string, allowed []int32, allowNull bool) int32OneOfValidator {
	return int32OneOfValidator{DisplayName: displayName, Allowed: allowed, AllowNull: allowNull}
}

type int32OneOfValidator struct {
	DisplayName string
	Allowed     []int32
	AllowNull   bool
}

func (v int32OneOfValidator) Description(_ context.Context) string {
	vals := make([]string, len(v.Allowed))
	for i, a := range v.Allowed {
		vals[i] = strconv.FormatInt(int64(a), 10)
	}
	return v.DisplayName + " must be one of: " + strings.Join(vals, ", ") + "."
}

func (v int32OneOfValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v int32OneOfValidator) ValidateInt32(_ context.Context, req validator.Int32Request, resp *validator.Int32Response) {
	if req.ConfigValue.IsUnknown() {
		return
	}
	if v.AllowNull && req.ConfigValue.IsNull() {
		return
	}
	val := req.ConfigValue.ValueInt32()
	for _, a := range v.Allowed {
		if val == a {
			return
		}
	}
	vals := make([]string, len(v.Allowed))
	for i, a := range v.Allowed {
		vals[i] = strconv.FormatInt(int64(a), 10)
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid "+v.DisplayName,
		fmt.Sprintf("%s %d is not supported. Must be one of: %s.", v.DisplayName, val, strings.Join(vals, ", ")),
	)
}

// ----- GENERIC: JSON STRING --------------------------------
// JSONStringValidator validates that a string value is valid JSON.
// DisplayName is used in error messages (e.g. "Service Account Key").
func JSONStringValidator(displayName string) jsonStringValidator {
	return jsonStringValidator{DisplayName: displayName}
}

type jsonStringValidator struct {
	DisplayName string
}

func (v jsonStringValidator) Description(_ context.Context) string {
	return v.DisplayName + " must be a valid JSON string."
}

func (v jsonStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v jsonStringValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	val := req.ConfigValue.ValueString()
	if !json.Valid([]byte(val)) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid "+v.DisplayName,
			v.DisplayName+" must be a valid JSON string.",
		)
	}
}

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
