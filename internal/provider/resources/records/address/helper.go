// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package address

import (
	"context"
	"strings"

	records "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func buildRecordAddCommand(data AddressResourceModel) string {
	var extra []string

	if j, err := buildAddressJSON(data); err == nil && strings.TrimSpace(j) != "" {
		records.AppendOptionalJSONAdd(&extra, FlagAddress, j)
	}

	custom := records.NormalizeCustomFromPlan(data.Custom)
	return records.BuildRecordAdd(data.Folder, data.Title.ValueString(), records.RecordTypeAddress, extra, custom, data.Notes)
}

func updateHasMutations(plan, state AddressResourceModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) {
		return true
	}
	if !addressFieldsEqual(plan, state) {
		return true
	}
	if !records.CustomFieldsEqual(plan.Custom, state.Custom) {
		return true
	}
	return false
}

func buildRecordUpdateCommand(recordUID string, plan, state AddressResourceModel) string {
	var extra []string

	planJSON, planErr := buildAddressJSON(plan)
	stateJSON, stateErr := buildAddressJSON(state)
	changed := planJSON != stateJSON || planErr != stateErr
	records.AppendChangedJSONField(&extra, FlagAddress, planJSON, stateJSON, changed)

	customPlan := records.NormalizeCustomFromPlan(plan.Custom)
	customState := records.NormalizeCustomFromPlan(state.Custom)
	return records.BuildRecordUpdate(recordUID, plan.Title, state.Title, extra, customPlan, customState, plan.Notes, state.Notes)
}

// buildAddressJSON serializes the flat address fields as a single Keeper
// `address` value object. Returns an empty string when no field is set.
func buildAddressJSON(m AddressResourceModel) (string, error) {
	if !addressHasValue(m) {
		return "", nil
	}
	av := records.AddressValue{
		Street1: m.Street1,
		Street2: m.Street2,
		City:    m.City,
		State:   m.State,
		Zip:     m.Zip,
		Country: m.Country,
	}
	return av.ToJSON()
}

func addressHasValue(m AddressResourceModel) bool {
	for _, f := range []types.String{m.Street1, m.Street2, m.City, m.State, m.Zip, m.Country} {
		if !f.IsNull() && !f.IsUnknown() && strings.TrimSpace(f.ValueString()) != "" {
			return true
		}
	}
	return false
}

func addressFieldsEqual(a, b AddressResourceModel) bool {
	return a.Street1.Equal(b.Street1) &&
		a.Street2.Equal(b.Street2) &&
		a.City.Equal(b.City) &&
		a.State.Equal(b.State) &&
		a.Zip.Equal(b.Zip) &&
		a.Country.Equal(b.Country)
}

func mapVaultRecordToModel(ctx context.Context, rec *utils.VaultRecordGetResponse, stateFolder types.String, m *AddressResourceModel) diag.Diagnostics {
	records.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)

	if av := records.AddressFromFields(rec.Fields, ""); av != nil {
		m.Street1 = av.Street1
		m.Street2 = av.Street2
		m.City = av.City
		m.State = av.State
		m.Zip = av.Zip
		m.Country = av.Country
	} else {
		m.Street1 = types.StringNull()
		m.Street2 = types.StringNull()
		m.City = types.StringNull()
		m.State = types.StringNull()
		m.Zip = types.StringNull()
		m.Country = types.StringNull()
	}

	// Parse share record permissions from the API response.
	shareMap, diags := records.ParseSharePermissionsFromResponse(ctx, rec.UserPermissions)
	if diags.HasError() {
		return diags
	}
	m.Share = shareMap
	return nil
}
