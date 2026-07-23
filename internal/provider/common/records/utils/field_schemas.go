// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// BaseRecordAttributes returns id, title, notes, folder shared by all standard vault record resources.
func BaseRecordAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			Description:         "The unique identifier (UID) of the vault record.",
			MarkdownDescription: "The unique identifier (UID) of the vault record.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"title": schema.StringAttribute{
			Required:            true,
			Description:         "Record title.",
			MarkdownDescription: "Record title.",
			Validators: []validator.String{
				utils.StringMinLengthValidator("Title", 1, false),
			},
		},
		"notes": schema.StringAttribute{
			Optional:            true,
			Description:         "Manage note for the record.",
			MarkdownDescription: "Manage note for the record.",
			Validators: []validator.String{
				utils.StringMinLengthValidator("Notes", 1, true),
			},
		},
		"folder_location": schema.StringAttribute{
			Optional:            true,
			Description:         "Folder path or UID where the record is to be stored.",
			MarkdownDescription: "Folder `path` or `UID` where the record is to be stored.",
			Validators: []validator.String{
				utils.StringMinLengthValidator("Folder", 1, true),
			},
		},
	}
}

func DataSourceBaseRecordAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"id": dschema.StringAttribute{
			Computed:            true,
			Description:         "Unique identifier (UID) of the vault record.",
			MarkdownDescription: "Unique identifier (UID) of the vault record.",
		},
		"title": dschema.StringAttribute{
			Computed:            true,
			Description:         "Record title.",
			MarkdownDescription: "Record title.",
		},
		"notes": dschema.StringAttribute{
			Computed:            true,
			Description:         "Note of the record.",
			MarkdownDescription: "Note of the record.",
		},
		"folder_location": dschema.StringAttribute{
			Computed:            true,
			Description:         "Folder path or UID where the record is to be stored.",
			MarkdownDescription: "Folder `path` or `UID` where the record is to be stored.",
		},
	}
}

// NameNestedSchema is a single nested object for a Keeper `name` field.
func NameNestedSchema(optional bool) schema.SingleNestedAttribute {
	req := !optional
	return schema.SingleNestedAttribute{
		Optional:            optional,
		Required:            req,
		Description:         "Person name (first, middle, last).",
		MarkdownDescription: "Person name (`name` field): first, middle, last.",
		Attributes: map[string]schema.Attribute{
			"first": schema.StringAttribute{
				Required:            true,
				Description:         "First name.",
				MarkdownDescription: "First name.",
			},
			"middle": schema.StringAttribute{
				Optional:            true,
				Description:         "Middle name.",
				MarkdownDescription: "Middle name.",
			},
			"last": schema.StringAttribute{
				Required:            true,
				Description:         "Last name.",
				MarkdownDescription: "Last name.",
			},
		},
	}
}

// PhoneListSchema is a list of phone entries (Keeper `phone` field value array).
func PhoneListSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:            true,
		Description:         "Manage phone numbers for the record.",
		MarkdownDescription: "Manage phone numbers for the record.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"region": schema.StringAttribute{
					Optional:            true,
					Description:         "Region or country code (e.g. US, +1).",
					MarkdownDescription: "Region or country code.",
				},
				"number": schema.StringAttribute{
					Optional:            true,
					Description:         "Phone number.",
					MarkdownDescription: "Phone number.",
				},
				"ext": schema.StringAttribute{
					Optional:            true,
					Description:         "Extension.",
					MarkdownDescription: "Extension.",
				},
				"type": schema.StringAttribute{
					Optional:            true,
					Description:         "Phone type: Mobile, Home, or Work.",
					MarkdownDescription: "Phone type: `Mobile`, `Home`, or `Work`.",
					Validators: []validator.String{
						utils.StringOneOfValidator("Phone type", []string{"Mobile", "Home", "Work"}, true),
					},
				},
			},
		},
	}
}

// AddressNestedSchema maps Keeper `address` object.
func AddressNestedSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Description:         "Physical address.",
		MarkdownDescription: "Physical address (`address` field).",
		Attributes: map[string]schema.Attribute{
			"street1": schema.StringAttribute{Optional: true, Description: "Street line 1.", MarkdownDescription: "Street line 1."},
			"street2": schema.StringAttribute{Optional: true, Description: "Street line 2.", MarkdownDescription: "Street line 2."},
			"city":    schema.StringAttribute{Optional: true, Description: "City.", MarkdownDescription: "City."},
			"state":   schema.StringAttribute{Optional: true, Description: "State or province.", MarkdownDescription: "State or province."},
			"zip":     schema.StringAttribute{Optional: true, Description: "Postal code.", MarkdownDescription: "Postal code."},
			"country": schema.StringAttribute{Optional: true, Description: "ISO 3166-1 alpha-2 country code.", MarkdownDescription: "ISO 3166-1 alpha-2 country code."},
		},
	}
}

// HostNestedSchema maps Keeper `host` object.
func HostNestedSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Description:         "Hostname and port.",
		MarkdownDescription: "Hostname and port (`host` field).",
		Attributes: map[string]schema.Attribute{
			"host_name": schema.StringAttribute{
				Optional:            true,
				Description:         "Hostname or IP address.",
				MarkdownDescription: "Hostname or IP address.",
				Validators: []validator.String{
					utils.StringMinLengthValidator("Hostname", 1, true),
				},
			},
			"port": schema.StringAttribute{
				Optional:            true,
				Description:         "Port (string per Keeper API).",
				MarkdownDescription: "Port (string per Keeper API).",
				Validators: []validator.String{
					utils.NumericStringValidator("Port", true),
				},
			},
		},
	}
}

// PaymentCardNestedSchema maps Keeper `paymentCard`.
func PaymentCardNestedSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Description:         "Payment card details.",
		MarkdownDescription: "Payment card details (`paymentCard` field).",
		Attributes: map[string]schema.Attribute{
			"card_number": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         "Card number.",
				MarkdownDescription: "Card number.",
			},
			"card_expiration_date": schema.StringAttribute{
				Optional:            true,
				Description:         "Expiration MM/YYYY.",
				MarkdownDescription: "Expiration `MM/YYYY`.",
			},
			"card_security_code": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         "Security code.",
				MarkdownDescription: "Security code.",
			},
		},
	}
}

// BankAccountNestedSchema maps Keeper `bankAccount`.
func BankAccountNestedSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Description:         "Bank account details.",
		MarkdownDescription: "Bank account details (`bankAccount` field).",
		Attributes: map[string]schema.Attribute{
			"account_type":   stringField(true, false, false, false, "Account type", "Checking, Savings, or Other.", "**Account type**: `Checking`, `Savings`, or `Other`.").withValidators(utils.StringOneOfValidator("Account type", []string{"Checking", "Savings", "Other"}, true)).resource(),
			"other_type":     stringField(true, false, false, false, "Other type", "Description when account_type is Other.", "Description when **account_type** is **Other**.").resource(),
			"routing_number": schema.StringAttribute{Optional: true, Description: "Routing number.", MarkdownDescription: "Routing number."},
			"account_number": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         "Account number.",
				MarkdownDescription: "Account number.",
			},
		},
	}
}

// SecurityQuestionNestedSchema maps Keeper `securityQuestion` (single Q/A pair).
func SecurityQuestionNestedSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Description:         "Security question and answer.",
		MarkdownDescription: "**Security question** and **answer**.",
		Attributes: map[string]schema.Attribute{
			"question": stringField(true, false, false, false, "Question", "Question.", "**Question**.").resource(),
			"answer":   stringField(true, false, true, false, "Answer", "Answer.", "**Answer**.").resource(),
		},
	}
}

// KeyPairNestedSchema maps Keeper `keyPair`.
func KeyPairNestedSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Description:         "SSH public/private key pair.",
		MarkdownDescription: "SSH public/private key pair (`keyPair` field).",
		Attributes: map[string]schema.Attribute{
			"public_key":  stringField(true, false, false, false, "Public key", "Public key.", "**Public key**.").resource(),
			"private_key": stringField(true, false, true, false, "Private key", "Private key.", "**Private key**.").resource(),
		},
	}
}

// ---------------------------------------------------------------------------
// stringField builder
//
// Reusable string attribute factory for generic vault record schemas. Callers
// compose attributes with stringField(...).resource() or .dataSource() instead
// of duplicating Optional/Required/Sensitive/Computed flags and validators.
//
// Typical usage:
//
//	// Optional resource field with default min-length validation:
//	stringField(true, false, false, false, "Login", "Login identifier.", "**Login** identifier.").resource()
//
//	// Required resource field:
//	stringField(false, true, false, false, "SSID", "...", "...").resource()
//
//	// Custom validator (overrides default min-length):
//	stringField(true, false, false, false, "Port", "...", "...").
//	    withValidators(utils.NumericStringValidator("Port", true)).resource()
//
//	// Computed data-source field (name is unused; pass ""):
//	stringField(false, false, false, true, "", "...", "...").dataSource()
// ---------------------------------------------------------------------------

// stringFieldConfig holds shared string attribute settings for resource and data-source builders.
type stringFieldConfig struct {
	// Terraform schema flags (resource fields use optional/required; data-source fields use computed).
	optional  bool
	required  bool
	sensitive bool
	computed  bool

	// validatorDisplayName is shown in validation error messages (e.g. "Login must be at least 1
	// character(s) long."). This is NOT the Terraform attribute key—that is set separately in
	// each record's schema map (e.g. "login": OptionalLoginField()). Ignored by dataSource().
	name string

	// Provider documentation strings surfaced in Terraform registry docs and IDE tooling.
	desc string
	md   string

	// Optional custom validators. When set, these replace the default min-length validator.
	validators []validator.String
}

// stringField builds a stringFieldConfig from schema flags and documentation.
//
// Parameters:
//   - optional:  true for optional resource attributes (null/omitted allowed).
//   - required:  true for required resource attributes (mutually exclusive with optional for resources).
//   - sensitive: true to mask the value in Terraform plan/output (passwords, keys, etc.).
//   - computed:  true for read-only data-source attributes; should be false for resources.
//   - name:      validator display name for default min-length errors on resources; pass "" for
//     data-source fields or when using withValidators.
//   - desc, md:  plain-text and markdown descriptions for provider documentation.
func stringField(optional, required, sensitive, computed bool, name, desc, md string) stringFieldConfig {
	return stringFieldConfig{
		optional:  optional,
		required:  required,
		sensitive: sensitive,
		computed:  computed,
		name:      name,
		desc:      desc,
		md:        md,
	}
}

// withValidators attaches custom validators, overriding the default min-length rule applied in resource().
// Use for fields that need NumericStringValidator, StringOneOfValidator, etc.
func (c stringFieldConfig) withValidators(validators ...validator.String) stringFieldConfig {
	c.validators = validators
	return c
}

// resource returns a Terraform resource schema.StringAttribute.
//
// Validator selection (first match wins):
//  1. If withValidators was called, use those validators as-is.
//  2. Else if name is non-empty, apply StringMinLengthValidator(name, 1, optional):
//     - required fields (optional=false): empty string fails validation.
//     - optional fields (optional=true): null/omitted is allowed, but a set value must be non-empty.
//  3. Else no validators are attached.
func (c stringFieldConfig) resource() schema.StringAttribute {
	attr := schema.StringAttribute{
		Optional:            c.optional,
		Required:            c.required,
		Sensitive:           c.sensitive,
		Description:         c.desc,
		MarkdownDescription: c.md,
	}
	if len(c.validators) > 0 {
		attr.Validators = c.validators
	} else if c.name != "" {
		attr.Validators = []validator.String{
			utils.StringMinLengthValidator(c.name, 1, c.optional),
		}
	}
	return attr
}

// dataSource returns a Terraform data-source schema.StringAttribute.
// Data-source string fields are always computed and read-only; name and validators are not used.
func (c stringFieldConfig) dataSource() dschema.StringAttribute {
	return dschema.StringAttribute{
		Computed:            c.computed,
		Sensitive:           c.sensitive,
		Description:         c.desc,
		MarkdownDescription: c.md,
	}
}

// computedBoolField builds a computed bool attribute for data sources.
func computedBoolField(desc, md string) dschema.BoolAttribute {
	return dschema.BoolAttribute{
		Computed:            true,
		Description:         desc,
		MarkdownDescription: md,
	}
}

// ---------------------------------------------------------------------------
// Named string field helpers for generic vault record schemas.
// Each helper wraps stringField with fixed descriptions and the correct
// optional/required/sensitive/computed combination for that Keeper field.
// ---------------------------------------------------------------------------

// RequiredLoginField is a required login/username string attribute.
func RequiredLoginField() schema.StringAttribute {
	return stringField(false, true, false, false, "Login", "Login identifier.", "**Login** identifier.").resource()
}

// OptionalLoginField is an optional login/username string attribute.
func OptionalLoginField() schema.StringAttribute {
	return stringField(true, false, false, false, "Login", "Login identifier.", "**Login** identifier.").resource()
}

// OptionalPasswordField is an optional sensitive password string attribute.
func OptionalPasswordField() schema.StringAttribute {
	return stringField(true, false, true, false, "Password", "Password.", "**Password**.").resource()
}

// OptionalWebsiteAddressField is an optional website address string attribute.
func OptionalWebsiteAddressField() schema.StringAttribute {
	return stringField(true, false, false, false, "Website address", "Website address.", "**Website address**.").resource()
}

// OptionalCompanyField is an optional company name string attribute.
func OptionalCompanyField() schema.StringAttribute {
	return stringField(true, false, false, false, "Company", "Company name.", "**Company name**.").resource()
}

// OptionalDatabaseTypeField is an optional database type string attribute (f.text.type).
func OptionalDatabaseTypeField() schema.StringAttribute {
	return stringField(true, false, false, false, "Database type", "Database type (e.g. SQL).", "**Database type** (e.g. `SQL`).").resource()
}

// OptionalEmailField is an optional email address string attribute.
func OptionalEmailField() schema.StringAttribute {
	return stringField(true, false, false, false, "Email", "Email address.", "**Email address**.").resource()
}

// OptionalAddressRefField is an optional linked record UID string attribute.
func OptionalAddressRefField() schema.StringAttribute {
	return stringField(true, false, false, false, "Record UID", "Existing address record UID.", "Existing address record **UID**.").resource()
}

// OptionalCardRefField is an optional linked bankCard record UID string attribute.
func OptionalCardRefField() schema.StringAttribute {
	return stringField(true, false, false, false, "Record UID", "Existing payment card record UID.", "Existing payment card **UID**.").resource()
}

// OptionalCardholderNameField is an optional cardholder name string attribute.
func OptionalCardholderNameField() schema.StringAttribute {
	return stringField(true, false, false, false, "Cardholder name", "Cardholder name.", "**Cardholder name**.").resource()
}

// OptionalPinCodeField is an optional sensitive PIN code string attribute.
func OptionalPinCodeField() schema.StringAttribute {
	return stringField(true, false, true, false, "PIN code", "PIN code.", "**PIN code**.").resource()
}

// OptionalAccountNumberField is an optional sensitive account number string attribute.
func OptionalAccountNumberField() schema.StringAttribute {
	return stringField(true, false, true, false, "Account number", "Account number.", "**Account number**.").resource()
}

// OptionalHostnameField is an optional hostname or IP address string attribute.
func OptionalHostnameField() schema.StringAttribute {
	return stringField(true, false, false, false, "Hostname", "Hostname or IP address.", "**Hostname** or **IP address**.").resource()
}

// OptionalPortField is an optional port string attribute (digits only).
// Uses NumericStringValidator instead of the default min-length rule.
func OptionalPortField() schema.StringAttribute {
	return stringField(true, false, false, false, "port", "Port.", "**Port** (numeric string) eg: \"22\".").
		withValidators(utils.NumericStringValidator("Port", true)).
		resource()
}

// OptionalPassphraseField is an optional sensitive passphrase string attribute.
func OptionalPassphraseField() schema.StringAttribute {
	return stringField(true, false, true, false, "Passphrase", "Passphrase.", "**Passphrase**.").resource()
}

// OptionalPublicKeyField is an optional sensitive public key string attribute.
func OptionalPublicKeyField() schema.StringAttribute {
	return stringField(true, false, true, false, "Public key", "Public key.", "**Public key**.").resource()
}

// OptionalPrivateKeyField is an optional sensitive private key string attribute.
func OptionalPrivateKeyField() schema.StringAttribute {
	return stringField(true, false, true, false, "Private key", "Private key.", "**Private key**.").resource()
}

// RequiredSSIDField is a required network SSID string attribute.
func RequiredSSIDField() schema.StringAttribute {
	return stringField(false, true, false, false, "SSID", "Network SSID (network name).", "**Network SSID** (network name).").resource()
}

// OptionalWifiEncryptionField is an optional WiFi encryption type string attribute.
// allowed must list valid Keeper wifiEncryption values (e.g. wep, wpa, noEncryption).
func OptionalWifiEncryptionField(allowed []string) schema.StringAttribute {
	return stringField(true, false, false, false, "Encryption", "Encryption type.", "**Encryption type**.").
		withValidators(utils.StringOneOfValidator("Encryption", allowed, true)).
		resource()
}

// OptionalSSIDHiddenField is an optional bool for hidden SSID.
func OptionalSSIDHiddenField() schema.BoolAttribute {
	return schema.BoolAttribute{
		Optional:            true,
		Description:         "Whether the SSID is hidden.",
		MarkdownDescription: "Whether the SSID is hidden (not broadcast).",
	}
}

// ---------------------------------------------------------------------------
// Computed string/bool field helpers for generic vault record data sources.
// ---------------------------------------------------------------------------

// ComputedLoginField is a computed login/username string attribute for data sources.
func ComputedLoginField() dschema.StringAttribute {
	return stringField(false, false, false, true, "", "Username or login identifier.", "**Username** or **login** identifier.").dataSource()
}

// ComputedPasswordField is a computed sensitive password string attribute for data sources.
func ComputedPasswordField() dschema.StringAttribute {
	return stringField(false, false, true, true, "", "Password.", "**Password**.").dataSource()
}

// ComputedWebsiteAddressField is a computed website address string attribute for data sources.
func ComputedWebsiteAddressField() dschema.StringAttribute {
	return stringField(false, false, false, true, "", "Website address.", "**Website address**.").dataSource()
}

// ComputedCompanyField is a computed company name string attribute for data sources.
func ComputedCompanyField() dschema.StringAttribute {
	return stringField(false, false, false, true, "", "Company name.", "**Company name**.").dataSource()
}

// ComputedDatabaseTypeField is a computed database type string attribute for data sources.
func ComputedDatabaseTypeField() dschema.StringAttribute {
	return stringField(false, false, false, true, "", "Database type (e.g. SQL).", "**Database type** (e.g. `SQL`).").dataSource()
}

// ComputedEmailField is a computed email address string attribute for data sources.
func ComputedEmailField() dschema.StringAttribute {
	return stringField(false, false, false, true, "", "Email address.", "**Email address**.").dataSource()
}

// ComputedCardholderNameField is a computed cardholder name string attribute for data sources.
func ComputedCardholderNameField() dschema.StringAttribute {
	return stringField(false, false, false, true, "", "Cardholder name.", "**Cardholder name**.").dataSource()
}

// ComputedPinCodeField is a computed sensitive PIN code string attribute for data sources.
func ComputedPinCodeField() dschema.StringAttribute {
	return stringField(false, false, true, true, "", "PIN code.", "**PIN code**.").dataSource()
}

// ComputedAccountNumberField is a computed sensitive account number string attribute for data sources.
func ComputedAccountNumberField() dschema.StringAttribute {
	return stringField(false, false, true, true, "", "Account number.", "**Account number**.").dataSource()
}

// ComputedAddressRefField is a computed linked record UID string attribute for data sources.
func ComputedAddressRefField() dschema.StringAttribute {
	return stringField(false, false, false, true, "", "Linked record UID.", "**UID** of a linked record.").dataSource()
}

// ComputedCardRefField is a computed linked bankCard record UID string attribute for data sources.
func ComputedCardRefField() dschema.StringAttribute {
	return stringField(false, false, false, true, "", "Linked payment card record UID.", "**UID** of a linked payment card record.").dataSource()
}

// ComputedHostnameField is a computed hostname or IP address string attribute for data sources.
func ComputedHostnameField() dschema.StringAttribute {
	return stringField(false, false, false, true, "", "Hostname or IP address.", "**Hostname** or **IP address**.").dataSource()
}

// ComputedPortField is a computed port string attribute for data sources.
func ComputedPortField() dschema.StringAttribute {
	return stringField(false, false, false, true, "", "Port.", "**Port** (numeric string).").dataSource()
}

// ComputedPassphraseField is a computed sensitive passphrase string attribute for data sources.
func ComputedPassphraseField() dschema.StringAttribute {
	return stringField(false, false, true, true, "", "Passphrase.", "**Passphrase**.").dataSource()
}

// ComputedPublicKeyField is a computed sensitive public key string attribute for data sources.
func ComputedPublicKeyField() dschema.StringAttribute {
	return stringField(false, false, true, true, "", "Public key.", "**Public key**.").dataSource()
}

// ComputedPrivateKeyField is a computed sensitive private key string attribute for data sources.
func ComputedPrivateKeyField() dschema.StringAttribute {
	return stringField(false, false, true, true, "", "Private key.", "**Private key**.").dataSource()
}

// ComputedSSIDField is a computed network SSID string attribute for data sources.
func ComputedSSIDField() dschema.StringAttribute {
	return stringField(false, false, false, true, "", "Network SSID (network name).", "**Network SSID** (network name).").dataSource()
}

// ComputedWifiEncryptionField is a computed WiFi encryption type string attribute for data sources.
func ComputedWifiEncryptionField() dschema.StringAttribute {
	return stringField(false, false, false, true, "", "Encryption type.", "**Encryption type**.").dataSource()
}

// ComputedSSIDHiddenField is a computed bool for hidden SSID on data sources.
func ComputedSSIDHiddenField() dschema.BoolAttribute {
	return computedBoolField("Whether the SSID is hidden.", "Whether the **SSID** is hidden (not broadcast).")
}

// OptionalDateStringField stores Keeper date/birthDate/expirationDate as YYYY-MM-DD.
func OptionalDateStringField(name string) schema.StringAttribute {
	return stringField(true, false, false, false, name, "Date value (YYYY-MM-DD).", "**Date value** (`YYYY-MM-DD`).").
		withValidators(utils.DateStringValidator(name, true)).
		resource()
}

// ComputedDateStringField is a computed date string attribute for data sources (YYYY-MM-DD).
func ComputedDateStringField(desc, md string) dschema.StringAttribute {
	return stringField(false, false, false, true, "", desc, md).dataSource()
}

// OptionalSoftwareLicenseKeyField is an optional sensitive software license key string attribute.
func OptionalSoftwareLicenseKeyField() schema.StringAttribute {
	return stringField(true, false, true, false, "Software license key", "Software license key.", "**Software license key**.").resource()
}

// ComputedSoftwareLicenseKeyField is a computed sensitive software license key string attribute for data sources.
func ComputedSoftwareLicenseKeyField() dschema.StringAttribute {
	return stringField(false, false, true, true, "", "Software license key.", "**Software license key**.").dataSource()
}

// OptionalSecuredNoteField is an optional sensitive secured note string attribute.
func OptionalSecuredNoteField() schema.StringAttribute {
	return stringField(true, false, true, false, "Secured note", "Secured note content.", "**Secured note** content.").resource()
}

// ComputedSecuredNoteField is a computed sensitive secured note string attribute for data sources.
func ComputedSecuredNoteField() dschema.StringAttribute {
	return stringField(false, false, true, true, "", "Secured note content.", "**Secured note** content.").dataSource()
}

// NameDataSourceNestedSchema is the computed name nested object for data sources.
func NameDataSourceNestedSchema() dschema.SingleNestedAttribute {
	return dschema.SingleNestedAttribute{
		Computed:            true,
		Description:         "Person name (first, middle, last).",
		MarkdownDescription: "Person name (`name` field): first, middle, last.",
		Attributes: map[string]dschema.Attribute{
			"first":  stringField(false, false, false, true, "First name", "First name.", "**First name**.").dataSource(),
			"middle": stringField(false, false, false, true, "Middle name", "Middle name.", "**Middle name**.").dataSource(),
			"last":   stringField(false, false, false, true, "Last name", "Last name.", "**Last name**.").dataSource(),
		},
	}
}

// AddressDataSourceNestedSchema is the computed address nested object for data sources.
func AddressDataSourceNestedSchema() dschema.SingleNestedAttribute {
	return dschema.SingleNestedAttribute{
		Computed:            true,
		Description:         "Physical address.",
		MarkdownDescription: "Physical address (`address` field).",
		Attributes: map[string]dschema.Attribute{
			"street1": stringField(false, false, false, true, "", "Street line 1.", "Street line 1.").dataSource(),
			"street2": stringField(false, false, false, true, "", "Street line 2.", "Street line 2.").dataSource(),
			"city":    stringField(false, false, false, true, "", "City.", "City.").dataSource(),
			"state":   stringField(false, false, false, true, "", "State or province.", "State or province.").dataSource(),
			"zip":     stringField(false, false, false, true, "", "Postal code.", "Postal code.").dataSource(),
			"country": stringField(false, false, false, true, "", "ISO 3166-1 alpha-2 country code.", "ISO 3166-1 alpha-2 country code.").dataSource(),
		},
	}
}

// PaymentCardDataSourceNestedSchema is the computed payment card nested object for data sources.
func PaymentCardDataSourceNestedSchema() dschema.SingleNestedAttribute {
	return dschema.SingleNestedAttribute{
		Computed:            true,
		Description:         "Payment card details.",
		MarkdownDescription: "Payment card details (`paymentCard` field).",
		Attributes: map[string]dschema.Attribute{
			"card_number":          stringField(false, false, true, true, "", "Card number.", "**Card number**.").dataSource(),
			"card_expiration_date": stringField(false, false, false, true, "", "Expiration MM/YYYY.", "Expiration `MM/YYYY`.").dataSource(),
			"card_security_code":   stringField(false, false, true, true, "", "Security code.", "**Security code**.").dataSource(),
		},
	}
}

// BankAccountDataSourceNestedSchema is the computed bank account nested object for data sources.
func BankAccountDataSourceNestedSchema() dschema.SingleNestedAttribute {
	return dschema.SingleNestedAttribute{
		Computed:            true,
		Description:         "Bank account details.",
		MarkdownDescription: "Bank account details (`bankAccount` field).",
		Attributes: map[string]dschema.Attribute{
			"account_type":   stringField(false, false, false, true, "", "Checking, Savings, or Other.", "**Account type**.").dataSource(),
			"other_type":     stringField(false, false, false, true, "", "Description when account_type is Other.", "Description when **account_type** is **Other**.").dataSource(),
			"routing_number": stringField(false, false, false, true, "", "Routing number.", "Routing number.").dataSource(),
			"account_number": stringField(false, false, true, true, "", "Account number.", "**Account number**.").dataSource(),
		},
	}
}

// PhoneDataSourceListSchema is the computed phone list for data sources.
func PhoneDataSourceListSchema() dschema.ListNestedAttribute {
	return dschema.ListNestedAttribute{
		Computed:            true,
		Description:         "Phone numbers.",
		MarkdownDescription: "Phone numbers.",
		NestedObject: dschema.NestedAttributeObject{
			Attributes: map[string]dschema.Attribute{
				"region": stringField(false, false, false, true, "Region or country code", "Region or country code (e.g. US, +1).", "**Region or country code** (e.g. US, +1).").dataSource(),
				"number": stringField(false, false, false, true, "Phone number", "Phone number.", "**Phone number**.").dataSource(),
				"ext":    stringField(false, false, false, true, "Extension", "Extension.", "**Extension**.").dataSource(),
				"type":   stringField(false, false, false, true, "Phone type", "Phone type: Mobile, Home, or Work.", "**Phone type**: `Mobile`, `Home`, or `Work`.").dataSource(),
			},
		},
	}
}
