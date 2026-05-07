// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package records

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBuildRecordAdd_contactMinimal(t *testing.T) {
	folder := types.StringValue("/My Folder")
	title := "My Title"
	extra := []string{
		formatFieldAssignment("email", "a@b.co", false),
	}
	cmd := BuildRecordAdd(folder, title, RecordTypeContact, extra, nil, types.StringNull())
	if !strings.Contains(cmd, "record-add") {
		t.Fatalf("expected record-add, got: %s", cmd)
	}
	if !strings.Contains(cmd, "--folder") || !strings.Contains(cmd, "My Folder") {
		t.Fatalf("expected folder flag: %s", cmd)
	}
	if !strings.Contains(cmd, "--title") || !strings.Contains(cmd, "My Title") {
		t.Fatalf("expected title: %s", cmd)
	}
	if !strings.Contains(cmd, "--record-type contact") {
		t.Fatalf("expected record type: %s", cmd)
	}
	if !strings.Contains(cmd, "email=") {
		t.Fatalf("expected email field: %s", cmd)
	}
}

func TestQuoteShellSingle_apostrophe(t *testing.T) {
	got := QuoteShellSingle("O'Brien")
	if !strings.Contains(got, "'") {
		t.Fatal("expected quoted string")
	}
}
