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
	Status         types.String `tfsdk:"status"`
	// Alias          types.Set    `tfsdk:"alias"` // NOTE: Not working in commander cli, so we are not supporting it for now
	// HideSharedFolders types.Bool `tfsdk:"hide_shared_folders"` // Need to check we have to implement it or not
}
