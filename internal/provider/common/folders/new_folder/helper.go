// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package new_folder

import (
	"context"
	"errors"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/new_share"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ErrKeeperDriveFolderNotFound is returned when nsf-get cannot find the folder
// (either the API returned a not-found error or the response is empty).
// Callers use errors.Is(err, ErrKeeperDriveFolderNotFound) to detect deletion
// out-of-band and drop the resource from Terraform state.
var ErrNestedSharedFolderNotFound = errors.New("nested shared folder not found")

// BuildNewFolderGetCommand builds: nsf-get "ID_OR_NAME" --format json.
func BuildNewFolderGetCommand(idOrName string) string {
	return fmt.Sprintf(`%s "%s" %s`, utils.CmdNsfGet, idOrName, utils.FlagFormatJSON)
}

// FetchNsfFolderByNameOrId loads a Keeper Drive folder by UID or name via
// nsf-get. Returns ErrKeeperDriveFolderNotFound when the backend reports the
// folder does not exist (so callers can RemoveResource) or any other error
// when the call itself fails.
func FetchNsfFolderByNameOrId(ctx context.Context, apiManager *api.ApiManager, nameOrUID string) (*NewFolderGetResponse, error) {
	if nameOrUID == "" {
		return nil, fmt.Errorf("new folder name or id is empty")
	}
	apiResp, err := apiManager.ExecuteCommand(ctx, BuildNewFolderGetCommand(nameOrUID), folderutils.ErrOpGet)
	if err != nil {
		if errors.Is(err, api.ErrResourceNotFound) {
			return nil, ErrNestedSharedFolderNotFound
		}
		return nil, err
	}

	if apiResp == nil || apiResp.Data == nil {
		return nil, ErrNestedSharedFolderNotFound
	}

	out, err := ParseNewFolderGetResponse(apiResp.Data)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, ErrNestedSharedFolderNotFound
	}
	if out.FolderUID == "" {
		return nil, fmt.Errorf("get response missing folder_uid")
	}
	return out, nil
}

// ParseNewFolderGetResponse unmarshals nsf-get result data into NewFolderGetResponse.
func ParseNewFolderGetResponse(data any) (*NewFolderGetResponse, error) {
	if data == nil {
		return nil, nil
	}
	var out NewFolderGetResponse
	if err := utils.UnmarshalApiResponse(data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MapResponseToModel maps nsf-get API data into the Terraform model's identity
// fields (Id, Name, FolderLocation) and Records. The share block is populated
// separately by callers via new_share.MapResponseToModel so the Optional-only
// semantics of `share` can be honored (skip the update when state.Share is null).
// When the response omits the records key entirely (apiData.Records == nil)
// m.Records is left untouched.
func MapResponseToModel(ctx context.Context, apiData *NewFolderGetResponse, m *Model) error {
	if apiData == nil {
		return fmt.Errorf("folder API response is nil")
	}

	m.Id = types.StringValue(apiData.FolderUID)
	m.Name = types.StringValue(apiData.Name)
	m.FolderLocation = utils.ExtractFolderValue(&apiData.FolderLocation, m.FolderLocation)

	// Link the records to the folder.
	if apiData.Records == nil {
		return nil
	}
	uids := make([]string, 0, len(apiData.Records))
	for _, r := range apiData.Records {
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

// CollectFolderSharePermissions merges user_permissions and team_permissions
// from an nsf-get response into a single slice for new_share.MapResponseToModel,
// dropping any entry whose access_type is AccessTypeApplication. Returns
// nil when apiData is nil.
func CollectFolderSharePermissions(apiData *NewFolderGetResponse) []new_share.UserPermissionEntry {
	if apiData == nil {
		return nil
	}
	out := make([]new_share.UserPermissionEntry, 0, len(apiData.UserPermissions)+len(apiData.TeamPermissions))
	out = appendNonApplicationEntries(out, apiData.UserPermissions)
	out = appendNonApplicationEntries(out, apiData.TeamPermissions)
	return out
}

// appendNonApplicationEntries copies entries from src into dst, skipping any
// whose access_type equals AccessTypeApplication.
func appendNonApplicationEntries(dst, src []new_share.UserPermissionEntry) []new_share.UserPermissionEntry {
	for _, e := range src {
		if e.AccessType == AccessTypeApplication {
			continue
		}
		dst = append(dst, e)
	}
	return dst
}
