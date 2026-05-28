// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package non_shared_folder

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ErrNonSharedFolderNotFound is returned when the get command yields no usable folder (deleted, wrong id/path, or empty response).
var ErrNonSharedFolderNotFound = errors.New("non-shared folder not found")

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

// EscapeDoubleQuotesForCLI escapes double quotes for use inside double-quoted shell arguments.
func EscapeDoubleQuotesForCLI(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

// BuildFolderPath constructs the full folder path from name and optional folder_location.
// If folderLocation is non-empty, the result is "folderLocation/name"; otherwise just "name".
func BuildFolderPath(name, folderLocation string) string {
	name = strings.TrimSpace(name)
	folderLocation = strings.TrimSpace(folderLocation)
	if folderLocation == "" {
		return name
	}
	return folderLocation + "/" + name
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

// BuildGetFolderCommand builds the Commander CLI: get '<id-or-path>' --format json.
func BuildGetFolderCommand(idOrPath string) string {
	return fmt.Sprintf("%s '%s' %s", utils.CmdGet, idOrPath, utils.FlagFormatJSON)
}

// FetchFolderByNameOrId loads a non-shared folder by UID or vault path via `get`.
// Returns ErrNonSharedFolderNotFound when the backend response is empty or contains
// no folder_uid so callers can drop the resource from state.
func FetchFolderByNameOrId(ctx context.Context, apiManager *api.ApiManager, nameOrUID string) (*NonSharedFolderResponse, error) {
	if nameOrUID == "" {
		return nil, fmt.Errorf("non-shared folder name or id is empty")
	}
	apiResp, err := apiManager.ExecuteCommand(ctx, BuildGetFolderCommand(nameOrUID), errOpGet)
	if err != nil {
		return nil, err
	}
	if apiResp == nil || apiResp.Data == nil {
		return nil, ErrNonSharedFolderNotFound
	}
	out, err := ParseFolderResponse(apiResp.Data)
	if err != nil {
		return nil, err
	}
	if out == nil || out.FolderUID == "" {
		return nil, ErrNonSharedFolderNotFound
	}
	return out, nil
}

// ParseFolderResponse unmarshals API result data into NonSharedFolderResponse.
// If data is nil, returns (nil, nil) so callers can treat as missing resource.
func ParseFolderResponse(data any) (*NonSharedFolderResponse, error) {
	if data == nil {
		return nil, nil
	}
	var out NonSharedFolderResponse
	if err := utils.UnmarshalApiResponse(data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MapResponseToModel maps a get FOLDER_UID --format json response into m.
// Sets identity (Id, Name), folder_location (parent of path; null when at root),
// and records (set of UIDs, or SetNull(StringType) when empty). Data sources that
// prefer empty values over nulls should coerce them after this call.
func MapResponseToModel(api *NonSharedFolderResponse, m *Model) error {
	if api == nil {
		return fmt.Errorf("non-shared folder API response is nil")
	}
	m.Id = types.StringValue(api.FolderUID)
	m.Name = types.StringValue(api.Name)

	if api.Path != "" {
		parent, _ := SplitFolderPath(api.Path)
		if parent != "" {
			m.FolderLocation = types.StringValue(parent)
		} else {
			m.FolderLocation = types.StringNull()
		}
	}

	if api.Records != nil {
		uids := make([]string, 0, len(api.Records))
		for _, r := range api.Records {
			if r.RecordUID != "" {
				uids = append(uids, r.RecordUID)
			}
		}
		if len(uids) > 0 {
			recordSet, diags := types.SetValueFrom(context.Background(), types.StringType, uids)
			if diags.HasError() {
				return fmt.Errorf("failed to build records set: %v", diags)
			}
			m.Records = recordSet
		} else {
			m.Records = types.SetNull(types.StringType)
		}
	}

	return nil
}
