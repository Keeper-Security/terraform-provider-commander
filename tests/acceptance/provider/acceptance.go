// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	commander "github.com/Keeper-Security/terraform-provider-commander/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
)

// TestAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
// The factory function is called for each Terraform CLI command to create a provider
// server that the CLI can connect to and interact with.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"commander": providerserver.NewProtocol6WithError(commander.New("test")()),
}

// TestAccProtoV6ProviderFactoriesWithEcho includes the echo provider alongside the commander provider.
// It allows for testing assertions on data returned by an ephemeral resource during Open.
var TestAccProtoV6ProviderFactoriesWithEcho = map[string]func() (tfprotov6.ProviderServer, error){
	"commander": providerserver.NewProtocol6WithError(commander.New("test")()),
	"echo":      echoprovider.NewProviderServer(),
}

// AccProviderConfig returns a provider "commander" block HCL string for use in acceptance test Config.
// Pass the mock server base URL (e.g. server.URL from helpers.StartCommandServer) and any non-empty
// API key; the mock does not validate the key. No real Service Mode URL or API key is required.
func AccProviderConfig(serviceModeURL, apiKey string) string {
	return `provider "commander" {
  service_mode_url    = "` + serviceModeURL + `"
  service_mode_api_key = "` + apiKey + `"
}
`
}
