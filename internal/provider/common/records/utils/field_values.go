// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ----- name -----------------------------------------------------------------

type NameValue struct {
	First  types.String `tfsdk:"first"`
	Middle types.String `tfsdk:"middle"`
	Last   types.String `tfsdk:"last"`
}

type nameJSON struct {
	First  string `json:"first,omitempty"`
	Middle string `json:"middle,omitempty"`
	Last   string `json:"last,omitempty"`
}

func (n *NameValue) IsNull() bool {
	if n == nil {
		return true
	}
	return (n.First.IsNull() || n.First.IsUnknown() || strings.TrimSpace(n.First.ValueString()) == "") &&
		(n.Middle.IsNull() || n.Middle.IsUnknown() || strings.TrimSpace(n.Middle.ValueString()) == "") &&
		(n.Last.IsNull() || n.Last.IsUnknown() || strings.TrimSpace(n.Last.ValueString()) == "")
}

func (n *NameValue) ToJSON() (string, error) {
	if n == nil || n.IsNull() {
		return "", nil
	}
	j := nameJSON{
		First:  stringOrEmpty(n.First),
		Middle: stringOrEmpty(n.Middle),
		Last:   stringOrEmpty(n.Last),
	}
	b, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func NameFromFields(fields []utils.VaultRecordFieldResponse, label string) *NameValue {
	for i := range fields {
		f := &fields[i]
		if f.Type != FieldTypeName || (label != "" && f.Label != label) {
			continue
		}
		var arr []nameJSON
		if err := json.Unmarshal(f.Value, &arr); err != nil || len(arr) == 0 {
			return nil
		}
		v := arr[0]
		return &NameValue{
			First:  utils.StringOrNull(v.First),
			Middle: utils.StringOrNull(v.Middle),
			Last:   utils.StringOrNull(v.Last),
		}
	}
	return nil
}

func NameEqual(a, b *NameValue) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.First.Equal(b.First) && a.Middle.Equal(b.Middle) && a.Last.Equal(b.Last)
}

// ----- phone (one entry) ----------------------------------------------------

type PhoneValue struct {
	Region types.String `tfsdk:"region"`
	Number types.String `tfsdk:"number"`
	Ext    types.String `tfsdk:"ext"`
	Type   types.String `tfsdk:"type"`
}

type phoneJSON struct {
	Region string `json:"region,omitempty"`
	Number string `json:"number,omitempty"`
	Ext    string `json:"ext,omitempty"`
	Type   string `json:"type,omitempty"`
}

func (p *PhoneValue) ToJSON() (string, error) {
	if p == nil {
		return "", nil
	}
	j := phoneJSON{
		Region: stringOrEmpty(p.Region),
		Number: stringOrEmpty(p.Number),
		Ext:    stringOrEmpty(p.Ext),
		Type:   stringOrEmpty(p.Type),
	}
	b, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func PhonesFromField(fields []utils.VaultRecordFieldResponse, label string) []PhoneValue {
	for i := range fields {
		f := &fields[i]
		if f.Type != FieldTypePhone || (label != "" && f.Label != label) {
			continue
		}
		var arr []phoneJSON
		if err := json.Unmarshal(f.Value, &arr); err != nil {
			return nil
		}
		out := make([]PhoneValue, 0, len(arr))
		for _, v := range arr {
			out = append(out, PhoneValue{
				Region: utils.StringOrNull(v.Region),
				Number: utils.StringOrNull(v.Number),
				Ext:    utils.StringOrNull(v.Ext),
				Type:   utils.StringOrNull(v.Type),
			})
		}
		return out
	}
	return nil
}

func PhonesToJSONArray(phones []PhoneValue) (string, error) {
	if len(phones) == 0 {
		return "", nil
	}
	arr := make([]phoneJSON, 0, len(phones))
	for i := range phones {
		p := &phones[i]
		arr = append(arr, phoneJSON{
			Region: stringOrEmpty(p.Region),
			Number: stringOrEmpty(p.Number),
			Ext:    stringOrEmpty(p.Ext),
			Type:   stringOrEmpty(p.Type),
		})
	}
	b, err := json.Marshal(arr)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func PhoneSliceEqual(a, b []PhoneValue) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Region.Equal(b[i].Region) || !a[i].Number.Equal(b[i].Number) ||
			!a[i].Ext.Equal(b[i].Ext) || !a[i].Type.Equal(b[i].Type) {
			return false
		}
	}
	return true
}

// ----- address --------------------------------------------------------------

type AddressValue struct {
	Street1 types.String `tfsdk:"street1"`
	Street2 types.String `tfsdk:"street2"`
	City    types.String `tfsdk:"city"`
	State   types.String `tfsdk:"state"`
	Zip     types.String `tfsdk:"zip"`
	Country types.String `tfsdk:"country"`
}

type addressJSON struct {
	Street1 string `json:"street1,omitempty"`
	Street2 string `json:"street2,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Zip     string `json:"zip,omitempty"`
	Country string `json:"country,omitempty"`
}

func (a *AddressValue) ToJSON() (string, error) {
	if a == nil {
		return "", nil
	}
	j := addressJSON{
		Street1: stringOrEmpty(a.Street1),
		Street2: stringOrEmpty(a.Street2),
		City:    stringOrEmpty(a.City),
		State:   stringOrEmpty(a.State),
		Zip:     stringOrEmpty(a.Zip),
		Country: stringOrEmpty(a.Country),
	}
	b, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func AddressFromFields(fields []utils.VaultRecordFieldResponse, label string) *AddressValue {
	for i := range fields {
		f := &fields[i]
		if f.Type != FieldTypeAddress || (label != "" && f.Label != label) {
			continue
		}
		var arr []addressJSON
		if err := json.Unmarshal(f.Value, &arr); err != nil || len(arr) == 0 {
			return nil
		}
		v := arr[0]
		return &AddressValue{
			Street1: utils.StringOrNull(v.Street1),
			Street2: utils.StringOrNull(v.Street2),
			City:    utils.StringOrNull(v.City),
			State:   utils.StringOrNull(v.State),
			Zip:     utils.StringOrNull(v.Zip),
			Country: utils.StringOrNull(v.Country),
		}
	}
	return nil
}

func AddressEqual(a, b *AddressValue) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Street1.Equal(b.Street1) && a.Street2.Equal(b.Street2) && a.City.Equal(b.City) &&
		a.State.Equal(b.State) && a.Zip.Equal(b.Zip) && a.Country.Equal(b.Country)
}

// ----- host -----------------------------------------------------------------

type HostValue struct {
	HostName types.String `tfsdk:"host_name"`
	Port     types.String `tfsdk:"port"`
}

type hostJSON struct {
	HostName string `json:"hostName,omitempty"`
	Port     string `json:"port,omitempty"`
}

func (h *HostValue) ToJSON() (string, error) {
	if h == nil {
		return "", nil
	}
	j := hostJSON{
		HostName: stringOrEmpty(h.HostName),
		Port:     stringOrEmpty(h.Port),
	}
	b, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func HostFromFields(fields []utils.VaultRecordFieldResponse, label string) *HostValue {
	for i := range fields {
		f := &fields[i]
		if f.Type != FieldTypeHost || (label != "" && f.Label != label) {
			continue
		}
		var arr []hostJSON
		if err := json.Unmarshal(f.Value, &arr); err != nil || len(arr) == 0 {
			return nil
		}
		v := arr[0]
		return &HostValue{
			HostName: utils.StringOrNull(v.HostName),
			Port:     utils.StringOrNull(v.Port),
		}
	}
	return nil
}

func HostEqual(a, b *HostValue) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.HostName.Equal(b.HostName) && a.Port.Equal(b.Port)
}

// IsEmpty reports whether both hostname and port are unset.
func (h *HostValue) IsEmpty() bool {
	if h == nil {
		return true
	}
	return StringUnset(h.HostName) && StringUnset(h.Port)
}

// HostFlatToJSON serializes flat hostname/port Terraform fields to Keeper host JSON.
func HostFlatToJSON(hostname, port types.String) (string, bool) {
	h := &HostValue{HostName: hostname, Port: port}
	if h.IsEmpty() {
		return "", false
	}
	j, err := h.ToJSON()
	if err != nil || strings.TrimSpace(j) == "" {
		return "", false
	}
	return j, true
}

// FlatHostFromFields reads a host field into flat hostname and port strings.
func FlatHostFromFields(fields []utils.VaultRecordFieldResponse) (hostname, port types.String) {
	if host := HostFromFields(fields, ""); host != nil {
		return host.HostName, host.Port
	}
	return types.StringNull(), types.StringNull()
}

// ----- paymentCard ----------------------------------------------------------

type PaymentCardValue struct {
	CardNumber         types.String `tfsdk:"card_number"`
	CardExpirationDate types.String `tfsdk:"card_expiration_date"`
	CardSecurityCode   types.String `tfsdk:"card_security_code"`
}

type paymentCardJSON struct {
	CardNumber         string `json:"cardNumber,omitempty"`
	CardExpirationDate string `json:"cardExpirationDate,omitempty"`
	CardSecurityCode   string `json:"cardSecurityCode,omitempty"`
}

func (p *PaymentCardValue) ToJSON() (string, error) {
	if p == nil {
		return "", nil
	}
	j := paymentCardJSON{
		CardNumber:         stringOrEmpty(p.CardNumber),
		CardExpirationDate: stringOrEmpty(p.CardExpirationDate),
		CardSecurityCode:   stringOrEmpty(p.CardSecurityCode),
	}
	b, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func PaymentCardFromFields(fields []utils.VaultRecordFieldResponse, label string) *PaymentCardValue {
	for i := range fields {
		f := &fields[i]
		if f.Type != FieldTypePaymentCard || (label != "" && f.Label != label) {
			continue
		}
		var arr []paymentCardJSON
		if err := json.Unmarshal(f.Value, &arr); err != nil || len(arr) == 0 {
			return nil
		}
		v := arr[0]
		return &PaymentCardValue{
			CardNumber:         utils.StringOrNull(v.CardNumber),
			CardExpirationDate: utils.StringOrNull(v.CardExpirationDate),
			CardSecurityCode:   utils.StringOrNull(v.CardSecurityCode),
		}
	}
	return nil
}

func PaymentCardEqual(a, b *PaymentCardValue) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.CardNumber.Equal(b.CardNumber) && a.CardExpirationDate.Equal(b.CardExpirationDate) &&
		a.CardSecurityCode.Equal(b.CardSecurityCode)
}

// ----- bankAccount ----------------------------------------------------------

type BankAccountValue struct {
	AccountType   types.String `tfsdk:"account_type"`
	OtherType     types.String `tfsdk:"other_type"`
	RoutingNumber types.String `tfsdk:"routing_number"`
	AccountNumber types.String `tfsdk:"account_number"`
}

type bankAccountJSON struct {
	AccountType   string `json:"accountType,omitempty"`
	OtherType     string `json:"otherType,omitempty"`
	RoutingNumber string `json:"routingNumber,omitempty"`
	AccountNumber string `json:"accountNumber,omitempty"`
}

func (b *BankAccountValue) ToJSON() (string, error) {
	if b == nil {
		return "", nil
	}
	j := bankAccountJSON{
		AccountType:   stringOrEmpty(b.AccountType),
		OtherType:     stringOrEmpty(b.OtherType),
		RoutingNumber: stringOrEmpty(b.RoutingNumber),
		AccountNumber: stringOrEmpty(b.AccountNumber),
	}
	out, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func BankAccountFromFields(fields []utils.VaultRecordFieldResponse, label string) *BankAccountValue {
	for i := range fields {
		f := &fields[i]
		if f.Type != FieldTypeBankAccount || (label != "" && f.Label != label) {
			continue
		}
		var arr []bankAccountJSON
		if err := json.Unmarshal(f.Value, &arr); err != nil || len(arr) == 0 {
			return nil
		}
		v := arr[0]
		return &BankAccountValue{
			AccountType:   utils.StringOrNull(v.AccountType),
			OtherType:     utils.StringOrNull(v.OtherType),
			RoutingNumber: utils.StringOrNull(v.RoutingNumber),
			AccountNumber: utils.StringOrNull(v.AccountNumber),
		}
	}
	return nil
}

func BankAccountEqual(a, b *BankAccountValue) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.AccountType.Equal(b.AccountType) && a.OtherType.Equal(b.OtherType) &&
		a.RoutingNumber.Equal(b.RoutingNumber) && a.AccountNumber.Equal(b.AccountNumber)
}

// ----- securityQuestion -----------------------------------------------------

type SecurityQuestionValue struct {
	Question types.String `tfsdk:"question"`
	Answer   types.String `tfsdk:"answer"`
}

type securityQuestionJSON struct {
	Question string `json:"question,omitempty"`
	Answer   string `json:"answer,omitempty"`
}

func (s *SecurityQuestionValue) ToJSON() (string, error) {
	if s == nil {
		return "", nil
	}
	arr := []securityQuestionJSON{{
		Question: stringOrEmpty(s.Question),
		Answer:   stringOrEmpty(s.Answer),
	}}
	b, err := json.Marshal(arr)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func SecurityQuestionFromFields(fields []utils.VaultRecordFieldResponse, label string) *SecurityQuestionValue {
	for i := range fields {
		f := &fields[i]
		if f.Type != FieldTypeSecurityQuestion || (label != "" && f.Label != label) {
			continue
		}
		var arr []securityQuestionJSON
		if err := json.Unmarshal(f.Value, &arr); err != nil || len(arr) == 0 {
			return nil
		}
		v := arr[0]
		return &SecurityQuestionValue{
			Question: utils.StringOrNull(v.Question),
			Answer:   utils.StringOrNull(v.Answer),
		}
	}
	return nil
}

func SecurityQuestionEqual(a, b *SecurityQuestionValue) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Question.Equal(b.Question) && a.Answer.Equal(b.Answer)
}

// ----- keyPair --------------------------------------------------------------

type KeyPairValue struct {
	PublicKey  types.String `tfsdk:"public_key"`
	PrivateKey types.String `tfsdk:"private_key"`
}

type keyPairJSON struct {
	PublicKey  string `json:"publicKey,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
}

func (k *KeyPairValue) ToJSON() (string, error) {
	if k == nil {
		return "", nil
	}
	j := keyPairJSON{
		PublicKey:  stringOrEmpty(k.PublicKey),
		PrivateKey: stringOrEmpty(k.PrivateKey),
	}
	b, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func KeyPairFromFields(fields []utils.VaultRecordFieldResponse, label string) *KeyPairValue {
	for i := range fields {
		f := &fields[i]
		if f.Type != FieldTypeKeyPair || (label != "" && f.Label != label) {
			continue
		}
		var arr []keyPairJSON
		if err := json.Unmarshal(f.Value, &arr); err != nil || len(arr) == 0 {
			return nil
		}
		v := arr[0]
		return &KeyPairValue{
			PublicKey:  utils.StringOrNull(v.PublicKey),
			PrivateKey: utils.StringOrNull(v.PrivateKey),
		}
	}
	return nil
}

func KeyPairEqual(a, b *KeyPairValue) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.PublicKey.Equal(b.PublicKey) && a.PrivateKey.Equal(b.PrivateKey)
}

// IsEmpty reports whether both public and private keys are unset.
func (k *KeyPairValue) IsEmpty() bool {
	if k == nil {
		return true
	}
	return StringUnset(k.PublicKey) && StringUnset(k.PrivateKey)
}

// KeyPairFlatToJSON serializes flat public/private key Terraform fields to Keeper keyPair JSON.
func KeyPairFlatToJSON(publicKey, privateKey types.String) (string, bool) {
	kp := &KeyPairValue{PublicKey: publicKey, PrivateKey: privateKey}
	if kp.IsEmpty() {
		return "", false
	}
	j, err := kp.ToJSON()
	if err != nil || strings.TrimSpace(j) == "" {
		return "", false
	}
	return j, true
}

// FlatKeyPairFromFields reads a keyPair field into flat public and private key strings.
func FlatKeyPairFromFields(fields []utils.VaultRecordFieldResponse) (publicKey, privateKey types.String) {
	if kp := KeyPairFromFields(fields, ""); kp != nil {
		return kp.PublicKey, kp.PrivateKey
	}
	return types.StringNull(), types.StringNull()
}

// ----- scalar / ref helpers -------------------------------------------------

// StringUnset reports whether a Terraform string is null, unknown, or blank.
func StringUnset(s types.String) bool {
	return s.IsNull() || s.IsUnknown() || strings.TrimSpace(s.ValueString()) == ""
}

func stringOrEmpty(s types.String) string {
	if s.IsNull() || s.IsUnknown() {
		return ""
	}
	return s.ValueString()
}

// FirstStringField returns the first string from a typed field's value array (standard fields).
func FirstStringField(fields []utils.VaultRecordFieldResponse, fieldType, label string) types.String {
	for i := range fields {
		f := &fields[i]
		if f.Type != fieldType {
			continue
		}
		if label != "" && f.Label != label {
			continue
		}
		var vals []string
		if err := json.Unmarshal(f.Value, &vals); err != nil {
			return types.StringNull()
		}
		if len(vals) > 0 && strings.TrimSpace(vals[0]) != "" {
			return types.StringValue(vals[0])
		}
		return types.StringNull()
	}
	return types.StringNull()
}

// FirstStringFieldCustom searches both fields and custom arrays.
func FirstStringFieldCustom(rec *utils.VaultRecordGetResponse, fieldType, label string) types.String {
	v := FirstStringField(rec.Fields, fieldType, label)
	if !v.IsNull() {
		return v
	}
	return FirstStringField(rec.Custom, fieldType, label)
}

// FirstRefUID returns the first UID from addressRef / cardRef / fileRef style arrays.
func FirstRefUID(fields []utils.VaultRecordFieldResponse, fieldType, label string) types.String {
	for i := range fields {
		f := &fields[i]
		if f.Type != fieldType {
			continue
		}
		if label != "" && f.Label != label {
			continue
		}
		var vals []string
		if err := json.Unmarshal(f.Value, &vals); err != nil || len(vals) == 0 {
			return types.StringNull()
		}
		if strings.TrimSpace(vals[0]) != "" {
			return types.StringValue(strings.TrimSpace(vals[0]))
		}
	}
	return types.StringNull()
}

// FirstStringFieldAnyLabel returns the first string from a typed field's value
// array, ignoring the label. Use for record-type-specific fields like
// `password`, `wifiEncryption` whose label is sometimes empty and sometimes
// echoed by the server as the type name.
func FirstStringFieldAnyLabel(fields []utils.VaultRecordFieldResponse, fieldType string) types.String {
	for i := range fields {
		f := &fields[i]
		if f.Type != fieldType {
			continue
		}
		var vals []string
		if err := json.Unmarshal(f.Value, &vals); err != nil {
			continue
		}
		if len(vals) > 0 && strings.TrimSpace(vals[0]) != "" {
			return types.StringValue(vals[0])
		}
	}
	return types.StringNull()
}

// FirstBoolFieldAnyLabel returns the first boolean from a typed field's value
// array, ignoring the label. See FirstStringFieldAnyLabel for rationale.
func FirstBoolFieldAnyLabel(fields []utils.VaultRecordFieldResponse, fieldType string) types.Bool {
	for i := range fields {
		f := &fields[i]
		if f.Type != fieldType {
			continue
		}
		var vals []bool
		if err := json.Unmarshal(f.Value, &vals); err == nil && len(vals) > 0 {
			return types.BoolValue(vals[0])
		}
	}
	return types.BoolNull()
}

// FirstBoolField returns the first boolean from a typed field's value array.
// Returns a null types.Bool when the field is missing or empty.
func FirstBoolField(fields []utils.VaultRecordFieldResponse, fieldType, label string) types.Bool {
	for i := range fields {
		f := &fields[i]
		if f.Type != fieldType {
			continue
		}
		if label != "" && f.Label != label {
			continue
		}
		if label == "" && f.Label != "" {
			continue
		}
		var vals []bool
		if err := json.Unmarshal(f.Value, &vals); err == nil && len(vals) > 0 {
			return types.BoolValue(vals[0])
		}
	}
	return types.BoolNull()
}

// EpochMillisField reads first numeric element as int64 epoch ms → ISO date string for Terraform.
func EpochMillisField(fields []utils.VaultRecordFieldResponse, fieldType, label string) types.String {
	for i := range fields {
		f := &fields[i]
		if f.Type != fieldType {
			continue
		}
		if label != "" && f.Label != label {
			continue
		}
		var vals []json.Number
		if err := json.Unmarshal(f.Value, &vals); err == nil && len(vals) > 0 {
			if ms, err := vals[0].Int64(); err == nil && ms != 0 {
				return types.StringValue(epochMillisToDateString(ms))
			}
		}
		// try float array
		var fvals []float64
		if err := json.Unmarshal(f.Value, &fvals); err == nil && len(fvals) > 0 {
			ms := int64(fvals[0])
			if ms != 0 {
				return types.StringValue(epochMillisToDateString(ms))
			}
		}
	}
	return types.StringNull()
}

// epochMillisToDateString converts Keeper epoch-ms to RFC3339 UTC for Terraform.
func epochMillisToDateString(ms int64) string {
	if ms == 0 {
		return ""
	}
	return time.UnixMilli(ms).UTC().Format(time.RFC3339)
}

// EpochMillisFieldUnlabeled returns the first epoch-ms field of the given type (ignores label).
func EpochMillisFieldUnlabeled(fields []utils.VaultRecordFieldResponse, fieldType string) types.String {
	for i := range fields {
		f := &fields[i]
		if f.Type != fieldType {
			continue
		}
		var vals []json.Number
		if err := json.Unmarshal(f.Value, &vals); err == nil && len(vals) > 0 {
			if ms, err := vals[0].Int64(); err == nil && ms != 0 {
				return types.StringValue(epochMillisToDateString(ms))
			}
		}
		var fvals []float64
		if err := json.Unmarshal(f.Value, &fvals); err == nil && len(fvals) > 0 {
			ms := int64(fvals[0])
			if ms != 0 {
				return types.StringValue(epochMillisToDateString(ms))
			}
		}
	}
	return types.StringNull()
}

// DateStringToEpochMillisOrZero parses RFC3339 or YYYY-MM-DD into epoch milliseconds; empty → 0, nil error.
func DateStringToEpochMillisOrZero(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UnixMilli(), nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return 0, err
	}
	return t.UnixMilli(), nil
}
