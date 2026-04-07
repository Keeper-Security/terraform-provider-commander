// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *SecretsManagerAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Creates and manages a Keeper Secrets Manager application.",
		MarkdownDescription: "Creates and manages a **Keeper Secrets Manager** application.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Application UID assigned by Keeper after creation.",
				MarkdownDescription: "**Application UID** assigned by Keeper after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Name of the Secrets Manager application.",
				MarkdownDescription: "**Name** of the Secrets Manager application.",
			},
			"shares": schema.SetNestedAttribute{
				Optional:            true,
				Description:         "Secrets (records or shared folders) shared with this application.",
				MarkdownDescription: "**Secrets** (records or shared folders) shared with this application.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"secret": schema.StringAttribute{
							Required:            true,
							Description:         "Record or shared folder UID to share with this application.",
							MarkdownDescription: "**Record or shared folder UID** to share with this application.",
						},
						"editable": schema.BoolAttribute{
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
							Description:         "Allow secrets to be editable by the client. Defaults to false.",
							MarkdownDescription: "Allow secrets to be **editable** by the client. Defaults to `false`.",
						},
					},
				},
			},
			"app_users": schema.SetAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				Description:         "Set of user emails to grant access to this application.",
				MarkdownDescription: "Set of **user emails** to grant access to this application.",
			},
		},
	}
}
