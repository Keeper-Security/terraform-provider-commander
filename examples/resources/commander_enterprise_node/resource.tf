# Enterprise node: name, parent (required). Optional: toggle_isolated, managed_company (MSP only).

resource "commander_enterprise_node" "example" {
  name   = "Engineering"
  parent = "Root"

  # Optional: isolate this node from other nodes
  toggle_isolated = false

  # Optional, MSP only: scope to a specific managed company
  # managed_company = "Acme Corp"
}
