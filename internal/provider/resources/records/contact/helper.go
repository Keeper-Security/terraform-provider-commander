// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	"context"
	"strings"

	records "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func buildRecordAddCommand(data ContactResourceModel) string {
	var extra []string

	// name.
	if data.Name != nil && !data.Name.IsNull() {
		if j, err := data.Name.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
			records.AppendOptionalJSONAdd(&extra, FlagName, j)
		}
	}
	// company.
	records.AppendOptionalTextField(&extra, FlagTextCompany, data.Company)

	// email.
	records.AppendOptionalTextField(&extra, FlagEmail, data.Email)

	// phone.
	for i := range data.Phone {
		p := data.Phone[i]
		phoneType := strings.TrimSpace(p.Type.ValueString())
		if phoneType == "" {
			continue
		}
		if j, err := p.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
			// Commander accepts per-type keys like: phone.Mobile='$JSON:{...}'
			records.AppendOptionalJSONAdd(&extra, FlagPhonePrefix+phoneType, j)
		}
	}

	// address.
	records.AppendOptionalTextField(&extra, FlagAddressRef, data.AddressRef)

	// custom.
	custom := records.NormalizeCustomFromPlan(data.Custom)

	return records.BuildRecordAdd(data.Folder, data.Title.ValueString(), records.RecordTypeContact, extra, custom, data.Notes)
}

func updateHasMutations(plan, state ContactResourceModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.Company.Equal(state.Company) ||
		!plan.Email.Equal(state.Email) ||
		!plan.AddressRef.Equal(state.AddressRef) {
		return true
	}
	if !records.NameEqual(plan.Name, state.Name) {
		return true
	}
	if !records.PhoneSliceEqual(plan.Phone, state.Phone) {
		return true
	}
	if !records.CustomFieldsEqual(plan.Custom, state.Custom) {
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
	records.AppendChangedJSONField(&extra, FlagName, planJSON, stateJSON, changed)

	records.AppendChangedStringField(&extra, FlagTextCompany, plan.Company, state.Company)
	records.AppendChangedStringField(&extra, FlagEmail, plan.Email, state.Email)

	if !records.PhoneSliceEqual(plan.Phone, state.Phone) {
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
			records.AppendChangedJSONField(&extra, FlagPhonePrefix+t, pj, sj, pj != sj)
		}
		// Clear phone.<Type> that existed but is removed in plan.
		for t, sj := range stateByType {
			if _, stillPresent := planByType[t]; stillPresent {
				continue
			}
			records.AppendChangedJSONField(&extra, FlagPhonePrefix+t, "", sj, true)
		}
	}

	records.AppendChangedStringField(&extra, FlagAddressRef, plan.AddressRef, state.AddressRef)

	customPlan := records.NormalizeCustomFromPlan(plan.Custom)
	customState := records.NormalizeCustomFromPlan(state.Custom)
	return records.BuildRecordUpdate(recordUID, plan.Title, state.Title, extra, customPlan, customState, plan.Notes, state.Notes)
}

func nameToJSON(n *records.NameValue) (string, error) {
	if n == nil || n.IsNull() {
		return "", nil
	}
	return n.ToJSON()
}

func mapVaultRecordToModel(ctx context.Context, rec *utils.VaultRecordGetResponse, stateFolder types.String, m *ContactResourceModel) diag.Diagnostics {
	records.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.Name = records.NameFromFields(rec.Fields, "")
	m.Company = records.FirstStringField(rec.Fields, records.FieldTypeText, "company")
	m.Email = records.FirstStringField(rec.Fields, records.FieldTypeEmail, "")
	m.Phone = records.PhonesFromField(rec.Fields, "")
	m.AddressRef = records.FirstRefUID(rec.Fields, records.FieldTypeAddressRef, "")

	// Parse share record permissions from the API response.
	shareMap, diags := records.ParseSharePermissionsFromResponse(ctx, rec.UserPermissions)
	if diags.HasError() {
		return diags
	}
	m.Share = shareMap
	return nil
}
