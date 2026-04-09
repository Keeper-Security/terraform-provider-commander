// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package pamconfiguration

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

// environmentBlocksValidator ensures at most one environment block is set and, if any is set, it matches `environment`.
type environmentBlocksValidator struct{}

func (environmentBlocksValidator) Description(_ context.Context) string {
	return "At most one of local_network, aws, azure, domain, or gcp may be set; if set, it must match `environment`."
}

func (environmentBlocksValidator) MarkdownDescription(_ context.Context) string {
	return "At most one of `local_network`, `aws`, `azure`, `domain`, or `gcp` may be set; if set, it must match `environment`."
}

func (environmentBlocksValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config PamConfigurationResourceModel
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
	if config.LocalNetwork != nil {
		present = append(present, "local_network")
	}
	if config.Aws != nil {
		present = append(present, "aws")
	}
	if config.Azure != nil {
		present = append(present, "azure")
	}
	if config.Domain != nil {
		present = append(present, "domain")
	}
	if config.Gcp != nil {
		present = append(present, "gcp")
	}

	if len(present) == 0 {
		return
	}
	if len(present) > 1 {
		resp.Diagnostics.AddError(
			"Conflicting environment configuration blocks",
			fmt.Sprintf("Only one of local_network, aws, azure, domain, or gcp may be set; found: %s", strings.Join(present, ", ")),
		)
		return
	}
	if present[0] != expected {
		resp.Diagnostics.AddError(
			"Environment does not match configuration block",
			fmt.Sprintf("environment is %q but %q is set; use the %q block instead, and remove %q.", env, present[0], expected, present[0]),
		)
	}
}

func expectedBlockName(env string) string {
	switch env {
	case EnvLocal:
		return "local_network"
	case EnvAWS:
		return "aws"
	case EnvAzure:
		return "azure"
	case EnvDomain:
		return "domain"
	case EnvGCP:
		return "gcp"
	default:
		return ""
	}
}
