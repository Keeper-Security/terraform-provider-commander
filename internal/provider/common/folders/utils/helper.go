// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"fmt"
	"strings"
)

// BuildFolderPath constructs the full vault folder path from a leaf name and
// optional parent folder_location. If folderLocation is non-empty the result is
// "folderLocation/name"; otherwise just "name". Shared by non_shared_folder,
// shared_folder and new_folder resources.
func BuildFolderPath(name, folderLocation string) string {
	name = strings.TrimSpace(name)
	folderLocation = strings.TrimSpace(folderLocation)
	if folderLocation == "" {
		return name
	}
	return folderLocation + "/" + name
}

// EscapeDoubleQuotesForCLI escapes double quotes for use inside double-quoted
// shell arguments passed to Commander CLI commands.
func EscapeDoubleQuotesForCLI(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

// SplitFolderPath splits a full vault path into parent and leaf name.
// Example: "Parent/MyFolder" -> parent "Parent", leaf "MyFolder".
// A path with no "/" returns empty parent and the whole string as leaf.
func SplitFolderPath(full string) (parent, leaf string) {
	full = strings.TrimSpace(full)
	if full == "" {
		return "", ""
	}
	i := strings.LastIndex(full, "/")
	if i < 0 {
		return "", full
	}
	return strings.TrimSpace(full[:i]), strings.TrimSpace(full[i+1:])
}

// MvPathForCommander normalizes a vault path for Commander `mv`. Paths with no parent
// (no `/` — at vault root) are prefixed with `/` so the CLI targets the root folder.
func MvPathForCommander(full string) string {
	full = strings.TrimSpace(full)
	if full == "" {
		return full
	}
	if strings.HasPrefix(full, "/") {
		return full
	}
	parent, leaf := SplitFolderPath(full)
	if parent == "" {
		return "/" + leaf
	}
	return full
}

// MvMoveTargetParent returns the destination parent folder for Commander `mv`.
// Example: "Templates/test4/MyFolder" -> "Templates/test4".
// "MyFolder" (vault root) -> "/".
func MvMoveTargetParent(planPath string) string {
	planPath = strings.TrimSpace(planPath)
	if planPath == "" {
		return planPath
	}
	trim := planPath
	if strings.HasPrefix(trim, "/") {
		trim = strings.TrimSpace(trim[1:])
	}
	parent, _ := SplitFolderPath(trim)
	parent = strings.TrimSpace(parent)
	if parent == "" {
		return "/"
	}
	return parent
}

// ExtractFolderUIDFromCreateResponse pulls KeyFolderUID ("folder_uid") out of
// the API response Data (a JSON-decoded map) returned by Commander folder
// create commands (mkdir/mkdir --shared-folder). If folder_uid is missing,
// empty, or not a string (e.g. a number), an error or its stringified form is
// returned accordingly. Shared by non_shared_folder and classic shared_folder
// resource Create flows.
func ExtractFolderUIDFromCreateResponse(data any) (string, error) {
	m, _ := data.(map[string]interface{})
	v, ok := m[KeyFolderUID]
	if !ok || v == nil {
		return "", fmt.Errorf("API response missing %s", KeyFolderUID)
	}
	if s, ok := v.(string); ok {
		if s == "" {
			return "", fmt.Errorf("API response %s is empty", KeyFolderUID)
		}
		return s, nil
	}
	return fmt.Sprintf("%v", v), nil
}
