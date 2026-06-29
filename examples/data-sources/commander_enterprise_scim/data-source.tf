data "commander_enterprise_scim" "example" {
  scim = "1169425105420640"

  # Optional, MSP only
  # managed_company = "Acme Corp"
}

output "scim_id" {
  value = data.commander_enterprise_scim.example.scim_id
}

output "scim_url" {
  value = data.commander_enterprise_scim.example.scim_url
}

output "node_id" {
  value = data.commander_enterprise_scim.example.node_id
}

output "node_name" {
  value = data.commander_enterprise_scim.example.node_name
}

output "status" {
  value = data.commander_enterprise_scim.example.status
}

output "unique_groups" {
  value = data.commander_enterprise_scim.example.unique_groups
}
