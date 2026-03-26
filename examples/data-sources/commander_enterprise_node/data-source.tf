# Look up an enterprise node by name or ID. Returns id, name, parent, parent_id.
# Optional: managed_company (MSP only) to scope the lookup.

data "commander_enterprise_node" "example" {
  node = "Engineering"

  # Optional, MSP only
  # managed_company = "Acme Corp"
}

output "node_id" {
  value = data.commander_enterprise_node.example.id
}

output "node_name" {
  value = data.commander_enterprise_node.example.name
}

output "node_parent" {
  value = data.commander_enterprise_node.example.parent
}

output "node_parent_id" {
  value = data.commander_enterprise_node.example.parent_id
}
