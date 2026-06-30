// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package cronvalidate

import (
	"fmt"
	"strconv"
	"strings"
)

// validateMinimumInterval enforces Keeper's rule that rotation cron schedules
// must not run more frequently than once per hour. See:
// https://docs.keeper.io/keeperpam/privileged-access-manager/references/cron-spec
func validateMinimumInterval(seconds, minutes string) error {
	if !isSingleNumericValue(seconds) {
		return fmt.Errorf("seconds field must be a single value between 0 and 59; rotation schedules require at least a 1-hour interval")
	}

	minutes = strings.TrimSpace(minutes)
	if minutes == "*" {
		return fmt.Errorf("minutes field cannot be '*' because rotation schedules require at least a 1-hour interval")
	}

	for _, part := range splitList(minutes) {
		part = strings.TrimSpace(part)
		if part == "*" {
			return fmt.Errorf("minutes field cannot be '*' because rotation schedules require at least a 1-hour interval")
		}
		if step, ok := incrementStep(part); ok && step < 60 {
			return fmt.Errorf("minute increments shorter than 60 are not supported; rotation schedules require at least a 1-hour interval")
		}
	}
	return nil
}

func isSingleNumericValue(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	n, err := strconv.Atoi(s)
	return err == nil && n >= 0
}

func incrementStep(part string) (int, bool) {
	if !strings.Contains(part, "/") {
		return 0, false
	}
	left, stepStr, ok := strings.Cut(part, "/")
	if !ok {
		return 0, false
	}
	step, err := strconv.Atoi(strings.TrimSpace(stepStr))
	if err != nil || step <= 0 {
		return 0, false
	}
	if strings.TrimSpace(left) == "*" {
		return step, true
	}
	return step, true
}
