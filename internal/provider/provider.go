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
	nonsharedfolderdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/folders/classic_folders/non_shared_folder"
	classicsharedfolderdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/folders/classic_folders/shared_folder"
	newfolderdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/folders/new_folder"
	managedcompanydatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/managed_company"
	pamconfigurationdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/pam_configuration"
	classicaddressdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/address"
	classicbankaccountdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/bank_account"
	classicbirthcertificatedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/birth_certificate"
	classiccontactdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/contact"
	classicdatabasedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/database"
	classicdriverlicensedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/driver_license"
	classichealthinsurancedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/health_insurance"
	classiclogindatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/login"
	classicmembershipdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/membership"
	classicpassportdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/passport"
	classicpaymentcarddatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/payment_card"
	classicsaasconfigurationdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/saas_configuration"
	classicsecurenotedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/secure_note"
	classicserverdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/server"
	classicsoftwarelicensedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/software_license"
	classicsshkeysdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/ssh_keys"
	classicssncarddatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/ssn_card"
	classicwifidatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/generic/wifi"
	classicpamdatabasedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/pam/pam_database"
	classicpamdirectorydatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/pam/pam_directory"
	classicpammachinedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/pam/pam_machine"
	classicpamremotebrowserdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/pam/pam_remote_browser"
	classicpamuserdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/classic/pam/pam_user"
	newdatabasedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/new/generic/database"
	newlogindatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/new/generic/login"
	newsaasconfigurationdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/new/generic/saas_configuration"
	newserverdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/new/generic/server"
	newwifidatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/new/generic/wifi"
	newpamdatabasedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/new/pam/pam_database"
	newpamdirectorydatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/new/pam/pam_directory"
	newpammachinedatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/new/pam/pam_machine"
	newpamremotebrowserdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/new/pam/pam_remote_browser"
	newpamuserdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/records/new/pam/pam_user"
	secretsmanagerdatasource "github.com/Keeper-Security/terraform-provider-commander/internal/provider/datasources/secrets_manager"
	enterprisenode "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_node"
	enterprisepush "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_push"
	enterpriserole "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_role"
	enterprisescim "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_scim"
	enterprisescimpush "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_scim_push"
	enterpriseteam "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_team"
	enterpriseuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/enterprise_user"
	epmpolicy "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/epm_policy"
	nonsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/classic_folders/non_shared_folder"
	classicsharedfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/classic_folders/shared_folder"
	newfolder "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/folders/new_folder"
	managedcompany "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/managed_company"
	pamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/pam_configuration"
	classicaddress "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/address"
	classicbankaccount "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/bank_account"
	classicbirthcertificate "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/birth_certificate"
	classiccontact "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/contact"
	classicdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/database"
	classicdriverlicense "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/driver_license"
	classichealthinsurance "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/health_insurance"
	classiclogin "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/login"
	classicmembership "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/membership"
	classicpassport "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/passport"
	classicpaymentcard "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/payment_card"
	classicsaasconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/saas_configuration"
	classicsecurenote "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/secure_note"
	classicserver "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/server"
	classicsoftwarelicense "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/software_license"
	classicsshkeys "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/ssh_keys"
	classicssncard "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/ssn_card"
	classicwifi "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/generic/wifi"
	classicpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/pam/pam_database"
	classicpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/pam/pam_directory"
	classicpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/pam/pam_machine"
	classicpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/pam/pam_remote_browser"
	classicpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/classic/pam/pam_user"
	newdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/new/generic/database"
	newlogin "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/new/generic/login"
	newsaasconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/new/generic/saas_configuration"
	newserver "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/new/generic/server"
	newwifi "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/new/generic/wifi"
	newpamdatabase "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/new/pam/pam_database"
	newpamdirectory "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/new/pam/pam_directory"
	newpammachine "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/new/pam/pam_machine"
	newpamremotebrowser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/new/pam/pam_remote_browser"
	newpamuser "github.com/Keeper-Security/terraform-provider-commander/internal/provider/resources/records/new/pam/pam_user"
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
		classicpamremotebrowser.NewPamRemoteBrowserResource,
		classicpamuser.NewPamUserResource,
		classicpamdatabase.NewPamDatabaseResource,
		classicpamdirectory.NewPamDirectoryResource,
		classicpammachine.NewPamMachineResource,
		classiclogin.NewLoginResource,
		newlogin.NewLoginResource,
		newdatabase.NewDatabaseResource,
		classicwifi.NewWifiResource,
		classiccontact.NewContactResource,
		classicaddress.NewAddressResource,
		classicpaymentcard.NewPaymentCardResource,
		classicbankaccount.NewBankAccountResource,
		classicmembership.NewMembershipResource,
		classichealthinsurance.NewHealthInsuranceResource,
		classicdriverlicense.NewDriverLicenseResource,
		classicpassport.NewPassportResource,
		classicssncard.NewSsnCardResource,
		classicbirthcertificate.NewBirthCertificateResource,
		classicsshkeys.NewSshKeysResource,
		classicsaasconfiguration.NewSaasConfigurationResource,
		classicserver.NewServerResource,
		classicdatabase.NewDatabaseResource,
		classicsoftwarelicense.NewSoftwareLicenseResource,
		classicsecurenote.NewSecureNoteResource,
		classicsharedfolder.NewClassicSharedFolderResource,
		newfolder.NewNewFolderResource,
		secretsmanager.NewSecretsManagerAppResource,
		nonsharedfolder.NewNonSharedFolderResource,
		newpamremotebrowser.NewPamRemoteBrowserResource,
		newpamuser.NewPamUserResource,
		newpamdatabase.NewPamDatabaseResource,
		newpamdirectory.NewPamDirectoryResource,
		newpammachine.NewPamMachineResource,
		newsaasconfiguration.NewSaasConfigurationResource,
		newwifi.NewWifiResource,
		newserver.NewServerResource,
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
		classicpamremotebrowserdatasource.NewPamRemoteBrowserDataSource,
		classicpamuserdatasource.NewPamUserDataSource,
		classiclogindatasource.NewLoginDataSource,
		newlogindatasource.NewLoginDataSource,
		newdatabasedatasource.NewDatabaseDataSource,
		classiccontactdatasource.NewContactDataSource,
		classicaddressdatasource.NewAddressDataSource,
		classicpaymentcarddatasource.NewPaymentCardDataSource,
		classicbankaccountdatasource.NewBankAccountDataSource,
		classicmembershipdatasource.NewMembershipDataSource,
		classichealthinsurancedatasource.NewHealthInsuranceDataSource,
		classicdriverlicensedatasource.NewDriverLicenseDataSource,
		classicpassportdatasource.NewPassportDataSource,
		classicssncarddatasource.NewSsnCardDataSource,
		classicbirthcertificatedatasource.NewBirthCertificateDataSource,
		classicwifidatasource.NewWifiDataSource,
		classicsshkeysdatasource.NewSshKeysDataSource,
		classicsaasconfigurationdatasource.NewSaasConfigurationDataSource,
		classicserverdatasource.NewServerDataSource,
		classicdatabasedatasource.NewDatabaseDataSource,
		classicsoftwarelicensedatasource.NewSoftwareLicenseDataSource,
		classicsecurenotedatasource.NewSecureNoteDataSource,
		pamconfigurationdatasource.NewPamConfigurationDataSource,
		classicpamdatabasedatasource.NewPamDatabaseDataSource,
		classicpamdirectorydatasource.NewPamDirectoryDataSource,
		classicpammachinedatasource.NewPamMachineDataSource,
		classicsharedfolderdatasource.NewClassicSharedFolderDataSource,
		newfolderdatasource.NewNewFolderDataSource,
		nonsharedfolderdatasource.NewNonSharedFolderDataSource,
		newpamremotebrowserdatasource.NewPamRemoteBrowserDataSource,
		newpamuserdatasource.NewPamUserDataSource,
		newpamdatabasedatasource.NewPamDatabaseDataSource,
		newpamdirectorydatasource.NewPamDirectoryDataSource,
		newpammachinedatasource.NewPamMachineDataSource,
		newsaasconfigurationdatasource.NewSaasConfigurationDataSource,
		newwifidatasource.NewWifiDataSource,
		newserverdatasource.NewServerDataSource,
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
