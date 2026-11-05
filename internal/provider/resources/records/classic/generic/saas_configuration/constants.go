// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

const (
	SchemaDescription         = "Creates and manages a SaaS configuration record in the vault. \n\n" + "The plugin to use must be set as the value of a custom field labeled `SaaS Type`. Because the list of available plugins is fetched from the gateway itself, `gateway` and `configuration` must both be set for Commander to look up and validate plugins for that gateway; see the [SaaS Plugins](https://docs.keeper.io/keeperpam/privileged-access-manager/password-rotation/rotation-use-cases/saas-plugins#keeperpam-saas-rotation-plugins) documentation for the full list of supported plugins. \n\n" + "After importing an existing SaaS configuration record, `gateway` and `configuration` are not populated; set them manually so Commander can validate the `SaaS Type` against that gateway."
	SchemaMarkdownDescription = "Creates and manages a **SaaS Configuration** record in the vault. \n\n" + "The plugin to use must be set as the value of a custom field labeled `SaaS Type`. Because the list of available plugins is fetched from the gateway itself, `gateway` and `configuration` must both be set for Commander to look up and validate plugins for that gateway; see the [SaaS Plugins](https://docs.keeper.io/keeperpam/privileged-access-manager/password-rotation/rotation-use-cases/saas-plugins#keeperpam-saas-rotation-plugins) documentation for the full list of supported plugins. \n\n" + "After importing an existing SaaS configuration record, `gateway` and `configuration` are not populated; set them manually so Commander can validate the `SaaS Type` against that gateway."

	ErrSummaryCreateFailed = "SaaS Configuration Record Create Failed"
	ErrSummaryReadFailed   = "SaaS Configuration Record Read Failed"
	ErrSummaryUpdateFailed = "SaaS Configuration Record Update Failed"

	ErrDetailCreateFailed = "Something went wrong when creating the SaaS Configuration record."
	ErrDetailReadFailed   = "Something went wrong when reading the SaaS Configuration record."
	ErrDetailUpdateFailed = "Something went wrong when updating the SaaS Configuration record."
)
