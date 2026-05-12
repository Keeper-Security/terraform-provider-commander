// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package records

// Vault field type strings from `get <uid> --format json` (field.type).
const (
	FieldTypeLogin            = "login"
	FieldTypePassword         = "password"
	FieldTypeURL              = "url"
	FieldTypeEmail            = "email"
	FieldTypeText             = "text"
	FieldTypeMultiline        = "multiline"
	FieldTypeSecret           = "secret"
	FieldTypeNote             = "note"
	FieldTypeName             = "name"
	FieldTypePhone            = "phone"
	FieldTypeAddress          = "address"
	FieldTypeHost             = "host"
	FieldTypePaymentCard      = "paymentCard"
	FieldTypeBankAccount      = "bankAccount"
	FieldTypeSecurityQuestion = "securityQuestion"
	FieldTypeKeyPair          = "keyPair"
	FieldTypeAccountNumber    = "accountNumber"
	FieldTypeLicenseNumber    = "licenseNumber"
	FieldTypePinCode          = "pinCode"
	FieldTypeBirthDate        = "birthDate"
	FieldTypeDate             = "date"
	FieldTypeExpirationDate   = "expirationDate"
	FieldTypeAddressRef       = "addressRef"
	FieldTypeCardRef          = "cardRef"
	FieldTypeFileRef          = "fileRef"
	FieldTypeCheckbox         = "checkbox"
	FieldTypeOneTimeCode      = "oneTimeCode"
	FieldTypeOTP              = "otp"
	FieldTypeWifiEncryption   = "wifiEncryption"
	FieldTypeIsSSIDHidden     = "isSSIDHidden"
)

// Standard Keeper record types (Commander --record-type values).
const (
	RecordTypeContact             = "contact"
	RecordTypeLogin               = "login"
	RecordTypeAddress             = "address"
	RecordTypeBankAccount         = "bankAccount"
	RecordTypeBankCard            = "bankCard"
	RecordTypeBirthCertificate    = "birthCertificate"
	RecordTypeDatabaseCredentials = "databaseCredentials"
	RecordTypeDriverLicense       = "driverLicense"
	RecordTypeEncryptedNotes      = "encryptedNotes"
	RecordTypeHealthInsurance     = "healthInsurance"
	RecordTypeMembership          = "membership"
	RecordTypePassport            = "passport"
	RecordTypeServerCredentials   = "serverCredentials"
	RecordTypeSoftwareLicense     = "softwareLicense"
	RecordTypeSsnCard             = "ssnCard"
	RecordTypeSshKeys             = "sshKeys"
	RecordTypeWifiCredentials     = "wifiCredentials"
)
