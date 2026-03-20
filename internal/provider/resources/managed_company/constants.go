// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedcompany

// Plan constants.
const (
	PlanBusiness       = "business"
	PlanBusinessPlus   = "businessPlus"
	PlanEnterprise     = "enterprise"
	PlanEnterprisePlus = "enterprisePlus"
)

// File plan constants.
const (
	FilePlan100GB = "100gb"
	FilePlan1TB   = "1tb"
	FilePlan10TB  = "10tb"
)

/*
COMMANDER ADD-ONS TO KEEPER ADMIN CONSOLE NAMING MAPPING ---> FOR REFERENCE ONLY

keeper_endpoint_privilege_manager -> Endpoint Manager
enterprise_breach_watch -> Breach Watch
compliance_report -> Compliance Reporting
enterprise_audit_and_reporting -> Advanced Reporting & Alerts Module
msp_service_and_support -> Dedicated Service & Support
privileged_access_manager -> Privileged Access Manager
secrets_manager -> Keeper Secrets Manager (KSM)
connection_manager -> Keeper Connection Manager (On-Prem)
remote_browser_isolation -> Remote Browser Isolation
chat -> KeeperChat

*/

// Add-on name constants.
const (
	// Add-ons with number suffix.
	AddOnConnectionManager              = "connection_manager"
	AddOnPrivilegedAccessManager        = "privileged_access_manager"
	AddOnKeeperEndpointPrivilegeManager = "keeper_endpoint_privilege_manager"

	// Base add-ons (no number suffix).
	AddOnSecretsManager              = "secrets_manager"
	AddOnRemoteBrowserIsolation      = "remote_browser_isolation"
	AddOnChat                        = "chat"
	AddOnEnterpriseAuditAndReporting = "enterprise_audit_and_reporting"
	AddOnMspServiceAndSupport        = "msp_service_and_support"
	AddOnEnterpriseBreachWatch       = "enterprise_breach_watch"
	AddOnComplianceReport            = "compliance_report"
	AddOnPasswordRotation            = "password_rotation"
)
