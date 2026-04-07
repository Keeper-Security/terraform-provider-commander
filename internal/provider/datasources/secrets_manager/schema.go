// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *SecretsManagerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to look up a Keeper Secrets Manager application by name or UID.",
		MarkdownDescription: "Use this data source to look up a **Keeper Secrets Manager** application by **name** or **UID**.",
		Attributes: map[string]schema.Attribute{
			"application": schema.StringAttribute{
				Required:            true,
				Description:         "Application name or UID to look up.",
				MarkdownDescription: "Application **name** or **UID** to look up.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Application UID.",
				MarkdownDescription: "**Application UID**.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the application.",
				MarkdownDescription: "**Name** of the application.",
			},
			"shares": schema.SetNestedAttribute{
				Computed:            true,
				Description:         "Secrets (records or shared folders) shared with this application.",
				MarkdownDescription: "**Secrets** (records or shared folders) shared with this application.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"secret": schema.StringAttribute{
							Computed:            true,
							Description:         "Record or shared folder UID.",
							MarkdownDescription: "**Record or shared folder UID**.",
						},
						"editable": schema.BoolAttribute{
							Computed:            true,
							Description:         "Whether the secret is editable by the client.",
							MarkdownDescription: "Whether the secret is **editable** by the client.",
						},
					},
				},
			},
			"app_users": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				Description:         "Set of user emails with access to this application.",
				MarkdownDescription: "Set of **user emails** with access to this application.",
			},
		},
	}
}
