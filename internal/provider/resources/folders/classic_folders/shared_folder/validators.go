// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sharedfolder

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Expiration format: "never" | yyyy-MM-ddTHH:mm:ss (e.g. 2026-04-02T11:11:00).
var expirationNever = regexp.MustCompile(`(?i)^never$`)

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
	if _, err := time.Parse(TimeLayoutExpiration, s); err == nil {
		return
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

// userExpirationManageUsersValidator rejects manage_users = true when expiration is a datetime (not "never").
type userExpirationManageUsersValidator struct{}

func (userExpirationManageUsersValidator) Description(ctx context.Context) string {
	return "manage_users cannot be true when expiration is a datetime; use \"never\" for users who manage other users."
}

func (userExpirationManageUsersValidator) MarkdownDescription(ctx context.Context) string {
	return "`manage_users` cannot be `true` when `expiration` is a datetime; use `\"never\"` for users who manage other users."
}

// expirationIsTimeLimited is true only when expiration is a valid yyyy-MM-ddTHH:mm:ss.
// It is false when expiration is omitted (null/unknown), empty, "never", or invalid (leave that to ExpirationValidator).
func expirationIsTimeLimited(expVal types.String) bool {
	if expVal.IsNull() || expVal.IsUnknown() {
		return false
	}
	s := strings.TrimSpace(expVal.ValueString())
	if s == "" || expirationNever.MatchString(s) {
		return false
	}
	_, err := time.Parse(TimeLayoutExpiration, s)
	return err == nil
}

func (userExpirationManageUsersValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	attrs := req.ConfigValue.Attributes()

	expVal, expOk := attrs[AttrExpiration].(types.String)
	if !expOk || !expirationIsTimeLimited(expVal) {
		return // no error when expiration not set, "never", empty, or invalid
	}

	muVal, muOk := attrs[AttrManageUsers].(types.Bool)
	if !muOk || muVal.IsNull() || muVal.IsUnknown() || !muVal.ValueBool() {
		return
	}

	resp.Diagnostics.AddAttributeError(
		req.Path.AtName(AttrManageUsers),
		ErrMsgInvalidUserPermissionsCombo,
		ErrMsgManageUsersWithTimeLimitedExpiration,
	)
}

// UserExpirationManageUsersValidator returns a validator for each users map entry object.
func UserExpirationManageUsersValidator() userExpirationManageUsersValidator {
	return userExpirationManageUsersValidator{}
}
