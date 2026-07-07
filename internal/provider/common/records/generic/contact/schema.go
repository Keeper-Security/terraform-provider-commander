// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// SharedAttributes returns the Contact resource attribute map shared between
// classic and new resources. Callers add any share-extension attribute separately.
func SharedAttributes() map[string]schema.Attribute {
	return utils.MergeResourceAttributes(
		commonrecordsutils.BaseRecordAttributes(),
		map[string]schema.Attribute{
			"name":        nameResourceAttribute(),
			"company":     optionalStringField("Company", CompanyDescription, CompanyMarkdownDescription),
			"email":       optionalStringField("Email", EmailDescription, EmailMarkdownDescription),
			"phone":       phoneResourceAttribute(),
			"address_ref": refUIDField(AddressRefDescription, AddressRefMarkdownDescription),
			"custom":      commonrecordsutils.CustomFieldAttributeSchema(),
		},
	)
}

// SharedDataSourceAttributes returns computed Contact data source attributes
// shared between classic and new data sources. Callers add the lookup key
// (e.g. contact) and share-extension attributes separately.
func SharedDataSourceAttributes() map[string]dschema.Attribute {
	return utils.MergeDataSourceAttributes(
		commonrecordsutils.DataSourceBaseRecordAttributes(),
		map[string]dschema.Attribute{
			"name":        nameDataSourceAttribute(),
			"company":     computedStringAttribute(CompanyDescription, CompanyMarkdownDescription),
			"email":       computedStringAttribute(EmailDescription, EmailMarkdownDescription),
			"phone":       phoneDataSourceAttribute(),
			"address_ref": computedStringAttribute(AddressRefDescription, AddressRefMarkdownDescription),
			"custom":      commonrecordsutils.CustomFieldDataSourceAttributeSchema(),
		},
	)
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
