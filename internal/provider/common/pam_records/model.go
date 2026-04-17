package pamrecords

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CommonPamRecordsResourceModel struct {
	Id     types.String `tfsdk:"id"`
	Title  types.String `tfsdk:"title"`
	Notes  types.String `tfsdk:"notes"`
	Folder types.String `tfsdk:"folder"`
}

type CommonPamSettingsRotationResourceModel struct {
}

type CommonPamSettingsConnectionResourceModel struct {
}

// This is structure of portForward that we get from the API.
type CommonPamSettingsTunnelResourceModel struct {
	Enabled               types.Bool   `tfsdk:"enabled"`
	Port                  types.Int32  `tfsdk:"port"`
	ReUsePort             types.Bool   `tfsdk:"re_use_port"`
	UseSpecifiedLocalPort types.Bool   `tfsdk:"use_specified_local_port"`
	LocalPort             types.String `tfsdk:"local_port"`
}
type CommonPamRecordsFieldResourceModel struct {
	AllowSupplyHost types.Bool                                `tfsdk:"allow_supply_host"`
	Connection      *CommonPamSettingsConnectionResourceModel `tfsdk:"connection"`
	Tunnel          *CommonPamSettingsTunnelResourceModel     `tfsdk:"tunnel"`
}
