// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package cronvalidate provides Keeper PAM rotation cron parsing and Terraform
// string validators aligned with the Keeper Quartz cron specification:
// https://docs.keeper.io/keeperpam/privileged-access-manager/references/cron-spec
package cronvalidate

import (
	"fmt"
	"strings"
)

const (
	fieldCountSix   = 6
	fieldCountSeven = 7
)

// ValidateKeeperRotationCron returns nil if s is a non-empty 6- or 7-field
// Quartz-style cron expression accepted for Keeper PAM rotation schedules.
func ValidateKeeperRotationCron(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("cron expression cannot be empty")
	}

	fields := strings.Fields(s)
	switch len(fields) {
	case fieldCountSix:
		return validateQuartzFields(fields)
	case fieldCountSeven:
		if err := validateQuartzFields(fields[:fieldCountSix]); err != nil {
			return err
		}
		return validateYearField(fields[6])
	default:
		return fmt.Errorf(
			"cron expression must have 6 or 7 whitespace-separated fields (seconds minutes hours day-of-month month day-of-week [year]); got %d fields",
			len(fields),
		)
	}
}

func validateQuartzFields(fields []string) error {
	seconds, minutes, hours, dom, month, dow := fields[0], fields[1], fields[2], fields[3], fields[4], fields[5]

	checks := []struct {
		name     string
		value    string
		allowed  string
		validate func(string) error
	}{
		{"seconds", seconds, "0123456789*,-/", validateSecondsField},
		{"minutes", minutes, "0123456789*,-/", validateMinutesField},
		{"hours", hours, "0123456789*,-/", validateHoursField},
		{"day of month", dom, "0123456789*,-/?LW", validateDayOfMonthField},
		{"month", month, "0123456789*,-/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", validateMonthField},
		{"day of week", dow, "0123456789*,-/?L#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", validateDayOfWeekField},
	}

	for _, c := range checks {
		if err := validateFieldCharacters(c.name, c.value, c.allowed); err != nil {
			return err
		}
		if err := c.validate(c.value); err != nil {
			return err
		}
	}

	if err := validateDOMDOWCombination(dom, dow); err != nil {
		return err
	}

	return validateMinimumInterval(seconds, minutes)
}

// Quartz requires that day-of-month and day-of-week are not both specified
// unless one of them is '?'.
func validateDOMDOWCombination(dom, dow string) error {
	domSpecified := !strings.Contains(dom, "?")
	dowSpecified := !strings.Contains(dow, "?")
	if domSpecified && dowSpecified {
		return fmt.Errorf("day-of-month and day-of-week cannot both be specified; set one field to '?' (for example: \"0 28 17 ? * *\" or \"0 15 10 15 * ?\")")
	}
	if !domSpecified && !dowSpecified {
		return fmt.Errorf("day-of-month and day-of-week cannot both be '?'; specify one scheduling dimension and set the other to '?'")
	}
	return nil
}

// AttributeLabel returns name trimmed, or a generic label if name is empty.
func AttributeLabel(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "cron expression"
	}
	return name
}
