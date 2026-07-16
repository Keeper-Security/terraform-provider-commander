// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package non_shared_folder

import (
	"context"
	"errors"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ErrNonSharedFolderNotFound is returned when the get command yields no usable folder (deleted, wrong id/path, or empty response).
var ErrNonSharedFolderNotFound = errors.New("non-shared folder not found")

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
// and records (set of UIDs, or SetNull(StringType) when empty). When the API
// response omits the records key entirely (api.Records == nil) m.Records is
// left untouched so data sources that prefer empty values over nulls can
// coerce them after this call.
func MapResponseToModel(ctx context.Context, api *NonSharedFolderResponse, m *Model) error {
	if api == nil {
		return fmt.Errorf("non-shared folder API response is nil")
	}

	m.Id = types.StringValue(api.FolderUID)
	m.Name = types.StringValue(api.Name)
	m.FolderLocation = utils.ExtractFolderValue(&api.FolderLocation, m.FolderLocation)

	// Link the records to the folder.
	if api.Records == nil {
		return nil
	}
	uids := make([]string, 0, len(api.Records))
	for _, r := range api.Records {
		if r.RecordUID != "" {
			uids = append(uids, r.RecordUID)
		}
	}
	set, err := folderutils.FolderRecordsToSet(ctx, uids)
	if err != nil {
		return err
	}
	m.Records = set
	return nil
}
