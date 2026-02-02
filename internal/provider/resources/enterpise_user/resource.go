package enterpiseuser

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &EnterpriseUserResource{}
var _ resource.ResourceWithConfigure = &EnterpriseUserResource{}

type EnterpriseUserResource struct {
	apiManager *api.ApiManager
}

func (r *EnterpriseUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enterprise_user"
}

func (r *EnterpriseUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.apiManager = apiManager
}

// ensureApiManager validates that apiManager is configured and returns an error if not
func (r *EnterpriseUserResource) ensureApiManager() error {
	if r.apiManager == nil {
		return fmt.Errorf("the Keeper Commander provider is not properly configured. Please ensure the provider is set up with valid service_mode_url and service_mode_api_key")
	}
	return nil
}

func NewEnterpriseUserResource() resource.Resource {
	return &EnterpriseUserResource{}
}
