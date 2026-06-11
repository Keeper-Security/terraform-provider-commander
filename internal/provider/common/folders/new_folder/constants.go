// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package new_folder holds types and API helpers for Keeper Drive (new) folder resources.
package new_folder

// AccessTypeApplication is the access_type value for Keeper-managed
// application principals. These entries appear in nsf-get user_permissions /
// team_permissions but are not user-controllable, so they are filtered out
// of the Terraform share map. Matched strictly (case-sensitive) against the
// API value.
const AccessTypeApplication = "AT_APPLICATION"

// ExecuteCommand operation descriptions (logs / error context).
const (
	ErrOpCreate = "Unable to create Keeper Drive folder"
	ErrOpGet    = "Unable to get Keeper Drive folder"
	ErrOpDelete = "Unable to delete Keeper Drive folder"
)
