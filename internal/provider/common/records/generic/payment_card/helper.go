// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package paymentcard

import (
	"strings"

	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BuildAddCommand builds a record-add command for a payment card record.
func BuildAddCommand(cmd string, data PaymentCardModel) string {
	var extra []string

	if data.PaymentCard != nil && !data.PaymentCard.IsNull() {
		if j, err := data.PaymentCard.ToJSON(); err == nil && strings.TrimSpace(j) != "" {
			commonrecordsutils.AppendOptionalJSONAdd(&extra, FlagPaymentCard, j)
		}
	}
	commonrecordsutils.AppendOptionalTextField(&extra, FlagCardholderName, data.CardholderName)
	commonrecordsutils.AppendOptionalTextField(&extra, FlagPinCode, data.PinCode)
	commonrecordsutils.AppendOptionalTextField(&extra, FlagAddressRef, data.AddressRef)

	custom := commonrecordsutils.NormalizeCustomFromPlan(data.Custom)
	return commonrecordsutils.BuildRecordAdd(cmd, commonrecordsutils.RecordTypeBankCard, data.Title.ValueString(), data.FolderLocation, extra, custom, data.Notes)
}

// UpdateHasMutations reports whether plan differs from state on updatable payment card fields.
func UpdateHasMutations(plan, state PaymentCardModel) bool {
	if !plan.Title.Equal(state.Title) ||
		!plan.Notes.Equal(state.Notes) ||
		!plan.CardholderName.Equal(state.CardholderName) ||
		!plan.PinCode.Equal(state.PinCode) ||
		!plan.AddressRef.Equal(state.AddressRef) {
		return true
	}
	if !commonrecordsutils.PaymentCardEqual(plan.PaymentCard, state.PaymentCard) {
		return true
	}
	return !commonrecordsutils.CustomFieldsEqual(plan.Custom, state.Custom)
}

// BuildUpdateCommand builds a record-update command for changed payment card fields.
func BuildUpdateCommand(cmd string, recordUID string, plan, state PaymentCardModel) string {
	var extra []string

	planJSON, planErr := plan.PaymentCard.ToJSON()
	stateJSON, stateErr := state.PaymentCard.ToJSON()
	changed := planJSON != stateJSON || planErr != stateErr
	commonrecordsutils.AppendChangedJSONField(&extra, FlagPaymentCard, planJSON, stateJSON, changed)

	commonrecordsutils.AppendChangedStringField(&extra, FlagCardholderName, plan.CardholderName, state.CardholderName)
	commonrecordsutils.AppendChangedStringField(&extra, FlagPinCode, plan.PinCode, state.PinCode)
	commonrecordsutils.AppendChangedStringField(&extra, FlagAddressRef, plan.AddressRef, state.AddressRef)

	customPlan := commonrecordsutils.NormalizeCustomFromPlan(plan.Custom)
	customState := commonrecordsutils.NormalizeCustomFromPlan(state.Custom)
	return commonrecordsutils.BuildRecordUpdate(cmd, recordUID, plan.Title, state.Title, extra, customPlan, customState, plan.Notes, state.Notes)
}

// MapVaultRecordGetResponseToPaymentCardModel fills state from a `get <uid> --format json` payload.
func MapVaultRecordGetResponseToPaymentCardModel(rec *utils.VaultRecordGetResponse, stateFolder types.String, m *PaymentCardModel) diag.Diagnostics {
	commonrecordsutils.MapBaseVaultRecord(rec, stateFolder, &m.BaseVaultRecordModel)
	m.PaymentCard = commonrecordsutils.PaymentCardFromFields(rec.Fields, "")
	m.CardholderName = commonrecordsutils.FirstStringField(rec.Fields, commonrecordsutils.FieldTypeText, CardholderNameFieldLabel)
	m.PinCode = commonrecordsutils.FirstStringField(rec.Fields, commonrecordsutils.FieldTypePinCode, "")
	m.AddressRef = commonrecordsutils.FirstRefUID(rec.Fields, commonrecordsutils.FieldTypeAddressRef, "")
	m.Custom = commonrecordsutils.ParseCustomFields(rec.Custom)
	return nil
}
