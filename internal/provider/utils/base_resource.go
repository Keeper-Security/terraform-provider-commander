// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// BaseResource provides shared configuration and API manager validation for resources.
// Embed it in a resource struct and call Configure from your resource's Configure,
// and use EnsureApiManager() at the start of CRUD operations.
type BaseResource struct {
	ApiManager *api.ApiManager
}

// ConfigureResource sets the API manager from provider data.
// Call this from your resource's Configure: r.ConfigureResource(ctx, req, resp).
func (b *BaseResource) ConfigureResource(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	b.ApiManager = apiManager
}

// EnsureApiManager validates that ApiManager is configured and returns an error if not.
func (b *BaseResource) EnsureApiManager() error {
	if b.ApiManager == nil {
		return fmt.Errorf("the Keeper Commander provider is not properly configured. Please ensure the provider is set up with valid service_mode_url and service_mode_api_key")
	}
	return nil
}
