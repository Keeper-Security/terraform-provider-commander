// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const (
	schedTypeDaily            = "DAILY"
	schedTypeWeekly           = "WEEKLY"
	schedTypeMonthlyByWeekday = "MONTHLY_BY_WEEKDAY"
	schedTypeYearly           = "YEARLY"
)

var (
	rotationScheduleWeekdays = map[string]struct{}{
		"SUNDAY": {}, "MONDAY": {}, "TUESDAY": {}, "WEDNESDAY": {},
		"THURSDAY": {}, "FRIDAY": {}, "SATURDAY": {},
	}
	rotationScheduleMonths = map[string]struct{}{
		"JANUARY": {}, "FEBRUARY": {}, "MARCH": {}, "APRIL": {}, "MAY": {}, "JUNE": {},
		"JULY": {}, "AUGUST": {}, "SEPTEMBER": {}, "OCTOBER": {}, "NOVEMBER": {}, "DECEMBER": {},
	}
	rotationScheduleOccurrences = map[string]struct{}{
		"FIRST": {}, "SECOND": {}, "THIRD": {}, "FOURTH": {}, "LAST": {},
	}
	rotationScheduleTimeHHMMSS = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d:[0-5]\d$`)
	rotationScheduleTimeHHMM   = regexp.MustCompile(`^([01]?\d|2[0-3]):[0-5]\d$`)
)

// RotationScheduleJSONValidator validates schedule_json for pam rotation edit.
type rotationScheduleJSONValidator struct{}

func RotationScheduleJSONValidator() rotationScheduleJSONValidator {
	return rotationScheduleJSONValidator{}
}

func (rotationScheduleJSONValidator) Description(_ context.Context) string {
	return "must be a valid Commander pam rotation schedule JSON object (DAILY, WEEKLY, MONTHLY_BY_WEEKDAY, or YEARLY)"
}

func (rotationScheduleJSONValidator) MarkdownDescription(_ context.Context) string {
	return "Must be a single JSON object accepted by Commander `pam rotation edit --schedulejson` (`DAILY`, `WEEKLY`, `MONTHLY_BY_WEEKDAY`, or `YEARLY`) with type-specific required fields. For cron schedules use `schedule_cron` instead."
}

func (rotationScheduleJSONValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if err := ValidateRotationScheduleJSON(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid schedule_json", err.Error())
	}
}

// ValidateRotationScheduleJSON parses and validates a Commander --schedulejson payload.
func ValidateRotationScheduleJSON(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("schedule_json cannot be empty; omit the attribute if you do not want a JSON schedule")
	}

	var sched map[string]interface{}
	if err := json.Unmarshal([]byte(s), &sched); err != nil {
		return fmt.Errorf("schedule_json must be a valid JSON object: %w", err)
	}
	if len(sched) == 0 {
		return fmt.Errorf("schedule_json must be a JSON object with a \"type\" field")
	}

	schedType, ok := jsonStringField(sched, "type")
	if !ok {
		return fmt.Errorf("schedule_json must include a non-empty \"type\" field")
	}
	schedType = strings.ToUpper(schedType)

	switch schedType {
	case schedTypeDaily:
		return validateDailySchedule(sched)
	case schedTypeWeekly:
		return validateWeeklySchedule(sched)
	case schedTypeMonthlyByWeekday:
		return validateMonthlyByWeekdaySchedule(sched)
	case schedTypeYearly:
		return validateYearlySchedule(sched)
	case "CRON":
		return fmt.Errorf("schedule_json does not support type CRON; use schedule_cron instead")
	case "MONTHLY_BY_DAY":
		return fmt.Errorf("schedule_json does not support type MONTHLY_BY_DAY")
	default:
		return fmt.Errorf(
			"schedule_json \"type\" must be one of DAILY, WEEKLY, MONTHLY_BY_WEEKDAY, or YEARLY; got %q",
			schedType,
		)
	}
}

func validateDailySchedule(sched map[string]interface{}) error {
	if err := forbidScheduleFields(sched, "weekday", "monthDay", "occurrence", "month", "cron"); err != nil {
		return err
	}
	return validateCommonScheduleFields(sched)
}

func validateWeeklySchedule(sched map[string]interface{}) error {
	if err := forbidScheduleFields(sched, "monthDay", "occurrence", "month", "cron"); err != nil {
		return err
	}
	if err := validateWeekdayField(sched); err != nil {
		return err
	}
	return validateCommonScheduleFields(sched)
}

func validateMonthlyByWeekdaySchedule(sched map[string]interface{}) error {
	if err := forbidScheduleFields(sched, "monthDay", "month", "cron"); err != nil {
		return err
	}
	if err := validateWeekdayField(sched); err != nil {
		return err
	}
	if err := validateOccurrenceField(sched); err != nil {
		return err
	}
	return validateCommonScheduleFields(sched)
}

func validateYearlySchedule(sched map[string]interface{}) error {
	if err := forbidScheduleFields(sched, "weekday", "occurrence", "cron"); err != nil {
		return err
	}
	month, ok := jsonStringField(sched, "month")
	if !ok {
		return fmt.Errorf("schedule_json with type YEARLY requires \"month\" (JANUARY..DECEMBER)")
	}
	if _, ok := rotationScheduleMonths[strings.ToUpper(month)]; !ok {
		return fmt.Errorf("schedule_json \"month\" must be one of JANUARY..DECEMBER; got %q", month)
	}
	if err := validateMonthDayField(sched); err != nil {
		return err
	}
	return validateCommonScheduleFields(sched)
}

func validateCommonScheduleFields(sched map[string]interface{}) error {
	if err := validateTimeFields(sched); err != nil {
		return err
	}
	return validateIntervalCountField(sched)
}

func validateTimeFields(sched map[string]interface{}) error {
	_, hasTime := jsonStringField(sched, "time")
	_, hasUTCTime := jsonStringField(sched, "utcTime")
	if !hasTime && !hasUTCTime {
		return fmt.Errorf("schedule_json requires either \"time\" (HH:MM:SS) or \"utcTime\" (HH:MM)")
	}
	if hasTime && hasUTCTime {
		return fmt.Errorf("schedule_json must use either \"time\" or \"utcTime\", not both")
	}
	if hasTime {
		timeVal, _ := jsonStringField(sched, "time")
		if !rotationScheduleTimeHHMMSS.MatchString(timeVal) {
			return fmt.Errorf("schedule_json \"time\" must be in HH:MM:SS 24-hour format; got %q", timeVal)
		}
	}
	if hasUTCTime {
		utcTime, _ := jsonStringField(sched, "utcTime")
		if !rotationScheduleTimeHHMM.MatchString(utcTime) {
			return fmt.Errorf("schedule_json \"utcTime\" must be in HH:MM 24-hour format; got %q", utcTime)
		}
	}
	return nil
}

func validateIntervalCountField(sched map[string]interface{}) error {
	if !hasJSONField(sched, "intervalCount") {
		return nil
	}
	n, ok := jsonIntField(sched, "intervalCount")
	if !ok || n < 1 {
		return fmt.Errorf("schedule_json \"intervalCount\" must be a positive integer")
	}
	return nil
}

func validateWeekdayField(sched map[string]interface{}) error {
	weekday, ok := jsonStringField(sched, "weekday")
	if !ok {
		return fmt.Errorf("schedule_json with type WEEKLY or MONTHLY_BY_WEEKDAY requires \"weekday\" (SUNDAY..SATURDAY)")
	}
	if _, ok := rotationScheduleWeekdays[strings.ToUpper(weekday)]; !ok {
		return fmt.Errorf("schedule_json \"weekday\" must be one of SUNDAY..SATURDAY; got %q", weekday)
	}
	return nil
}

func validateOccurrenceField(sched map[string]interface{}) error {
	occurrence, ok := jsonStringField(sched, "occurrence")
	if !ok {
		return fmt.Errorf("schedule_json with type MONTHLY_BY_WEEKDAY requires \"occurrence\" (FIRST, SECOND, THIRD, FOURTH, or LAST)")
	}
	if _, ok := rotationScheduleOccurrences[strings.ToUpper(occurrence)]; !ok {
		return fmt.Errorf("schedule_json \"occurrence\" must be one of FIRST, SECOND, THIRD, FOURTH, or LAST; got %q", occurrence)
	}
	return nil
}

func validateMonthDayField(sched map[string]interface{}) error {
	if !hasJSONField(sched, "monthDay") {
		return fmt.Errorf("schedule_json requires \"monthDay\" (1-28)")
	}
	n, ok := jsonIntField(sched, "monthDay")
	if !ok || n < 1 || n > 28 {
		return fmt.Errorf("schedule_json \"monthDay\" must be an integer between 1 and 28")
	}
	return nil
}

func forbidScheduleFields(sched map[string]interface{}, fields ...string) error {
	for _, field := range fields {
		if hasJSONField(sched, field) {
			return fmt.Errorf("schedule_json must not include %q for this schedule type", field)
		}
	}
	return nil
}

func jsonStringField(m map[string]interface{}, key string) (string, bool) {
	v, ok := m[key]
	if !ok || v == nil {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	return s, true
}

func jsonIntField(m map[string]interface{}, key string) (int, bool) {
	v, ok := m[key]
	if !ok || v == nil {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		if n != math.Trunc(n) {
			return 0, false
		}
		return int(n), true
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			return 0, false
		}
		return int(i), true
	default:
		return 0, false
	}
}

func hasJSONField(m map[string]interface{}, key string) bool {
	_, ok := m[key]
	return ok
}
