# Managed company (MSP accounts): name, node, plan, optional seats, file plan, and add-ons.

resource "commander_managed_company" "example" {
  name      = "Acme Corp"
  node      = "Root"
  plan      = "enterprise"
  seats     = 25
  file_plan = "1tb"

  add_ons = [
    "chat",
    "compliance_report",
    "connection_manager:2",
    "enterprise_audit_and_reporting",
    "enterprise_breach_watch",
    "keeper_endpoint_privilege_manager:3",
    "msp_service_and_support",
    "password_rotation",
    "privileged_access_manager:1",
    "remote_browser_isolation",
    "secrets_manager",
  ]
}
