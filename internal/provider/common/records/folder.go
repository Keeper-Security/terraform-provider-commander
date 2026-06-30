// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package records

import (
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ExtractFolderValue resolves the folder value from the API response against the
// current Terraform state. If the user originally specified a UID or path that
// matches the API response, the state value is preserved so Terraform does not
// show a spurious diff. Otherwise the folder UID from the response is used.
// Path comparison normalizes spaces around "/" so "Test / My Folder" matches "Test/My Folder".
func ExtractFolderValue(folder *utils.FolderLocationResponse, stateFolder types.String) types.String {
	if folder == nil || (strings.TrimSpace(folder.UID) == "" && strings.TrimSpace(folder.Path) == "") {
		return types.StringNull()
	}
	if !stateFolder.IsNull() && !stateFolder.IsUnknown() {
		sv := strings.TrimSpace(stateFolder.ValueString())
		if sv == strings.TrimSpace(folder.UID) || normalizeFolderPath(sv) == normalizeFolderPath(folder.Path) {
			return stateFolder
		}
	}
	return types.StringValue(strings.TrimSpace(folder.UID))
}

// normalizeFolderPath removes spaces around "/" separators so that
// "Test / My Folder" and "Test/My Folder" compare as equal.
func normalizeFolderPath(p string) string {
	parts := strings.Split(p, "/")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return strings.Join(parts, "/")
}
