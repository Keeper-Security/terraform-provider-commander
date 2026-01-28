package enterprisenode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *EnterpriseNodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					nameValidator{},
				},
			},
			"parent": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					parentValidator{},
				},
			},
			"wipe_out": schema.BoolAttribute{
				Optional: true,
			},
			"toggle_isolated": schema.BoolAttribute{
				Optional: true,
			},
			// "logo_file": schema.StringAttribute{
			// 	Optional: true,
			// },
			"managed_company": schema.StringAttribute{
				Optional:   true,
				Validators: []validator.String{
					// managedCompanyValidator{},
				},
			},
		},
	}
}
