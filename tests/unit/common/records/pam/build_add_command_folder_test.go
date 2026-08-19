// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamrecords_test

import (
	"strings"
	"testing"

	commonpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_database"
	commonpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_directory"
	commonpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_machine"
	commonpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_remote_browser"
	commonpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/pam/pam_user"
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestPamDatabaseBuildAddCommand_OmitsBlankFolderLocation(t *testing.T) {
	t.Parallel()

	base := commonpamdatabase.PamDatabaseResourceModel{
		BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
			Title: types.StringValue("Test DB"),
		},
	}

	tests := []struct {
		name   string
		folder types.String
	}{
		{name: "null", folder: types.StringNull()},
		{name: "unknown", folder: types.StringUnknown()},
		{name: "empty", folder: types.StringValue("")},
		{name: "whitespace", folder: types.StringValue("   ")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := base
			data.FolderLocation = tt.folder
			cmd := commonpamdatabase.BuildAddCommand(utils.CmdNsfRecordAdd, data)
			if strings.Contains(cmd, utils.FlagFolder) {
				t.Fatalf("command should omit %s for blank folder_location, got: %s", utils.FlagFolder, cmd)
			}
		})
	}

	withFolder := base
	withFolder.FolderLocation = types.StringValue("NSF PAM/Resources")
	cmd := commonpamdatabase.BuildAddCommand(utils.CmdNsfRecordAdd, withFolder)
	if !strings.Contains(cmd, utils.FlagFolder) || !strings.Contains(cmd, "NSF PAM/Resources") {
		t.Fatalf("command should include folder flag and path, got: %s", cmd)
	}
}

func TestPamMachineBuildAddCommand_OmitsBlankFolderLocation(t *testing.T) {
	t.Parallel()

	base := commonpammachine.PamMachineResourceModel{
		BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
			Title: types.StringValue("Test Machine"),
		},
	}

	data := base
	data.FolderLocation = types.StringValue("")
	cmd := commonpammachine.BuildAddCommand(utils.CmdNsfRecordAdd, data)
	if strings.Contains(cmd, utils.FlagFolder) {
		t.Fatalf("command should omit %s for blank folder_location, got: %s", utils.FlagFolder, cmd)
	}

	data.FolderLocation = types.StringValue("NSF PAM/Machines")
	cmd = commonpammachine.BuildAddCommand(utils.CmdNsfRecordAdd, data)
	if !strings.Contains(cmd, utils.FlagFolder) || !strings.Contains(cmd, "NSF PAM/Machines") {
		t.Fatalf("command should include folder flag and path, got: %s", cmd)
	}
}

func TestPamDirectoryBuildAddCommand_OmitsBlankFolderLocation(t *testing.T) {
	t.Parallel()

	base := commonpamdirectory.PamDirectoryResourceModel{
		BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
			Title: types.StringValue("Test Directory"),
		},
	}

	data := base
	data.FolderLocation = types.StringNull()
	cmd := commonpamdirectory.BuildAddCommand(utils.CmdNsfRecordAdd, data)
	if strings.Contains(cmd, utils.FlagFolder) {
		t.Fatalf("command should omit %s for null folder_location, got: %s", utils.FlagFolder, cmd)
	}

	data.FolderLocation = types.StringValue("NSF PAM/Directories")
	cmd = commonpamdirectory.BuildAddCommand(utils.CmdNsfRecordAdd, data)
	if !strings.Contains(cmd, utils.FlagFolder) || !strings.Contains(cmd, "NSF PAM/Directories") {
		t.Fatalf("command should include folder flag and path, got: %s", cmd)
	}
}

func TestPamRemoteBrowserBuildAddCommand_OmitsBlankFolderLocation(t *testing.T) {
	t.Parallel()

	base := commonpamremotebrowser.PamRemoteBrowserResourceModel{
		BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
			Title: types.StringValue("Test RBI"),
		},
		Url: types.StringValue("https://example.com"),
	}

	data := base
	data.FolderLocation = types.StringValue("  ")
	cmd := commonpamremotebrowser.BuildAddCommand(utils.CmdNsfRecordAdd, data)
	if strings.Contains(cmd, utils.FlagFolder) {
		t.Fatalf("command should omit %s for blank folder_location, got: %s", utils.FlagFolder, cmd)
	}

	data.FolderLocation = types.StringValue("NSF PAM/Browsers")
	cmd = commonpamremotebrowser.BuildAddCommand(utils.CmdNsfRecordAdd, data)
	if !strings.Contains(cmd, utils.FlagFolder) || !strings.Contains(cmd, "NSF PAM/Browsers") {
		t.Fatalf("command should include folder flag and path, got: %s", cmd)
	}
}

func TestPamUserBuildAddCommand_OmitsBlankFolderLocation(t *testing.T) {
	t.Parallel()

	base := commonpamuser.PamUserSharedModel{
		BaseVaultRecordModel: commonrecordsutils.BaseVaultRecordModel{
			Title: types.StringValue("Test User"),
		},
		Login: types.StringValue("svc_user"),
	}

	data := base
	data.FolderLocation = types.StringValue("")
	cmd := commonpamuser.BuildAddCommand(utils.CmdNsfRecordAdd, data)
	if strings.Contains(cmd, utils.FlagFolder) {
		t.Fatalf("command should omit %s for blank folder_location, got: %s", utils.FlagFolder, cmd)
	}

	data.FolderLocation = types.StringValue("NSF PAM/Users")
	cmd = commonpamuser.BuildAddCommand(utils.CmdNsfRecordAdd, data)
	if !strings.Contains(cmd, utils.FlagFolder) || !strings.Contains(cmd, "NSF PAM/Users") {
		t.Fatalf("command should include folder flag and path, got: %s", cmd)
	}
}
