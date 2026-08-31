// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package bankaccount

import (
	"strings"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for a bank account record.
func BuildAddCommand(cmd string, data BankAccountModel) string {
	var extra []string

	if data.BankAccount != nil && !data.BankAccount.IsNull() {
		if j, err := data.BankAccount.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
			commonrecordsutils.AppendOptionalJSONAdd(&extra, FlagBankAccount, j)
		}
	}
	if data.Name != nil && !data.Name.IsNull() {
		if j, err := data.Name.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
			commonrecordsutils.AppendOptionalJSONAdd(&extra, FlagName, j)
		}
	}
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagLogin, data.Login)
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagPassword, data.Password)
	commonrecordsutils.AppendOptionalScalarAdd(&extra, FlagURL, data.WebsiteAddress)
	commonrecordsutils.AppendOptionalTextField(&extra, FlagCardRef, data.CardRef)

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(cmd, commonrecordsutils.RecordTypeBankAccount, data.Title.ValueString(), data.FolderLocation, extra, custom, data.Notes)
}

// UpdateHasMutations reports whether plan differs from state on updatable bank account fields.
func UpdateHasMutations(plan, state BankAccountModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.Login.Equal(state.Login) ||
		!plan.Password.Equal(state.Password) ||
		!plan.WebsiteAddress.Equal(state.WebsiteAddress) ||
		!plan.CardRef.Equal(state.CardRef) {
		return true
	}
	if !commonrecordsutils.BankAccountEqual(plan.BankAccount, state.BankAccount) {
		return true
	}
	if !commonrecordsutils.NameEqual(plan.Name, state.Name) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed bank account fields.
func BuildUpdateCommand(cmd string, recordUID string, plan, state BankAccountModel) string {
	var extra []string

	bankAccountPlanJSON, bankAccountPlanErr := plan.BankAccount.ToJSON()
	bankAccountStateJSON, bankAccountStateErr := state.BankAccount.ToJSON()
	bankAccountChanged := bankAccountPlanJSON != bankAccountStateJSON || bankAccountPlanErr != bankAccountStateErr
	commonrecordsutils.AppendChangedJSONField(&extra, FlagBankAccount, bankAccountPlanJSON, bankAccountStateJSON, bankAccountChanged)

	namePlanJSON, namePlanErr := plan.Name.ToJSON()
	nameStateJSON, nameStateErr := state.Name.ToJSON()
	nameChanged := namePlanJSON != nameStateJSON || namePlanErr != nameStateErr
	commonrecordsutils.AppendChangedJSONField(&extra, FlagName, namePlanJSON, nameStateJSON, nameChanged)

	commonrecordsutils.AppendChangedStringField(&extra, FlagLogin, plan.Login, state.Login)
	commonrecordsutils.AppendChangedStringField(&extra, FlagPassword, plan.Password, state.Password)
	commonrecordsutils.AppendChangedStringField(&extra, FlagURL, plan.WebsiteAddress, state.WebsiteAddress)
	commonrecordsutils.AppendChangedStringField(&extra, FlagCardRef, plan.CardRef, state.CardRef)

	customPlan := commonrecordsutils.NormalizeCustomFromPlan(plan.Custom)
	customState := commonrecordsutils.NormalizeCustomFromPlan(state.Custom)
	return commonrecordsutils.BuildRecordUpdate(cmd, recordUID, plan.Title, state.Title, extra, customPlan, customState, plan.Notes, state.Notes)
}

// MapVaultRecordGetResponseToBankAccountModel fills state from a `get <uid> --format json` payload.
func MapVaultRecordGetResponseToBankAccountModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *BankAccountModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.BankAccount = commonrecordsutils.BankAccountFromFields(rec.Fields, "")
	m.Name = commonrecordsutils.NameFromFields(rec.Fields, "")
	m.Login = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeLogin)
	m.Password = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypePassword)
	m.WebsiteAddress = commonrecordsutils.FirstStringFieldAnyLabel(rec.Fields, commonrecordsutils.FieldTypeURL)
	m.CardRef = commonrecordsutils.FirstRefUID(rec.Fields, commonrecordsutils.FieldTypeCardRef, "")
	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}
