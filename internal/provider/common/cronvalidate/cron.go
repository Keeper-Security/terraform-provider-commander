// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

// Package cronvalidate provides Keeper/Commander-style rotation cron parsing and
// Terraform string validators for reuse across resources (e.g. pam_user schedule_cron).
package cronvalidate

import (
	"fmt"
	"strings"

	"github.com/robfig/cron/v3"
)

var (
	// Quartz-style 6-field (seconds first) plus descriptors such as @daily.
	parserQuartzOrDescriptor = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	// Standard 5-field (minutes first).
	parserFiveField = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
)

// ValidateKeeperRotationCron returns nil if s is a non-empty cron expression
// accepted by the same rules used for Commander rotation schedules: 6-field
// Quartz (e.g. "0 0 3 1 * ?"), 5-field (e.g. "56 17 * * *"), or a descriptor
// preset such as @daily when supported by the underlying parser.
func ValidateKeeperRotationCron(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("cron expression cannot be empty")
	}
	if _, err := parserQuartzOrDescriptor.Parse(s); err == nil {
		return nil
	}
	if _, err := parserFiveField.Parse(s); err == nil {
		return nil
	}
	_, err6 := parserQuartzOrDescriptor.Parse(s)
	_, err5 := parserFiveField.Parse(s)
	return fmt.Errorf("not a valid cron expression (examples: 6-field Quartz \"0 0 3 1 * ?\", 5-field \"56 17 * * *\"): 6-field: %v; 5-field: %v", err6, err5)
}

// AttributeLabel returns name trimmed, or a generic label if name is empty.
func AttributeLabel(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "cron expression"
	}
	return name
}
