// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterpriseuser

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &EnterpriseUserDataSource{}
var _ datasource.DataSourceWithConfigure = &EnterpriseUserDataSource{}

type EnterpriseUserDataSource struct {
	apiManager *api.ApiManager
}

func (d *EnterpriseUserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_user"
}

func (d *EnterpriseUserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EnterpriseUserDataSource) ensureApiManager() error {
	if d.apiManager == nil {
		return fmt.Errorf("API manager not configured")
	}
	return nil
}

func NewEnterpriseUserDataSource() datasource.DataSource {
	return &EnterpriseUserDataSource{}
}
