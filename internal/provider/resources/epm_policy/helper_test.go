// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package epmpolicy

import (
	"strings"
	"testing"

	commonepm "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/epm_policy"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func epmModelForCommand(status string) *EpmPolicyResourceModel {
	emptySet := types.SetNull(types.StringType)
	return &EpmPolicyResourceModel{
		PolicyName:         types.StringValue("p1"),
		PolicyType:         types.StringValue("command"),
		Status:             types.StringValue(status),
		Control:            emptySet,
		UserGroups:         emptySet,
		MachineCollections: emptySet,
		Applications:       emptySet,
		DayFilter:          emptySet,
		TimeFilter:         emptySet,
		DateFilter:         emptySet,
	}
}

func TestBuildCreateCommand_StatusOffUsesEnableFlag(t *testing.T) {
	cmd := buildCreateCommand(epmModelForCommand(commonepm.StatusOff))
	if !strings.Contains(cmd, commonepm.FlagEnable+" '"+commonepm.StatusOff+"'") {
		t.Fatalf("create with status off should use enable flag, got: %s", cmd)
	}
	if strings.Contains(cmd, commonepm.FlagStatus) {
		t.Fatalf("create with status off should not use status flag, got: %s", cmd)
	}
}

func TestBuildCreateCommand_NonOffStatusUsesStatusFlag(t *testing.T) {
	cmd := buildCreateCommand(epmModelForCommand(commonepm.StatusEnforce))
	if !strings.Contains(cmd, commonepm.FlagStatus+" '"+commonepm.StatusEnforce+"'") {
		t.Fatalf("create with enforce should use status flag, got: %s", cmd)
	}
	if strings.Contains(cmd, commonepm.FlagEnable) {
		t.Fatalf("create with enforce should not use enable flag, got: %s", cmd)
	}
}

func TestBuildUpdateCommand_StatusOffUsesEnableFlag(t *testing.T) {
	cmd := buildUpdateCommand("42", epmModelForCommand(commonepm.StatusOff))
	if !strings.Contains(cmd, commonepm.FlagEnable+" '"+commonepm.StatusOff+"'") {
		t.Fatalf("update with status off should use enable flag, got: %s", cmd)
	}
}

func TestBuildUpdateCommand_NonOffStatusUsesStatusFlag(t *testing.T) {
	cmd := buildUpdateCommand("42", epmModelForCommand(commonepm.StatusMonitor))
	if !strings.Contains(cmd, commonepm.FlagStatus+" '"+commonepm.StatusMonitor+"'") {
		t.Fatalf("update with monitor should use status flag, got: %s", cmd)
	}
}
