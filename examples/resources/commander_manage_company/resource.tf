# Managed company (MSP accounts): name, node, plan, optional seats, file plan, and add-ons.

resource "commander_manage_company" "example" {
  name      = "Acme Corp"
  node      = "Root"
  plan      = "enterprise"
  seats     = 25
  file_plan = "1tb"

  add_ons = [
    "secrets_manager",
    "enterprise_breach_watch",
    "enterprise_audit_and_reporting",
    "password_rotation",
    "privileged_access_manager:1"
  ]
}
