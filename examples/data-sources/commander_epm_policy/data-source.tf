# Look up an EPM policy by Keeper-assigned policy ID.
# Optional: managed_company (MSP only) to scope the lookup.

data "commander_epm_policy" "example" {
  policy = "your-epm-policy-id"

  # Optional, MSP only
  # managed_company = "Acme Corp"
}

output "epm_policy_name" {
  value = data.commander_epm_policy.example.policy_name
}

output "epm_policy_type" {
  value = data.commander_epm_policy.example.policy_type
}
