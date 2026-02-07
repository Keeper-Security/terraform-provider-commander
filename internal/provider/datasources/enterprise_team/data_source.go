// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseteam

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &EnterpriseTeamDataSource{}
var _ datasource.DataSourceWithConfigure = &EnterpriseTeamDataSource{}

type EnterpriseTeamDataSource struct {
	apiManager *api.ApiManager
}

func (d *EnterpriseTeamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_team"
}

func (d *EnterpriseTeamDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiManager, ok := req.ProviderData.(*api.ApiManager)
	if !ok {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			fmt.Sprintf("The provider was not configured correctly. Expected API manager, but got: %T. Please check your provider configuration.", req.ProviderData),
		)
		return
	}

	d.apiManager = apiManager
}

func (d *EnterpriseTeamDataSource) ensureApiManager() error {
	if d.apiManager == nil {
		return fmt.Errorf("API manager not configured")
	}
	return nil
}

func NewEnterpriseTeamDataSource() datasource.DataSource {
	return &EnterpriseTeamDataSource{}
}
