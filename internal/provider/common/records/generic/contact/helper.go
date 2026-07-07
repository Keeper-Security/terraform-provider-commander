// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	"strings"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for a contact record.
func BuildAddCommand(data ContactModel) string {
	var extra []string

	if data.Name != nil && !data.Name.IsNull() {
		if j, err := data.Name.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
			commonrecordsutils.AppendOptionalJSONAdd(&extra, FlagName, j)
		}
	}
	commonrecordsutils.AppendOptionalTextField(&extra, FlagTextCompany, data.Company)
	commonrecordsutils.AppendOptionalTextField(&extra, FlagEmail, data.Email)

	for i := range data.Phone {
		p := data.Phone[i]
		phoneType := strings.TrimSpace(p.Type.ValueString())
		if phoneType == "" {
			continue
		}
		if j, err := p.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
			commonrecordsutils.AppendOptionalJSONAdd(&extra, FlagPhonePrefix+phoneType, j)
		}
	}

	commonrecordsutils.AppendOptionalTextField(&extra, FlagAddressRef, data.AddressRef)

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(utils.CmdRecordAdd, commonrecordsutils.RecordTypeContact, data.Title.ValueString(), data.FolderLocation, extra, custom, data.Notes)
}

// UpdateHasMutations reports whether plan differs from state on updatable contact fields.
func UpdateHasMutations(plan, state ContactModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.Company.Equal(state.Company) ||
		!plan.Email.Equal(state.Email) ||
		!plan.AddressRef.Equal(state.AddressRef) {
		return true
	}
	if !commonrecordsutils.NameEqual(plan.Name, state.Name) {
		return true
	}
	if !commonrecordsutils.PhoneSliceEqual(plan.Phone, state.Phone) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed contact fields.
func BuildUpdateCommand(recordUID string, plan, state ContactModel) string {
	var extra []string

	planJSON, planErr := nameToJSON(plan.Name)
	stateJSON, stateErr := nameToJSON(state.Name)
	changed := planJSON != stateJSON || planErr != stateErr
	commonrecordsutils.AppendChangedJSONField(&extra, FlagName, planJSON, stateJSON, changed)

	commonrecordsutils.AppendChangedStringField(&extra, FlagTextCompany, plan.Company, state.Company)
	commonrecordsutils.AppendChangedStringField(&extra, FlagEmail, plan.Email, state.Email)

	if !commonrecordsutils.PhoneSliceEqual(plan.Phone, state.Phone) {
		planByType := map[string]string{}
		stateByType := map[string]string{}

		for i := range plan.Phone {
			p := plan.Phone[i]
			t := strings.TrimSpace(p.Type.ValueString())
			if t == "" {
				continue
			}
			if j, err := p.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
				planByType[t] = j
			}
		}
		for i := range state.Phone {
			p := state.Phone[i]
			t := strings.TrimSpace(p.Type.ValueString())
			if t == "" {
				continue
			}
			if j, err := p.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
				stateByType[t] = j
			}
		}

		for t, pj := range planByType {
			sj := stateByType[t]
			commonrecordsutils.AppendChangedJSONField(&extra, FlagPhonePrefix+t, pj, sj, pj != sj)
		}
		for t, sj := range stateByType {
			if _, stillPresent := planByType[t]; stillPresent {
				continue
			}
			commonrecordsutils.AppendChangedJSONField(&extra, FlagPhonePrefix+t, "", sj, true)
		}
	}

	commonrecordsutils.AppendChangedStringField(&extra, FlagAddressRef, plan.AddressRef, state.AddressRef)

	customPlan := commonrecordsutils.NormalizeCustomFromPlan(plan.Custom)
	customState := commonrecordsutils.NormalizeCustomFromPlan(state.Custom)
	return commonrecordsutils.BuildRecordUpdate(utils.CmdRecordUpdate, recordUID, plan.Title, state.Title, extra, customPlan, customState, plan.Notes, state.Notes)
}

// MapVaultRecordGetResponseToContactModel fills state from a `get <uid> --format json` payload.
func MapVaultRecordGetResponseToContactModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *ContactModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.Name = commonrecordsutils.NameFromFields(rec.Fields, "")
	m.Company = commonrecordsutils.FirstStringField(rec.Fields, commonrecordsutils.FieldTypeText, "company")
	m.Email = commonrecordsutils.FirstStringField(rec.Fields, commonrecordsutils.FieldTypeEmail, "")
	m.Phone = commonrecordsutils.PhonesFromField(rec.Fields, "")
	m.AddressRef = commonrecordsutils.FirstRefUID(rec.Fields, commonrecordsutils.FieldTypeAddressRef, "")
	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}

func nameToJSON(n *commonrecordsutils.NameValue) (string, error) {
	if n == nil || n.IsNull() {
		return "", nil
	}
	return n.ToJSON()
}
