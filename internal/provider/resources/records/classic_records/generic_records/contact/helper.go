// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	"strings"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func buildRecordAddCommand(data ContactResourceModel) string {
	var extra []string

	// name.
	if data.Name != nil && !data.Name.IsNull() {
		if j, err := data.Name.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
			commonrecordsutils.AppendOptionalJSONAdd(&extra, FlagName, j)
		}
	}
	// company.
	commonrecordsutils.AppendOptionalTextField(&extra, FlagTextCompany, data.Company)

	// email.
	commonrecordsutils.AppendOptionalTextField(&extra, FlagEmail, data.Email)

	// phone.
	for i := range data.Phone {
		p := data.Phone[i]
		phoneType := strings.TrimSpace(p.Type.ValueString())
		if phoneType == "" {
			continue
		}
		if j, err := p.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
			// Commander accepts per-type keys like: phone.Mobile='$JSON:{...}'
			commonrecordsutils.AppendOptionalJSONAdd(&extra, FlagPhonePrefix+phoneType, j)
		}
	}

	// address.
	commonrecordsutils.AppendOptionalTextField(&extra, FlagAddressRef, data.AddressRef)

	// custom.
	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)

	return commonrecordsutils.BuildRecordAdd(data.FolderLocation, data.Title.ValueString(), commonrecordsutils.RecordTypeContact, extra, custom, data.Notes)
}

func updateHasMutations(plan, state ContactResourceModel) bool {
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
	if !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom) {
		return true
	}
	return false
}

func buildRecordUpdateCommand(recordUID string, plan, state ContactResourceModel) string {
	var extra []string

	// name.
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

		// Emit changed/new phone.<Type> values.
		for t, pj := range planByType {
			sj := stateByType[t]
			commonrecordsutils.AppendChangedJSONField(&extra, FlagPhonePrefix+t, pj, sj, pj != sj)
		}
		// Clear phone.<Type> that existed but is removed in plan.
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
	return commonrecordsutils.BuildRecordUpdate(recordUID, plan.Title, state.Title, extra, customPlan, customState, plan.Notes, state.Notes)
}

func nameToJSON(n *commonrecordsutils.NameValue) (string, error) {
	if n == nil || n.IsNull() {
		return "", nil
	}
	return n.ToJSON()
}

func mapVaultRecordToModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *ContactResourceModel) {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)

	m.Name = commonrecordsutils.NameFromFields(rec.Fields, "")
	m.Company = commonrecordsutils.FirstStringField(rec.Fields, commonrecordsutils.FieldTypeText, "company")
	m.Email = commonrecordsutils.FirstStringField(rec.Fields, commonrecordsutils.FieldTypeEmail, "")
	m.Phone = commonrecordsutils.PhonesFromField(rec.Fields, "")
	m.AddressRef = commonrecordsutils.FirstRefUID(rec.Fields, commonrecordsutils.FieldTypeAddressRef, "")

	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
}
