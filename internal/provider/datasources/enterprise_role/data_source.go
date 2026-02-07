package enterpriserole

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &EnterpriseRoleDataSource{}
var _ datasource.DataSourceWithConfigure = &EnterpriseRoleDataSource{}

type EnterpriseRoleDataSource struct {
	apiManager *api.ApiManager
}

func (d *EnterpriseRoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_role"
}

func (d *EnterpriseRoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EnterpriseRoleDataSource) ensureApiManager() error {
	if d.apiManager == nil {
		return fmt.Errorf("API manager not configured")
	}
	return nil
}

func NewEnterpriseRoleDataSource() datasource.DataSource {
	return &EnterpriseRoleDataSource{}
}
