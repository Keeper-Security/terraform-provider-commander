// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package sharedfolder_test

import (
	"testing"

	sharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/classic_folders/shared_folder"
)

func TestSplitSharedFolderPath(t *testing.T) {
	tests := []struct {
		full, parent, leaf string
	}{
		{"Templates/My Shared Folder 1", "Templates", "My Shared Folder 1"},
		{"My Shared Folder", "", "My Shared Folder"},
		{"", "", ""},
		{"  A/B  ", "A", "B"},
	}
	for _, tc := range tests {
		p, l := sharedfolder.SplitSharedFolderPath(tc.full)
		if p != tc.parent || l != tc.leaf {
			t.Errorf("SplitSharedFolderPath(%q) = (%q,%q), want (%q,%q)", tc.full, p, l, tc.parent, tc.leaf)
		}
	}
}

func TestMvPathForCommander(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"My Folder", "/My Folder"},
		{"/My Folder", "/My Folder"},
		{"Templates/x", "Templates/x"},
		{"", ""},
	}
	for _, tc := range tests {
		got := sharedfolder.MvPathForCommander(tc.in)
		if got != tc.want {
			t.Errorf("MvPathForCommander(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestMvMoveTargetParent(t *testing.T) {
	tests := []struct {
		planPath, want string
	}{
		{"Templates/test4/My Shared Folder 1", "Templates/test4"},
		{"My Shared Folder 1", "/"},
		{"/Leaf Only At Root", "/"},
		{"", ""},
	}
	for _, tc := range tests {
		got := sharedfolder.MvMoveTargetParent(tc.planPath)
		if got != tc.want {
			t.Errorf("MvMoveTargetParent(%q) = %q, want %q", tc.planPath, got, tc.want)
		}
	}
}
