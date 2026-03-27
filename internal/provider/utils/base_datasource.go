// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"fmt"

	"github.com/Keeper-Security/terraform-provider-commander/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// BaseDataSource provides shared configuration and API manager validation for data sources.
// Embed it in a datasource struct and call ConfigureDataSource from your datasource's Configure,
// and use EnsureApiManager() at the start of Read.
type BaseDataSource struct {
	ApiManager *api.ApiManager
}

// ConfigureDataSource sets the API manager from provider data.
// Call this from your datasource's Configure: d.ConfigureDataSource(ctx, req, resp).
func (b *BaseDataSource) ConfigureDataSource(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiManager, ok := req.ProviderData.(*api.ApiManager)
	if !ok {
		resp.Diagnostics.AddError(
			ERR_MSG_PROVIDER_CONFIGURATION_ERROR,
			fmt.Sprintf("The provider was not configured correctly. Expected API manager, but got: %T. Please check your provider configuration.", req.ProviderData),
		)
		return
	}

	b.ApiManager = apiManager
}

// EnsureApiManager validates that ApiManager is configured and returns an error if not.
func (b *BaseDataSource) EnsureApiManager() error {
	if b.ApiManager == nil {
		return fmt.Errorf("the Keeper Commander provider is not properly configured. Please ensure the provider is set up with valid service_mode_url and service_mode_api_key")
	}
	return nil
}
