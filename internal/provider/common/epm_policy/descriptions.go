// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import "strings"

// PolicyTypeDescription returns a comma-separated list of allowed policy types for plain-text descriptions.
func PolicyTypeDescription() string {
	return strings.Join(PolicyTypeValues, ", ")
}

// PolicyTypeMarkdown returns a markdown-formatted list of allowed policy types (backtick-wrapped).
func PolicyTypeMarkdown() string {
	return "`" + strings.Join(PolicyTypeValues, "`, `") + "`"
}

// StatusDescription returns a comma-separated list of allowed status values for plain-text descriptions.
func StatusDescription() string {
	return strings.Join(StatusValues, ", ")
}

// StatusMarkdown returns a markdown-formatted list of allowed status values (backtick-wrapped).
func StatusMarkdown() string {
	return "`" + strings.Join(StatusValues, "`, `") + "`"
}

// StatusDescriptionForLeastPrivilege returns the allowed statuses for least_privilege policy type.
func StatusDescriptionForLeastPrivilege() string {
	return strings.Join(AllowedStatusValuesForLeastPrivilege, ", ")
}

// StatusMarkdownForLeastPrivilege returns markdown-formatted allowed statuses for least_privilege.
func StatusMarkdownForLeastPrivilege() string {
	return "`" + strings.Join(AllowedStatusValuesForLeastPrivilege, "`, `") + "`"
}

// ControlDescription returns a comma-separated list of allowed control values for plain-text descriptions.
func ControlDescription() string {
	return strings.Join(ControlValues, ", ")
}

// ControlMarkdown returns a markdown-formatted list of allowed control values (backtick-wrapped).
func ControlMarkdown() string {
	return "`" + strings.Join(ControlValues, "`, `") + "`"
}

// DayFilterDescription returns a comma-separated list of day names for plain-text descriptions.
func DayFilterDescription() string {
	return strings.Join(DayFilterValues, ", ")
}

// DayFilterMarkdown returns a markdown-formatted list of day names for descriptions.
func DayFilterMarkdown() string {
	return strings.Join(DayFilterValues, ", ")
}
