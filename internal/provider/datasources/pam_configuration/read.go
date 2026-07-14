// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	"context"
	"fmt"
	"strings"

	commonpamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_configuration"
	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (d *PamConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PamConfigurationDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := d.EnsureApiManager(); err != nil {
		resp.Diagnostics.AddError(utils.ERR_MSG_PROVIDER_CONFIGURATION_ERROR, err.Error())
		return
	}

	configUID := strings.TrimSpace(data.PamConfiguration.ValueString())
	if configUID == "" {
		resp.Diagnostics.AddError(commonpamconfiguration.ErrSummaryReadFailed, "pam_configuration is empty")
		return
	}

	if err := utils.SyncDown(ctx, d.ApiManager); err != nil {
		resp.Diagnostics.AddError(utils.ErrSummarySyncDownFailed, err.Error())
		return
	}

	command := commonpamconfiguration.FetchPamConfigByUIDCommand(configUID)
	apiResp, err := d.ApiManager.ExecuteCommand(ctx, command, commonpamconfiguration.ErrOpFetchPamConfig)
	if err != nil {
		resp.Diagnostics.AddError(commonpamconfiguration.ErrSummaryReadFailed, err.Error())
		return
	}

	if apiResp == nil || apiResp.Data == nil {
		resp.Diagnostics.AddError(commonpamconfiguration.ErrSummaryReadFailed, fmt.Sprintf("PAM configuration %q not found or empty response", configUID))
		return
	}

	var apiData utils.PamConfigListResponse
	if err := utils.UnmarshalApiResponse(apiResp.Data, &apiData); err != nil {
		resp.Diagnostics.AddError(commonpamconfiguration.ErrSummaryReadFailed, err.Error())
		return
	}

	var state commonpamconfiguration.PamConfigurationResourceModel
	if err := commonpamconfiguration.MapPamConfigAPIResponseToModel(&state, &apiData); err != nil {
		resp.Diagnostics.AddError(commonpamconfiguration.ErrSummaryReadFailed, err.Error())
		return
	}

	data.PamConfiguration = types.StringValue(configUID)
	data.Id = types.StringValue(strings.TrimSpace(apiData.UID))
	data.Title = state.Title
	data.Environment = state.Environment
	data.Gateway = state.Gateway
	data.ApplicationFolder = state.ApplicationFolder
	data.Schedule = state.Schedule
	data.PortMapping = state.PortMapping
	data.Connections = state.Connections
	data.Tunneling = state.Tunneling
	data.Rotation = state.Rotation
	data.RemoteBrowserIsolation = state.RemoteBrowserIsolation
	data.ConnectionsRecording = state.ConnectionsRecording
	data.TypescriptRecording = state.TypescriptRecording
	data.AIThreatDetection = state.AIThreatDetection
	data.AITerminateSessionOnDetection = state.AITerminateSessionOnDetection
	data.LocalNetwork = state.LocalNetwork
	data.Aws = state.Aws
	data.Azure = state.Azure
	data.Domain = state.Domain
	data.Gcp = state.Gcp

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
