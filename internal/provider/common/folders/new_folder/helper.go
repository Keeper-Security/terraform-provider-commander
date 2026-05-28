// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package new_folder

import (
	"context"
	"errors"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	folderutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/folders/utils"
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

// FetchNewFolderByNameOrId loads a Keeper Drive folder by UID or name via
// nsf-get. Returns ErrKeeperDriveFolderNotFound when the backend reports the
// folder does not exist (so callers can RemoveResource) or any other error
// when the call itself fails.
func FetchNewFolderByNameOrId(ctx context.Context, apiManager *api.ApiManager, nameOrUID string) (*NewFolderGetResponse, error) {
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
		return nil, fmt.Errorf("get response missing nested_share_folder_uid")
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
// fields (Id, Name). folder_location is preserved from existing state since the
// nsf-get response does not include a parent path. The share block is populated
// separately by callers via new_share.MapResponseToModel so the Optional-only
// semantics of `share` can be honored (skip the update when state.Share is null).
func MapResponseToModel(apiData *NewFolderGetResponse, m *Model) error {
	if apiData == nil {
		return fmt.Errorf("folder API response is nil")
	}
	m.Id = types.StringValue(apiData.FolderUID)
	m.Name = types.StringValue(apiData.Name)
	if m.FolderLocation.IsUnknown() {
		m.FolderLocation = types.StringNull()
	}
	return nil
}
