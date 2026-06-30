// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SharedAttributes returns the Contact resource attribute map shared between
// classic and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			Description:         IDDescription,
			MarkdownDescription: IDMarkdownDescription,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"title": schema.StringAttribute{
			Required:            true,
			Description:         TitleDescription,
			MarkdownDescription: TitleMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Title", 1, false),
			},
		},
		"notes": schema.StringAttribute{
			Optional:            true,
			Description:         NotesDescription,
			MarkdownDescription: NotesMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Notes", 0, true),
			},
		},
		"folder_location": schema.StringAttribute{
			Optional:            true,
			Description:         FolderDescription,
			MarkdownDescription: FolderMarkdownDescription,
			Validators: []validator.String{
				utils.StringMinLengthValidator("Folder", 1, true),
			},
		},
		"name":        nameResourceAttribute(),
		"company":     optionalStringField("Company", CompanyDescription, CompanyMarkdownDescription),
		"email":       optionalStringField("Email", EmailDescription, EmailMarkdownDescription),
		"phone":       phoneResourceAttribute(),
		"address_ref": refUIDField(AddressRefDescription, AddressRefMarkdownDescription),
		"custom":      customFieldResourceAttribute(),
	}
}

// SharedDataSourceAttributes returns computed Contact data source attributes
// shared between classic and new data sources. Callers add the lookup key
// (e.g. contact) and share-extension attributes separately.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return map[string]dschema.Attribute{
		"id": dschema.StringAttribute{
			Computed:            true,
			Description:         IDDescription,
			MarkdownDescription: IDMarkdownDescription,
		},
		"title": dschema.StringAttribute{
			Computed:            true,
			Description:         TitleDescription,
			MarkdownDescription: TitleMarkdownDescription,
		},
		"notes": dschema.StringAttribute{
			Computed:            true,
			Description:         DSNotesDescription,
			MarkdownDescription: DSNotesMarkdownDescription,
		},
		"folder_location": dschema.StringAttribute{
			Computed:            true,
			Description:         DSFolderDescription,
			MarkdownDescription: DSFolderMarkdownDescription,
		},
		"name":        nameDataSourceAttribute(),
		"company":     computedStringAttribute(CompanyDescription, CompanyMarkdownDescription),
		"email":       computedStringAttribute(EmailDescription, EmailMarkdownDescription),
		"phone":       phoneDataSourceAttribute(),
		"address_ref": computedStringAttribute(AddressRefDescription, AddressRefMarkdownDescription),
		"custom":      customFieldDataSourceAttribute(),
	}
}

func nameResourceAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Required:            true,
		Description:         NameDescription,
		MarkdownDescription: NameMarkdownDescription,
		Attributes: map[string]schema.Attribute{
			"first": schema.StringAttribute{
				Required:            true,
				Description:         FirstNameDescription,
				MarkdownDescription: FirstNameMarkdownDescription,
			},
			"middle": schema.StringAttribute{
				Optional:            true,
				Description:         MiddleNameDescription,
				MarkdownDescription: MiddleNameMarkdownDescription,
			},
			"last": schema.StringAttribute{
				Required:            true,
				Description:         LastNameDescription,
				MarkdownDescription: LastNameMarkdownDescription,
			},
		},
	}
}

func phoneResourceAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:            true,
		Description:         PhoneDescription,
		MarkdownDescription: PhoneMarkdownDescription,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"region": schema.StringAttribute{
					Optional:            true,
					Description:         PhoneRegionDescription,
					MarkdownDescription: PhoneRegionMarkdownDescription,
				},
				"number": schema.StringAttribute{
					Optional:            true,
					Description:         PhoneNumberDescription,
					MarkdownDescription: PhoneNumberMarkdownDescription,
				},
				"ext": schema.StringAttribute{
					Optional:            true,
					Description:         PhoneExtDescription,
					MarkdownDescription: PhoneExtMarkdownDescription,
				},
				"type": schema.StringAttribute{
					Optional:            true,
					Description:         PhoneTypeDescription,
					MarkdownDescription: PhoneTypeMarkdownDescription,
					Validators: []validator.String{
						utils.StringOneOfValidator("Phone type", []string{"Mobile", "Home", "Work"}, true),
					},
				},
			},
		},
	}
}

func nameDataSourceAttribute() dschema.SingleNestedAttribute {
	return dschema.SingleNestedAttribute{
		Computed:            true,
		Description:         NameDescription,
		MarkdownDescription: NameMarkdownDescription,
		Attributes: map[string]dschema.Attribute{
			"first": dschema.StringAttribute{
				Computed:            true,
				Description:         FirstNameDescription,
				MarkdownDescription: FirstNameMarkdownDescription,
			},
			"middle": dschema.StringAttribute{
				Computed:            true,
				Description:         MiddleNameDescription,
				MarkdownDescription: MiddleNameMarkdownDescription,
			},
			"last": dschema.StringAttribute{
				Computed:            true,
				Description:         LastNameDescription,
				MarkdownDescription: LastNameMarkdownDescription,
			},
		},
	}
}

func phoneDataSourceAttribute() dschema.ListNestedAttribute {
	return dschema.ListNestedAttribute{
		Computed:            true,
		Description:         DSPhoneDescription,
		MarkdownDescription: DSPhoneMarkdownDescription,
		NestedObject: dschema.NestedAttributeObject{
			Attributes: map[string]dschema.Attribute{
				"region": dschema.StringAttribute{
					Computed:            true,
					Description:         PhoneRegionDescription,
					MarkdownDescription: PhoneRegionMarkdownDescription,
				},
				"number": dschema.StringAttribute{
					Computed:            true,
					Description:         PhoneNumberDescription,
					MarkdownDescription: PhoneNumberMarkdownDescription,
				},
				"ext": dschema.StringAttribute{
					Computed:            true,
					Description:         PhoneExtDescription,
					MarkdownDescription: PhoneExtMarkdownDescription,
				},
				"type": dschema.StringAttribute{
					Computed:            true,
					Description:         PhoneTypeDescription,
					MarkdownDescription: PhoneTypeMarkdownDescription,
				},
			},
		},
	}
}

func customFieldResourceAttribute() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:            true,
		Description:         CustomDescription,
		MarkdownDescription: CustomMarkdownDescription,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:            true,
					Description:         "Keeper field type (e.g. text, email, secret, phone, name, date).",
					MarkdownDescription: "Keeper field type (e.g. `text`, `email`, `secret`, `phone`, `name`, `date`).",
				},
				"label": schema.StringAttribute{
					Required:            true,
					Description:         "Field label.",
					MarkdownDescription: "Field label.",
					Validators: []validator.String{
						utils.StringMinLengthValidator("Custom field label", 1, false),
					},
				},
				"value": schema.StringAttribute{
					Required:            true,
					Description:         "Field value; for complex types use jsonencode(JSON) matching the Keeper field schema.",
					MarkdownDescription: "Field value; for complex types use `jsonencode(JSON)` matching the [Keeper field schema](https://docs.keeper.io/en/keeperpam/secrets-manager/about/field-record-types).",
				},
				"sensitive": schema.BoolAttribute{
					Optional:            true,
					Computed:            true,
					Description:         "Whether to mark the value as sensitive in Terraform state display.",
					MarkdownDescription: "Whether to mark the value as sensitive in Terraform state display.",
					Default:             booldefault.StaticBool(false),
				},
			},
		},
	}
}

func customFieldDataSourceAttribute() dschema.ListNestedAttribute {
	return dschema.ListNestedAttribute{
		Computed:            true,
		Description:         DSCustomDescription,
		MarkdownDescription: DSCustomMarkdownDescription,
		NestedObject: dschema.NestedAttributeObject{
			Attributes: map[string]dschema.Attribute{
				"type": dschema.StringAttribute{
					Computed:            true,
					Description:         DSCustomTypeDescription,
					MarkdownDescription: DSCustomTypeMarkdownDescription,
				},
				"label": dschema.StringAttribute{
					Computed:            true,
					Description:         "Field label.",
					MarkdownDescription: "Field label.",
				},
				"value": dschema.StringAttribute{
					Computed:            true,
					Description:         DSCustomValueDescription,
					MarkdownDescription: DSCustomValueMarkdownDescription,
				},
				"sensitive": dschema.BoolAttribute{
					Computed:            true,
					Description:         DSCustomSensitiveDescription,
					MarkdownDescription: DSCustomSensitiveMarkdownDescription,
				},
			},
		},
	}
}

func optionalStringField(name, desc, md string) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:            true,
		Description:         desc,
		MarkdownDescription: md,
		Validators: []validator.String{
			utils.StringMinLengthValidator(name, 1, true),
		},
	}
}

func refUIDField(desc, md string) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:            true,
		Description:         desc,
		MarkdownDescription: md,
		Validators: []validator.String{
			utils.StringMinLengthValidator("Record UID", 1, true),
		},
	}
}

func computedStringAttribute(desc, md string) dschema.StringAttribute {
	return dschema.StringAttribute{
		Computed:            true,
		Description:         desc,
		MarkdownDescription: md,
	}
}
