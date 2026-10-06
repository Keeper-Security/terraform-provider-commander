// Copyright Keeper Security, Inc. 2026
// SPDX-License-Identifier: MPL-2.0

package saasconfiguration

import (
	commonrecordsutils "github.com/Keeper-Security/terraform-provider-commander/internal/provider/common/records/utils"
)

const (
	SaaSTypeCustomFieldLabel = "SaaS Type"
	ActiveCustomFieldLabel   = "Active"
)

// SaaSTypeCustomFieldType is the Keeper field type required for the SaaS Type custom field.
const SaaSTypeCustomFieldType = commonrecordsutils.FieldTypeText

// ActiveCustomFieldType is the Keeper field type required for the Active custom field.
// Active is stored as text ("true"/"false"), not a boolean-type Keeper field.
const ActiveCustomFieldType = commonrecordsutils.FieldTypeText

// Configuration and Gateway are used at create time to validate the selected SaaS Type
// (plugin) and its fields against the given gateway and configuration, and to link the
// record to that gateway/configuration. They are not populated on import - they're plan
// inputs for that create-time validation/linking, not values Commander returns on read.
const (
	ConfigurationDescription         = "PAM Configuration UID. Used at create time to validate the selected SaaS Type (plugin) and its fields against the given gateway and configuration. Not populated on import."
	ConfigurationMarkdownDescription = "The PAM configuration `UID`. Used at create time to validate the selected **SaaS Type** (plugin) and its fields against the given gateway and configuration. Not populated on import."

	GatewayDescription         = "The configured gateway `UID` or `name`. Used at create time to validate the selected SaaS Type (plugin) and its fields against the given gateway and configuration. Not populated on import."
	GatewayMarkdownDescription = "The configured gateway `UID` or `name`. Used at create time to validate the selected **SaaS Type** (plugin) and its fields against the given gateway and configuration. Not populated on import."
)

const (
	// ErrOpListSaasPlugins is the error-summary context passed to ExecuteCommand when listing
	// the SaaS plugins available for a gateway.
	ErrOpListSaasPlugins = "Unable to fetch available plugins for gateway"

	// ErrSummaryInvalidSaasType is the diagnostic summary used when the "SaaS Type" custom
	// field's value doesn't match any plugin returned for the configured gateway.
	ErrSummaryInvalidSaasType = "Invalid SaaS Type"

	// ErrOpGetSaasPluginInfo is the error-summary context passed to ExecuteCommand when fetching
	// a plugin's field spec.
	ErrOpGetSaasPluginInfo = "Unable to fetch plugin field's"

	// ErrSummaryInvalidSaasFields is the diagnostic summary used when the custom fields a
	// practitioner provided don't satisfy the selected SaaS plugin's required/optional fields.
	ErrSummaryInvalidSaasFields = "Invalid SaaS Configuration Fields"

	// ErrOpLinkSaasConfigToGateway is the error-summary context passed to ExecuteCommand when
	// connecting a newly created SaaS configuration record to its gateway.
	ErrOpLinkSaasConfigToGateway = "Unable to link SaaS Configuration record to gateway"

	// ErrOpChangeGatewayOrConfiguration is the error-summary context used when a practitioner attempts to change the gateway or configuration attributes of an existing SaaS configuration record.
	ErrOpChangeGatewayOrConfiguration = "Changing the gateway or configuration is not supported. Please create a new saas configuration record with the desired values."
)
