# Enterprise node: name, parent (required). Optional: toggle_isolated.
# toggle_isolated defaults to false and is not supported on create; set it on update to turn isolation on or off.

resource "commander_enterprise_node" "example" {
  name            = "Engineering"
  parent          = "Root"
  toggle_isolated = true

  # Optional, MSP only: scope to a specific managed company
  # managed_company = "Acme Corp"
}
