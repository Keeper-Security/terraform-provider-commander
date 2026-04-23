// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamuser

import "github.com/hashicorp/terraform-plugin-framework/types"

type PamUserDataSourceModel struct {
	RecordUID         types.String                      `tfsdk:"record_uid"`
	Id                types.String                      `tfsdk:"id"`
	Title             types.String                      `tfsdk:"title"`
	Login             types.String                      `tfsdk:"login"`
	Password          types.String                      `tfsdk:"password"`
	Folder            types.String                      `tfsdk:"folder"`
	Notes             types.String                      `tfsdk:"notes"`
	DistinguishedName types.String                      `tfsdk:"distinguished_name"`
	PrivatePEMKey     types.String                      `tfsdk:"private_pem_key"`
	ConnectDatabase   types.String                      `tfsdk:"connect_database"`
	Managed           types.Bool                        `tfsdk:"managed"`
	RotationSettings  *PamUserDataSourceRotationSettings `tfsdk:"rotation_settings"`
}

type PamUserDataSourceRotationSettings struct {
	RotationProfile types.String `tfsdk:"rotation_profile"`
	Configuration   types.String `tfsdk:"configuration"`
	IamAadConfig    types.String `tfsdk:"iam_aad_config"`
	Resource        types.String `tfsdk:"resource"`
	AdminUser       types.String `tfsdk:"admin_user"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	ScheduleCron    types.String `tfsdk:"schedule_cron"`
	ScheduleJSON    types.String `tfsdk:"schedule_json"`
	OnDemand        types.Bool   `tfsdk:"on_demand"`
	ScheduleConfig  types.Bool   `tfsdk:"schedule_config"`
	Complexity      types.String `tfsdk:"complexity"`
}
