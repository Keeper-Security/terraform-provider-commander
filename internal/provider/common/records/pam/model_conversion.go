// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamrecords

import (
	"context"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// The shared helper functions in this package operate on the union model
// CommonPamSettingsFieldResourceModel. PAM Database resources expose only the
// database protocol sub-blocks in their schema, while PAM Machine / PAM
// Directory resources expose only the machine/directory protocol sub-blocks.
// The Terraform plugin framework requires the Go struct to match the schema
// object type exactly, so each record type uses its own typed model and the
// converters below bridge to the shared internal model.

func databaseConnectionToCommon(c *DatabaseConnectionResourceModel) *CommonPamSettingsConnectionResourceModel {
	if c == nil {
		return nil
	}
	return &CommonPamSettingsConnectionResourceModel{
		connectionScalarsModel: c.connectionScalarsModel,
		Mysql:                  c.Mysql,
		PostgreSql:             c.PostgreSql,
		SqlServer:              c.SqlServer,
		MariaDb:                c.MariaDb,
		Oracle:                 c.Oracle,
	}
}

func commonConnectionToDatabase(c *CommonPamSettingsConnectionResourceModel) *DatabaseConnectionResourceModel {
	if c == nil {
		return nil
	}
	return &DatabaseConnectionResourceModel{
		connectionScalarsModel: c.connectionScalarsModel,
		Mysql:                  c.Mysql,
		PostgreSql:             c.PostgreSql,
		SqlServer:              c.SqlServer,
		MariaDb:                c.MariaDb,
		Oracle:                 c.Oracle,
	}
}

func machineDirectoryConnectionToCommon(c *MachineDirectoryConnectionResourceModel) *CommonPamSettingsConnectionResourceModel {
	if c == nil {
		return nil
	}
	return &CommonPamSettingsConnectionResourceModel{
		connectionScalarsModel: c.connectionScalarsModel,
		Kubernetes:             c.Kubernetes,
		Rdp:                    c.Rdp,
		Ssh:                    c.Ssh,
		Telnet:                 c.Telnet,
		Vnc:                    c.Vnc,
	}
}

func commonConnectionToMachineDirectory(c *CommonPamSettingsConnectionResourceModel) *MachineDirectoryConnectionResourceModel {
	if c == nil {
		return nil
	}
	return &MachineDirectoryConnectionResourceModel{
		connectionScalarsModel: c.connectionScalarsModel,
		Kubernetes:             c.Kubernetes,
		Rdp:                    c.Rdp,
		Ssh:                    c.Ssh,
		Telnet:                 c.Telnet,
		Vnc:                    c.Vnc,
	}
}

func databasePamSettingsToCommon(m *DatabasePamSettingsFieldResourceModel) *CommonPamSettingsFieldResourceModel {
	if m == nil {
		return nil
	}
	return &CommonPamSettingsFieldResourceModel{
		pamSettingsCommonModel: m.pamSettingsCommonModel,
		Connection:             databaseConnectionToCommon(m.Connection),
	}
}

func commonPamSettingsToDatabase(m *CommonPamSettingsFieldResourceModel) *DatabasePamSettingsFieldResourceModel {
	if m == nil {
		return nil
	}
	return &DatabasePamSettingsFieldResourceModel{
		pamSettingsCommonModel: m.pamSettingsCommonModel,
		Connection:             commonConnectionToDatabase(m.Connection),
	}
}

func machineDirectoryPamSettingsToCommon(m *MachineDirectoryPamSettingsFieldResourceModel) *CommonPamSettingsFieldResourceModel {
	if m == nil {
		return nil
	}
	return &CommonPamSettingsFieldResourceModel{
		pamSettingsCommonModel: m.pamSettingsCommonModel,
		Connection:             machineDirectoryConnectionToCommon(m.Connection),
	}
}

func commonPamSettingsToMachineDirectory(m *CommonPamSettingsFieldResourceModel) *MachineDirectoryPamSettingsFieldResourceModel {
	if m == nil {
		return nil
	}
	return &MachineDirectoryPamSettingsFieldResourceModel{
		pamSettingsCommonModel: m.pamSettingsCommonModel,
		Connection:             commonConnectionToMachineDirectory(m.Connection),
	}
}

// ExtractDatabasePamSettingsFromResponse is the database-typed wrapper around
// ExtractPamSettingsFromResponse.
func ExtractDatabasePamSettingsFromResponse(rec *utils.VaultRecordGetResponse, existingState *DatabasePamSettingsFieldResourceModel) *DatabasePamSettingsFieldResourceModel {
	return commonPamSettingsToDatabase(ExtractPamSettingsFromResponse(rec, databasePamSettingsToCommon(existingState)))
}

// ExtractMachineDirectoryPamSettingsFromResponse is the machine/directory-typed
// wrapper around ExtractPamSettingsFromResponse.
func ExtractMachineDirectoryPamSettingsFromResponse(rec *utils.VaultRecordGetResponse, existingState *MachineDirectoryPamSettingsFieldResourceModel) *MachineDirectoryPamSettingsFieldResourceModel {
	return commonPamSettingsToMachineDirectory(ExtractPamSettingsFromResponse(rec, machineDirectoryPamSettingsToCommon(existingState)))
}

// ValidateDatabasePamSettingsFieldsNotRemoved is the database-typed wrapper
// around ValidatePamSettingsFieldsNotRemoved.
func ValidateDatabasePamSettingsFieldsNotRemoved(plan, state *DatabasePamSettingsFieldResourceModel) diag.Diagnostics {
	return ValidatePamSettingsFieldsNotRemoved(databasePamSettingsToCommon(plan), databasePamSettingsToCommon(state))
}

// ValidateMachineDirectoryPamSettingsFieldsNotRemoved is the
// machine/directory-typed wrapper around ValidatePamSettingsFieldsNotRemoved.
func ValidateMachineDirectoryPamSettingsFieldsNotRemoved(plan, state *MachineDirectoryPamSettingsFieldResourceModel) diag.Diagnostics {
	return ValidatePamSettingsFieldsNotRemoved(machineDirectoryPamSettingsToCommon(plan), machineDirectoryPamSettingsToCommon(state))
}

// ApplyDatabasePamSettings is the database-typed wrapper around ApplyPamSettings.
// recordUpdateCmd must be utils.CmdRecordUpdate (classic) or
// utils.CmdNsfRecordUpdate (nested-shared).
func ApplyDatabasePamSettings(ctx context.Context, apiManager *api.ApiManager, recordUpdateCmd, recordUID string, plan, state *DatabasePamSettingsFieldResourceModel) error {
	return ApplyPamSettings(ctx, apiManager, recordUpdateCmd, recordUID, databasePamSettingsToCommon(plan), databasePamSettingsToCommon(state))
}

// ApplyMachineDirectoryPamSettings is the machine/directory-typed wrapper
// around ApplyPamSettings. recordUpdateCmd must be utils.CmdRecordUpdate
// (classic) or utils.CmdNsfRecordUpdate (nested-shared).
func ApplyMachineDirectoryPamSettings(ctx context.Context, apiManager *api.ApiManager, recordUpdateCmd, recordUID string, plan, state *MachineDirectoryPamSettingsFieldResourceModel) error {
	return ApplyPamSettings(ctx, apiManager, recordUpdateCmd, recordUID, machineDirectoryPamSettingsToCommon(plan), machineDirectoryPamSettingsToCommon(state))
}
