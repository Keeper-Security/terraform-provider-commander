resource "commander_enterprise_scim" "example" {
  node          = "Engineering"
  unique_groups = true

  # Optional: SCIM groups starting with this prefix are imported as Roles
  # prefix = "SCIM-"

  # Optional, MSP only: scope to a specific managed company
  # managed_company = "Acme Corp"
}

output "scim_id" {
  value = commander_enterprise_scim.example.id
}

output "scim_url" {
  value = commander_enterprise_scim.example.scim_url
}

output "provisioning_token" {
  value     = commander_enterprise_scim.example.provisioning_token
  sensitive = true
}
