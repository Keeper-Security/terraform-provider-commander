// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package contact

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func (d *ContactDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dschema.Schema{
		Description:         SchemaDescription,
		MarkdownDescription: SchemaMarkdownDescription,
		Attributes: map[string]dschema.Attribute{
			"contact": dschema.StringAttribute{
				Required:            true,
				Description:         "Contact record title or UID to look up.",
				MarkdownDescription: "Contact record **title** or **UID** to look up.",
			},
			"id": dschema.StringAttribute{
				Computed:            true,
				Description:         "The unique identifier (UID) of the contact record.",
				MarkdownDescription: "The unique identifier (**UID**) of the contact record.",
			},
			"title": dschema.StringAttribute{
				Computed:            true,
				Description:         "Record title.",
				MarkdownDescription: "Record title.",
			},
			"notes": dschema.StringAttribute{
				Computed:            true,
				Description:         "Notes on the record.",
				MarkdownDescription: "Notes on the record.",
			},
			"folder": dschema.StringAttribute{
				Computed:            true,
				Description:         "Folder path or UID where the record is stored.",
				MarkdownDescription: "Folder path or UID where the record is stored.",
			},
			"name": dschema.SingleNestedAttribute{
				Computed:            true,
				Description:         "Person name (first, middle, last).",
				MarkdownDescription: "Person name (`name` field): first, middle, last.",
				Attributes: map[string]dschema.Attribute{
					"first": dschema.StringAttribute{
						Computed:            true,
						Description:         "First name.",
						MarkdownDescription: "First name.",
					},
					"middle": dschema.StringAttribute{
						Computed:            true,
						Description:         "Middle name.",
						MarkdownDescription: "Middle name.",
					},
					"last": dschema.StringAttribute{
						Computed:            true,
						Description:         "Last name.",
						MarkdownDescription: "Last name.",
					},
				},
			},
			"company": dschema.StringAttribute{
				Computed:            true,
				Description:         "Company name.",
				MarkdownDescription: "Company name.",
			},
			"email": dschema.StringAttribute{
				Computed:            true,
				Description:         "Email address.",
				MarkdownDescription: "Email address.",
			},
			"phone": dschema.ListNestedAttribute{
				Computed:            true,
				Description:         "Phone numbers.",
				MarkdownDescription: "Phone numbers.",
				NestedObject: dschema.NestedAttributeObject{
					Attributes: map[string]dschema.Attribute{
						"region": dschema.StringAttribute{
							Computed:            true,
							Description:         "Region or country code.",
							MarkdownDescription: "Region or country code.",
						},
						"number": dschema.StringAttribute{
							Computed:            true,
							Description:         "Phone number.",
							MarkdownDescription: "Phone number.",
						},
						"ext": dschema.StringAttribute{
							Computed:            true,
							Description:         "Extension.",
							MarkdownDescription: "Extension.",
						},
						"type": dschema.StringAttribute{
							Computed:            true,
							Description:         "Phone type: Mobile, Home, or Work.",
							MarkdownDescription: "Phone type: `Mobile`, `Home`, or `Work`.",
						},
					},
				},
			},
			"address_ref": dschema.StringAttribute{
				Computed:            true,
				Description:         "Linked Address record UID.",
				MarkdownDescription: "UID of an `address` record linked via `addressRef`.",
			},
			"custom": dschema.ListNestedAttribute{
				Computed:            true,
				Description:         "Custom fields on the record.",
				MarkdownDescription: "Custom fields on the record.",
				NestedObject: dschema.NestedAttributeObject{
					Attributes: map[string]dschema.Attribute{
						"type": dschema.StringAttribute{
							Computed:            true,
							Description:         "Keeper field type.",
							MarkdownDescription: "Keeper field type.",
						},
						"label": dschema.StringAttribute{
							Computed:            true,
							Description:         "Field label.",
							MarkdownDescription: "Field label.",
						},
						"value": dschema.StringAttribute{
							Computed:            true,
							Description:         "Field value.",
							MarkdownDescription: "Field value.",
						},
						"sensitive": dschema.BoolAttribute{
							Computed:            true,
							Description:         "Whether the value is sensitive.",
							MarkdownDescription: "Whether the value is sensitive.",
						},
					},
				},
			},
			"share": dschema.MapNestedAttribute{
				Computed:            true,
				Description:         "Users with share permissions on this record (keyed by email).",
				MarkdownDescription: "Users with share permissions on this record (keyed by email).",
				NestedObject: dschema.NestedAttributeObject{
					Attributes: map[string]dschema.Attribute{
						"can_share": dschema.BoolAttribute{
							Computed:            true,
							Description:         "Whether the user can share this record.",
							MarkdownDescription: "Whether the user can share this record.",
						},
						"can_edit": dschema.BoolAttribute{
							Computed:            true,
							Description:         "Whether the user can edit this record.",
							MarkdownDescription: "Whether the user can edit this record.",
						},
					},
				},
			},
		},
	}
}
