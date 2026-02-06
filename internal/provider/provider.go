// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterprisenodedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_node"
	managecompanydatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/manage_company"
	enterpiseuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterpise_user"
	enterprisenode "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_node"
	enterpriserole "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_role"
	enterpriseteam "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_team"
	managecompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/manage_company"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure CommanderProvider satisfies various provider interfaces.
var _ provider.Provider = &CommanderProvider{}
var _ provider.ProviderWithFunctions = &CommanderProvider{}
var _ provider.ProviderWithEphemeralResources = &CommanderProvider{}
var _ provider.ProviderWithActions = &CommanderProvider{}

// CommanderProvider defines the provider implementation.
type CommanderProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// CommanderProviderModel describes the provider data model.
type CommanderProviderModel struct {
	ServiceModeUrl    types.String `tfsdk:"service_mode_url"`
	ServiceModeApiKey types.String `tfsdk:"service_mode_api_key"`
}

func (p *CommanderProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "commander"
	resp.Version = p.version
}

func (p *CommanderProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"service_mode_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the running Keeper Commander Service Mode, for more information see [Keeper Commander Service Mode](https://docs.keeper.io/en/keeperpam/commander-cli/service-mode-rest-api#keeper-commander-service-mode)",
				Description:         "The URL of the running Keeper Commander Service Mode",
				Required:            true,
			},
			"service_mode_api_key": schema.StringAttribute{
				MarkdownDescription: "The API key for the running Keeper Commander Service Mode, for more information see [Keeper Commander Service Mode](https://docs.keeper.io/en/keeperpam/commander-cli/service-mode-rest-api#keeper-commander-service-mode)",
				Description:         "The API key for the running Keeper Commander Service Mode",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *CommanderProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CommanderProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if data.ServiceModeUrl.IsNull() || data.ServiceModeUrl.ValueString() == "" {
		resp.Diagnostics.AddError("Missing Service Mode URL", "The Service Mode URL is required and cannot be empty")
	}
	if data.ServiceModeApiKey.IsNull() || data.ServiceModeApiKey.ValueString() == "" {
		resp.Diagnostics.AddError("Missing Service Mode API Key", "The Service Mode API Key is required and cannot be empty")
	}

	// Return early if validation failed
	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP Client with request timeout
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	// Normalize the Service Mode URL to always end with "/api/v2/"
	serviceModeUrl := strings.TrimSuffix(data.ServiceModeUrl.ValueString(), "/")
	if before, found := strings.CutSuffix(serviceModeUrl, "/api/v2"); found {
		serviceModeUrl = before
	}
	processedServiceModeUrl := serviceModeUrl + "/api/v2"

	// Create ApiManager with configuration
	apiManager := &api.ApiManager{
		ServiceModeUrl:    processedServiceModeUrl,
		ServiceModeApiKey: data.ServiceModeApiKey.ValueString(),
		HttpClient:        httpClient,
		IsMspAccount:      false,
	}

	// Detect account type during provider configuration
	if err := apiManager.IsMspAccountType(ctx); err != nil {
		// If detection fails, default to Enterprise account
		// This prevents MSP commands from being called on Enterprise accounts
		apiManager.IsMspAccount = false
	}

	resp.DataSourceData = apiManager
	resp.ResourceData = apiManager
}

func (p *CommanderProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Add your resources here
		managecompany.NewManageCompanyResource,
		enterprisenode.NewEnterpriseNodeResource,
		enterpriseteam.NewEnterpriseTeamResource,
		enterpriserole.NewEnterpriseRoleResource,
		enterpiseuser.NewEnterpriseUserResource,
	}
}

func (p *CommanderProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *CommanderProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Add your data sources here
		managecompanydatasource.NewManageCompanyDataSource,
		enterprisenodedatasource.NewEnterpriseNodesDataSource,
	}
}

func (p *CommanderProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// Add your resources here
	}
}

func (p *CommanderProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{
		// Add your resources here
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CommanderProvider{
			version: version,
		}
	}
}
