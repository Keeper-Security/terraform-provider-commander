// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package records

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
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
			Description:         "Optional notes on the record.",
			MarkdownDescription: "Optional notes on the record.",
			Validators: []validator.String{
				utils.StringMinLengthValidator("Notes", 0, true),
			},
		},
		"folder": schema.StringAttribute{
			Optional:            true,
			Description:         "Folder path or UID where the record is stored.",
			MarkdownDescription: "Folder path or UID where the record is stored.",
			Validators: []validator.String{
				utils.StringMinLengthValidator("Folder", 1, true),
			},
		},
	}
}

// BaseRecordBlocks returns the custom fields nested list block.
func BaseRecordBlocks() map[string]schema.Block {
	return map[string]schema.Block{
		"custom": CustomFieldBlockSchema(),
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
		Description:         "Phone numbers (type, number, region, ext).",
		MarkdownDescription: "Phone numbers; maps to Keeper `phone` field array.",
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
			},
			"port": schema.StringAttribute{
				Optional:            true,
				Description:         "Port (string per Keeper API).",
				MarkdownDescription: "Port (string per Keeper API).",
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
			"account_type": schema.StringAttribute{
				Optional:            true,
				Description:         "Checking, Savings, or Other.",
				MarkdownDescription: "`Checking`, `Savings`, or `Other`.",
				Validators: []validator.String{
					utils.StringOneOfValidator("Account type", []string{"Checking", "Savings", "Other"}, true),
				},
			},
			"other_type": schema.StringAttribute{
				Optional:            true,
				Description:         "Description when account_type is Other.",
				MarkdownDescription: "Description when `account_type` is `Other`.",
			},
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
		MarkdownDescription: "Security question and answer (`securityQuestion` field).",
		Attributes: map[string]schema.Attribute{
			"question": schema.StringAttribute{Optional: true, Description: "Question.", MarkdownDescription: "Question."},
			"answer": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         "Answer.",
				MarkdownDescription: "Answer.",
			},
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
			"public_key": schema.StringAttribute{
				Optional:            true,
				Description:         "Public key.",
				MarkdownDescription: "Public key.",
			},
			"private_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				Description:         "Private key.",
				MarkdownDescription: "Private key.",
			},
		},
	}
}

// OptionalStringField is a reusable optional string attribute with min length 1 when set.
func OptionalStringField(name, desc, md string) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:            true,
		Description:         desc,
		MarkdownDescription: md,
		Validators: []validator.String{
			utils.StringMinLengthValidator(name, 1, true),
		},
	}
}

// OptionalSensitiveStringField is like OptionalStringField but sensitive.
func OptionalSensitiveStringField(name, desc, md string) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:            true,
		Sensitive:           true,
		Description:         desc,
		MarkdownDescription: md,
		Validators: []validator.String{
			utils.StringMinLengthValidator(name, 1, true),
		},
	}
}

// OptionalDateStringField stores Keeper date/birthDate/expirationDate as RFC3339 or YYYY-MM-DD.
func OptionalDateStringField(name, desc, md string) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:            true,
		Description:         desc,
		MarkdownDescription: md,
		Validators: []validator.String{
			utils.StringMinLengthValidator(name, 1, true),
		},
	}
}

// RefUIDField is addressRef / cardRef target record UID.
func RefUIDField(desc, md string) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:            true,
		Description:         desc,
		MarkdownDescription: md,
		Validators: []validator.String{
			utils.StringMinLengthValidator("Record UID", 1, true),
		},
	}
}
