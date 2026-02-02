package enterpiseuser

import "github.com/hashicorp/terraform-plugin-framework/types"

type EnterpriseUserResourceModel struct {
	Id             types.String `tfsdk:"id"` // Note: For now we are using string type, once we get id in commander cli response while creating user we will change to Int64 type
	Email          types.String `tfsdk:"email"`
	Name           types.String `tfsdk:"name"`
	JobTitle       types.String `tfsdk:"job_title"`
	Roles          types.Set    `tfsdk:"roles"`
	Teams          types.Set    `tfsdk:"teams"`
	Node           types.String `tfsdk:"node"`
	ManagedCompany types.String `tfsdk:"managed_company"`
}
