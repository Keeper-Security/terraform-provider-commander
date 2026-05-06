// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	"context"
	"fmt"
	"strings"

	commonpamconfiguration "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/pam_configuration"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// environmentValidator ensures `environment` is one of ValidEnvironments.
type environmentValidator struct{}

func (environmentValidator) Description(_ context.Context) string {
	return fmt.Sprintf("must be one of: %s", strings.Join(ValidEnvironments, ", "))
}

func (environmentValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Must be one of: %s.", ValidEnvironmentsMarkdown())
}

func (v environmentValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	val := req.ConfigValue.ValueString()
	for _, allowed := range ValidEnvironments {
		if val == allowed {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid environment",
		fmt.Sprintf("must be one of: %s", strings.Join(ValidEnvironments, ", ")),
	)
}

// environmentBlocksValidator ensures exactly one environment block matching
// `environment` is provided, and that the required fields within it are set.
// With SingleNestedBlock, absent blocks are decoded as non-nil structs with
// all-null fields, so we detect presence by checking for any non-null attribute.
type environmentBlocksValidator struct{}

func (environmentBlocksValidator) Description(_ context.Context) string {
	return "Exactly one of local_network, aws, azure, domain, or gcp must be set and must match `environment`."
}

func (environmentBlocksValidator) MarkdownDescription(_ context.Context) string {
	return "Exactly one of `local_network`, `aws`, `azure`, `domain`, or `gcp` must be set and must match `environment`."
}

func isStringSet(v types.String) bool {
	return !v.IsNull() && !v.IsUnknown() && strings.TrimSpace(v.ValueString()) != ""
}

func isBoolSet(v types.Bool) bool {
	return !v.IsNull() && !v.IsUnknown()
}

func isLocalNetworkPresent(m *commonpamconfiguration.PamLocalNetworkModel) bool {
	if m == nil {
		return false
	}
	return isStringSet(m.NetworkId) || isStringSet(m.NetworkCidr)
}

func isAwsPresent(m *commonpamconfiguration.PamAwsModel) bool {
	if m == nil {
		return false
	}
	return isStringSet(m.AwsId) || isStringSet(m.AccessKeyId) || isStringSet(m.AccessSecretKey) ||
		(!m.RegionNames.IsNull() && !m.RegionNames.IsUnknown())
}

func isAzurePresent(m *commonpamconfiguration.PamAzureModel) bool {
	if m == nil {
		return false
	}
	return isStringSet(m.AzureId) || isStringSet(m.ClientId) || isStringSet(m.ClientSecret) ||
		isStringSet(m.SubscriptionId) || isStringSet(m.TenantId) ||
		(!m.ResourceGroups.IsNull() && !m.ResourceGroups.IsUnknown())
}

func isDomainPresent(m *commonpamconfiguration.PamDomainModel) bool {
	if m == nil {
		return false
	}
	return isStringSet(m.DomainId) || isStringSet(m.DomainHostname) || isStringSet(m.DomainPort) ||
		isBoolSet(m.DomainUseSsl) || isBoolSet(m.DomainScanDcCidr) ||
		isStringSet(m.DomainNetworkCidr) || isStringSet(m.DomainAdmin)
}

func isGcpPresent(m *commonpamconfiguration.PamGcpModel) bool {
	if m == nil {
		return false
	}
	return isStringSet(m.GcpId) || isStringSet(m.ServiceAccountKey) ||
		isStringSet(m.GoogleAdminEmail) || isStringSet(m.GcpRegion)
}

func (environmentBlocksValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config commonpamconfiguration.PamConfigurationResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Environment.IsNull() || config.Environment.IsUnknown() {
		return
	}

	env := config.Environment.ValueString()
	expected := expectedBlockName(env)
	if expected == "" {
		return
	}

	present := []string{}
	if isLocalNetworkPresent(config.LocalNetwork) {
		present = append(present, "local_network")
	}
	if isAwsPresent(config.Aws) {
		present = append(present, "aws")
	}
	if isAzurePresent(config.Azure) {
		present = append(present, "azure")
	}
	if isDomainPresent(config.Domain) {
		present = append(present, "domain")
	}
	if isGcpPresent(config.Gcp) {
		present = append(present, "gcp")
	}

	if len(present) > 1 {
		resp.Diagnostics.AddError(
			"Conflicting environment configuration blocks",
			fmt.Sprintf("Only one of local_network, aws, azure, domain, or gcp may be set; found: %s", strings.Join(present, ", ")),
		)
		return
	}
	if len(present) == 1 && present[0] != expected {
		resp.Diagnostics.AddError(
			"Environment does not match configuration block",
			fmt.Sprintf("environment is %q but %q is set; use the %q block instead, and remove %q.", env, present[0], expected, present[0]),
		)
		return
	}

	validateRequiredBlockFields(env, &config, resp)
}

func validateRequiredBlockFields(env string, config *commonpamconfiguration.PamConfigurationResourceModel, resp *resource.ValidateConfigResponse) {
	addMissing := func(block, field string) {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			fmt.Sprintf("%q is required inside the %q block when environment is %q.", field, block, env),
		)
	}

	switch env {
	case commonpamconfiguration.EnvAWS:
		if config.Aws == nil || !isStringSet(config.Aws.AwsId) {
			addMissing("aws", "aws_id")
		}
	case commonpamconfiguration.EnvAzure:
		m := config.Azure
		if m == nil || !isStringSet(m.AzureId) {
			addMissing("azure", "azure_id")
		}
		if m == nil || !isStringSet(m.ClientId) {
			addMissing("azure", "client_id")
		}
		if m == nil || !isStringSet(m.ClientSecret) {
			addMissing("azure", "client_secret")
		}
		if m == nil || !isStringSet(m.SubscriptionId) {
			addMissing("azure", "subscription_id")
		}
		if m == nil || !isStringSet(m.TenantId) {
			addMissing("azure", "tenant_id")
		}
	case commonpamconfiguration.EnvDomain:
		m := config.Domain
		if m == nil || !isStringSet(m.DomainId) {
			addMissing("domain", "domain_id")
		}
		if m == nil || !isStringSet(m.DomainHostname) {
			addMissing("domain", "domain_hostname")
		}
		if m == nil || !isStringSet(m.DomainPort) {
			addMissing("domain", "domain_port")
		}
		if m == nil || !isBoolSet(m.DomainUseSsl) {
			addMissing("domain", "domain_use_ssl")
		}
		if m == nil || !isStringSet(m.DomainAdmin) {
			addMissing("domain", "domain_admin")
		}
	case commonpamconfiguration.EnvGCP:
		m := config.Gcp
		if m == nil || !isStringSet(m.GcpId) {
			addMissing("gcp", "gcp_id")
		}
		if m == nil || !isStringSet(m.ServiceAccountKey) {
			addMissing("gcp", "service_account_key")
		}
	}
}

func expectedBlockName(env string) string {
	switch env {
	case commonpamconfiguration.EnvLocal:
		return "local_network"
	case commonpamconfiguration.EnvAWS:
		return "aws"
	case commonpamconfiguration.EnvAzure:
		return "azure"
	case commonpamconfiguration.EnvDomain:
		return "domain"
	case commonpamconfiguration.EnvGCP:
		return "gcp"
	default:
		return ""
	}
}
