# Enterprise node: name, parent (required). Optional: managed_company (MSP only).
# toggle_isolated defaults to false and is not supported on create; set it on update to turn isolation on or off.

resource "commander_enterprise_node" "example" {
  name   = "Engineering"
  parent = "Root"

  # Optional, MSP only: scope to a specific managed company
  # managed_company = "Acme Corp"
}
