// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	enterprisenodedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_node"
	enterpriseroledatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_role"
	enterprisescimdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_scim"
	enterpriseteamdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_team"
	enterpriseuserdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/enterprise_user"
	epmpolicydatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/epm_policy"
	classicsharedfolderdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/folders/classic_folders/shared_folder"
	newfolderdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/folders/new_folder"
	managedcompanydatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/managed_company"
	pamconfigurationdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/pam_configuration"
	pamdatabasedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic_records/pam_records/pam_database"
	pamdirectorydatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic_records/pam_records/pam_directory"
	pammachinedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic_records/pam_records/pam_machine"
	pamremotebrowserdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic_records/pam_records/pam_remote_browser"
	pamuserdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic_records/pam_records/pam_user"
	secretsmanagerdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/secrets_manager"
	enterprisenode "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_node"
	enterprisepush "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_push"
	enterpriserole "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_role"
	enterprisescim "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_scim"
	enterprisescimpush "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_scim_push"
	enterpriseteam "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_team"
	enterpriseuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_user"
	epmpolicy "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/epm_policy"
	classicsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/classic_folders/shared_folder"
	newfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/new_folder"
	managedcompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/managed_company"
	pamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/pam_configuration"
	pamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic_records/pam_records/pam_database"
	pamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic_records/pam_records/pam_directory"
	pammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic_records/pam_records/pam_machine"
	pamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic_records/pam_records/pam_remote_browser"
	pamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic_records/pam_records/pam_user"
	secretsmanager "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/secrets_manager"
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
	Timeout           types.Int64  `tfsdk:"timeout"`
}

func (p *CommanderProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "commander"
	resp.Version = p.version
}

func (p *CommanderProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage Keeper enterprise and MSP configuration as code via the Commander Service Mode API. See the detailed documentation https://docs.keeper.io/en/keeperpam/secrets-manager/integrations/terraform-provider-commander for information about features, prerequisites, setup and installation, and examples.",
		MarkdownDescription: "Manage Keeper **enterprise** and **MSP** configuration as code via the [Commander Service Mode API](https://docs.keeper.io/en/keeperpam/commander-cli/service-mode-rest-api).\n\n" + "-> <i>**New to the Commander provider?**</i> See the detailed [documentation](https://docs.keeper.io/en/keeperpam/secrets-manager/integrations/terraform-provider-commander) for information about features, prerequisites, setup and installation.",
		Attributes: map[string]schema.Attribute{
			"service_mode_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the running Keeper Commander Service Mode. Can also be set via the `COMMANDER_SERVICE_MODE_URL` environment variable.",
				Description:         "The URL of the running Keeper Commander Service Mode. Can also be set via the COMMANDER_SERVICE_MODE_URL environment variable.",
				Optional:            true,
			},
			"service_mode_api_key": schema.StringAttribute{
				MarkdownDescription: "The API key for the running Keeper Commander Service Mode. Can also be set via the `COMMANDER_SERVICE_MODE_API_KEY` environment variable.",
				Description:         "The API key for the running Keeper Commander Service Mode. Can also be set via the COMMANDER_SERVICE_MODE_API_KEY environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "Timeout in seconds for HTTP requests to the Commander Service Mode and for waiting on async command results. If this value is **not provided or is set to 0 or less**, the provider will use the default timeout of `60` seconds.",
				Description:         "Timeout in seconds for HTTP requests to the Commander Service Mode and for waiting on async command results. If this value is not provided or is set to 0 or less, the provider will use the default timeout of 60 seconds.",
				Optional:            true,
			},
		},
	}
}

const (
	envServiceModeURL    = "COMMANDER_SERVICE_MODE_URL"
	envServiceModeAPIKey = "COMMANDER_SERVICE_MODE_API_KEY"
)

func (p *CommanderProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CommanderProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Resolve URL and API key: config takes precedence over environment variables.
	serviceModeUrl := data.ServiceModeUrl.ValueString()
	if serviceModeUrl == "" {
		serviceModeUrl = os.Getenv(envServiceModeURL)
	}
	serviceModeApiKey := data.ServiceModeApiKey.ValueString()
	if serviceModeApiKey == "" {
		serviceModeApiKey = os.Getenv(envServiceModeAPIKey)
	}

	if serviceModeUrl == "" {
		resp.Diagnostics.AddError("Missing Service Mode URL", "The Service Mode URL is required. Set the service_mode_url attribute or the COMMANDER_SERVICE_MODE_URL environment variable.")
	}
	if serviceModeApiKey == "" {
		resp.Diagnostics.AddError("Missing Service Mode API Key", "The Service Mode API Key is required. Set the service_mode_api_key attribute or the COMMANDER_SERVICE_MODE_API_KEY environment variable.")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Timeout: config value or default 60 seconds
	timeoutSec := int64(60)
	if !data.Timeout.IsNull() && !data.Timeout.IsUnknown() {
		timeoutSec = data.Timeout.ValueInt64()
	}
	if timeoutSec <= 0 {
		timeoutSec = 60
	}
	requestTimeout := time.Duration(timeoutSec) * time.Second

	// Create HTTP Client with request timeout (same as poll timeout by default)
	httpClient := &http.Client{
		Timeout: requestTimeout,
	}

	// Normalize the Service Mode URL to always end with "/api/v2/"
	serviceModeUrl = strings.TrimSuffix(serviceModeUrl, "/")
	if before, found := strings.CutSuffix(serviceModeUrl, "/api/v2"); found {
		serviceModeUrl = before
	}
	processedServiceModeUrl := serviceModeUrl + "/api/v2"

	// Create ApiManager with configuration
	apiManager := &api.ApiManager{
		ServiceModeUrl:    processedServiceModeUrl,
		ServiceModeApiKey: serviceModeApiKey,
		HttpClient:        httpClient,
		IsMspAccount:      false,
		RequestTimeout:    requestTimeout,
	}

	// Detect account type during provider configuration
	if err := apiManager.IsMspAccountType(ctx); err != nil {
		// If detection fails, default to Enterprise account
		// This prevents MSP commands from being called on Enterprise accounts
		apiManager.IsMspAccount = false
	}
	// Ensure we start with a clean context (MSP). Avoids a no-MC op wrongly calling switch-to-msp
	// when currentContext was left set from a previous run or from plan/apply ordering.
	apiManager.SetCurrentContext("")

	resp.DataSourceData = apiManager
	resp.ResourceData = apiManager
}

func (p *CommanderProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Add your resources here
		managedcompany.NewManagedCompanyResource,
		enterprisenode.NewEnterpriseNodeResource,
		enterprisepush.NewEnterprisePushResource,
		epmpolicy.NewEpmPolicyResource,
		enterprisescimpush.NewEnterpriseScimPushResource,
		enterprisescim.NewEnterpriseScimResource,
		enterpriseteam.NewEnterpriseTeamResource,
		enterpriserole.NewEnterpriseRoleResource,
		enterpriseuser.NewEnterpriseUserResource,
		pamconfiguration.NewPamConfigurationResource,
		pamremotebrowser.NewPamRemoteBrowserResource,
		pamuser.NewPamUserResource,
		pamdatabase.NewPamDatabaseResource,
		pamdirectory.NewPamDirectoryResource,
		pammachine.NewPamMachineResource,
		classicsharedfolder.NewClassicSharedFolderResource,
		newfolder.NewNewFolderResource,
		secretsmanager.NewSecretsManagerAppResource,
	}
}

func (p *CommanderProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *CommanderProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Add your data sources here
		managedcompanydatasource.NewManagedCompanyDataSource,
		enterprisenodedatasource.NewEnterpriseNodesDataSource,
		enterpriseroledatasource.NewEnterpriseRoleDataSource,
		enterprisescimdatasource.NewEnterpriseScimDataSource,
		enterpriseteamdatasource.NewEnterpriseTeamDataSource,
		enterpriseuserdatasource.NewEnterpriseUserDataSource,
		secretsmanagerdatasource.NewSecretsManagerDataSource,
		epmpolicydatasource.NewEpmPolicyDataSource,
		pamremotebrowserdatasource.NewPamRemoteBrowserDataSource,
		pamuserdatasource.NewPamUserDataSource,
		pamconfigurationdatasource.NewPamConfigurationDataSource,
		pamdatabasedatasource.NewPamDatabaseDataSource,
		pamdirectorydatasource.NewPamDirectoryDataSource,
		pammachinedatasource.NewPamMachineDataSource,
		classicsharedfolderdatasource.NewClassicSharedFolderDataSource,
		newfolderdatasource.NewNewFolderDataSource,
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
