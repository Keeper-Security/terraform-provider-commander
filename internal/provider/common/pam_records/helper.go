package pamrecords

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
)

func FetchVaultRecord(ctx context.Context, apiManager *api.ApiManager, recordUID string) (*api.RequestResultResponse, error) {
	command := fmt.Sprintf("%s '%s' %s", utils.CmdGetRecord, recordUID, utils.FlagFormatJSON)
	apiResp, err := apiManager.ExecuteCommand(ctx, command, utils.ErrSummaryFetchVaultRecordFailed)

	return apiResp, err
}

// ExtractFirstTextFieldValue finds a field with type "text" and the given label,
// then returns the first string from its value array. Returns "" if not found.
// Reusable across pamMachine, pamDatabase, pamDirectory, etc.
func ExtractFirstTextFieldValue(fields []utils.VaultRecordFieldResponse, label string) string {
	for i := range fields {
		f := &fields[i]
		if f.Type != "text" || f.Label != label {
			continue
		}
		var vals []string
		if err := json.Unmarshal(f.Value, &vals); err != nil {
			continue
		}
		if len(vals) > 0 {
			return strings.TrimSpace(vals[0])
		}
	}
	return ""
}

func MoveRecordFromSourceToDestination(ctx context.Context, apiManager *api.ApiManager, recordUID string, planFolderData string) error {
	src := recordUID

	dest := planFolderData
	if dest == "" {
		dest = "/"
	}

	if src == dest {
		return nil
	}

	command := fmt.Sprintf("%s '%s' '%s' %s", utils.CmdMv, src, dest, utils.FlagForce)
	_, err := apiManager.ExecuteCommand(ctx, command, utils.ErrSummaryMoveRecordFailed)
	if err != nil {
		return err
	}
	return nil
}
