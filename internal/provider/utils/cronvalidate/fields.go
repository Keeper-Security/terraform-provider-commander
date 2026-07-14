// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package cronvalidate

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	monthNames = map[string]int{
		"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4, "MAY": 5, "JUN": 6,
		"JUL": 7, "AUG": 8, "SEP": 9, "OCT": 10, "NOV": 11, "DEC": 12,
	}
	dowNames = map[string]int{
		"SUN": 1, "MON": 2, "TUE": 3, "WED": 4, "THU": 5, "FRI": 6, "SAT": 7,
	}
)

func validateSecondsField(s string) error {
	return validateNumericField("seconds", s, 0, 59)
}

func validateMinutesField(s string) error {
	return validateNumericField("minutes", s, 0, 59)
}

func validateHoursField(s string) error {
	return validateNumericField("hours", s, 0, 23)
}

func validateDayOfMonthField(s string) error {
	for _, part := range splitList(s) {
		if err := validateDayOfMonthPart(part); err != nil {
			return err
		}
	}
	return nil
}

func validateMonthField(s string) error {
	for _, part := range splitList(s) {
		if err := validateMonthPart(part); err != nil {
			return err
		}
	}
	return nil
}

func validateDayOfWeekField(s string) error {
	for _, part := range splitList(s) {
		if err := validateDayOfWeekPart(part); err != nil {
			return err
		}
	}
	return nil
}

func validateYearField(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("year field cannot be empty when present")
	}
	return validateNumericField("year", s, 1970, 2099)
}

func validateDayOfMonthPart(part string) error {
	part = strings.ToUpper(strings.TrimSpace(part))
	switch {
	case part == "*", part == "?":
		return nil
	case part == "L", part == "LW":
		return nil
	case strings.HasPrefix(part, "L-"):
		return validateSingleNumber("day of month", strings.TrimPrefix(part, "L-"), 1, 31)
	case strings.HasSuffix(part, "W"):
		return validateSingleNumber("day of month", strings.TrimSuffix(part, "W"), 1, 31)
	default:
		return validateNumericPart("day of month", part, 1, 31)
	}
}

func validateMonthPart(part string) error {
	part = strings.ToUpper(strings.TrimSpace(part))
	if part == "*" {
		return nil
	}
	if strings.Contains(part, "/") {
		return validateIncrementPart("month", part, 1, 12, parseMonthToken)
	}
	if strings.Contains(part, "-") {
		return validateRangePart("month", part, 1, 12, parseMonthToken)
	}
	if _, err := parseMonthToken(part); err != nil {
		return fmt.Errorf("invalid month field value %q: %w", part, err)
	}
	return nil
}

func validateDayOfWeekPart(part string) error {
	part = strings.ToUpper(strings.TrimSpace(part))
	switch {
	case part == "*", part == "?":
		return nil
	case part == "L":
		return nil
	case strings.Contains(part, "#"):
		return validateDOWHashPart(part)
	case strings.HasSuffix(part, "L"):
		return validateDOWLastPart(part)
	default:
		if strings.Contains(part, "/") {
			return validateIncrementPart("day of week", part, 1, 7, parseDOWToken)
		}
		if strings.Contains(part, "-") {
			return validateRangePart("day of week", part, 1, 7, parseDOWToken)
		}
		if _, err := parseDOWToken(part); err != nil {
			return fmt.Errorf("invalid day of week field value %q: %w", part, err)
		}
		return nil
	}
}

func validateDOWHashPart(part string) error {
	before, after, ok := strings.Cut(part, "#")
	if !ok || after == "" {
		return fmt.Errorf("invalid day of week field value %q: expected format like 6#3", part)
	}
	if _, err := parseDOWToken(before); err != nil {
		return fmt.Errorf("invalid day of week field value %q: %w", part, err)
	}
	n, err := strconv.Atoi(after)
	if err != nil || n < 1 || n > 5 {
		return fmt.Errorf("invalid day of week field value %q: # suffix must be 1-5", part)
	}
	return nil
}

func validateDOWLastPart(part string) error {
	token := strings.TrimSuffix(part, "L")
	if token == "" {
		return fmt.Errorf("invalid day of week field value %q", part)
	}
	if _, err := parseDOWToken(token); err != nil {
		return fmt.Errorf("invalid day of week field value %q: %w", part, err)
	}
	return nil
}

func validateNumericField(name, s string, minBound, maxBound int) error {
	for _, part := range splitList(s) {
		if err := validateNumericPart(name, part, minBound, maxBound); err != nil {
			return err
		}
	}
	return nil
}

func validateNumericPart(name, part string, minBound, maxBound int) error {
	part = strings.TrimSpace(part)
	if part == "*" || part == "?" {
		return nil
	}
	if strings.Contains(part, "/") {
		return validateIncrementPart(name, part, minBound, maxBound, parseIntToken(minBound, maxBound))
	}
	if strings.Contains(part, "-") {
		return validateRangePart(name, part, minBound, maxBound, parseIntToken(minBound, maxBound))
	}
	if _, err := parseIntToken(minBound, maxBound)(part); err != nil {
		return fmt.Errorf("invalid %s field value %q: %w", name, part, err)
	}
	return nil
}

type tokenParser func(string) (int, error)

func parseIntToken(minBound, maxBound int) tokenParser {
	return func(s string) (int, error) {
		return parseBoundedInt(s, minBound, maxBound)
	}
}

func parseMonthToken(s string) (int, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if n, ok := monthNames[s]; ok {
		return n, nil
	}
	return parseBoundedInt(s, 1, 12)
}

func parseDOWToken(s string) (int, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if n, ok := dowNames[s]; ok {
		return n, nil
	}
	return parseBoundedInt(s, 1, 7)
}

func parseBoundedInt(s string, minBound, maxBound int) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("value cannot be empty")
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return 0, fmt.Errorf("expected a number between %d and %d, got %q", minBound, maxBound, s)
		}
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("expected a number between %d and %d, got %q", minBound, maxBound, s)
	}
	if n < minBound || n > maxBound {
		return 0, fmt.Errorf("expected a number between %d and %d, got %d", minBound, maxBound, n)
	}
	return n, nil
}

func validateSingleNumber(name, part string, minBound, maxBound int) error {
	if _, err := parseBoundedInt(part, minBound, maxBound); err != nil {
		return fmt.Errorf("invalid %s field value %q: %w", name, part, err)
	}
	return nil
}

func validateIncrementPart(name, part string, minBound, maxBound int, parse tokenParser) error {
	left, stepStr, ok := strings.Cut(part, "/")
	if !ok || stepStr == "" {
		return fmt.Errorf("invalid %s field value %q: increment step is required after '/'", name, part)
	}
	step, err := parseBoundedInt(stepStr, 1, maxBound-minBound+1)
	if err != nil {
		return fmt.Errorf("invalid %s field increment step in %q: %w", name, part, err)
	}
	if step <= 0 {
		return fmt.Errorf("invalid %s field value %q: increment step must be positive", name, part)
	}
	left = strings.TrimSpace(left)
	if left == "" {
		return fmt.Errorf("invalid %s field value %q: missing value before '/'", name, part)
	}
	if left == "*" {
		return nil
	}
	if strings.Contains(left, "-") {
		return validateRangePart(name, left, minBound, maxBound, parse)
	}
	if _, err := parse(left); err != nil {
		return fmt.Errorf("invalid %s field value %q: %w", name, part, err)
	}
	return nil
}

func validateRangePart(name, part string, minBound, maxBound int, parse tokenParser) error {
	startStr, endStr, ok := strings.Cut(part, "-")
	if !ok || endStr == "" {
		return fmt.Errorf("invalid %s field value %q: expected start-end range", name, part)
	}
	start, err := parse(startStr)
	if err != nil {
		return fmt.Errorf("invalid %s field range start in %q: %w", name, part, err)
	}
	end, err := parse(endStr)
	if err != nil {
		return fmt.Errorf("invalid %s field range end in %q: %w", name, part, err)
	}
	if start > end {
		return fmt.Errorf("invalid %s field value %q: range start must not exceed end", name, part)
	}
	if start < minBound || end > maxBound {
		return fmt.Errorf("invalid %s field value %q: range must be within %d-%d", name, part, minBound, maxBound)
	}
	return nil
}

func splitList(s string) []string {
	parts := strings.FieldsFunc(strings.TrimSpace(s), func(r rune) bool { return r == ',' })
	if len(parts) == 0 {
		return []string{""}
	}
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{""}
	}
	return out
}

func validateFieldCharacters(fieldName, s string, allowed string) error {
	for _, r := range s {
		if !strings.ContainsRune(allowed, r) {
			return fmt.Errorf("%s field contains invalid character %q", fieldName, string(r))
		}
	}
	return nil
}
