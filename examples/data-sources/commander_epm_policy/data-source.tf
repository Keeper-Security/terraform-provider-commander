data "commander_epm_policy" "example" {
  policy = "your-epm-policy-id"

  # Optional, MSP only — scope lookup to a managed company.
  # managed_company = "Acme Corp"
}

# Reference a policy created in the same configuration:
# data "commander_epm_policy" "from_resource" {
#   policy = commander_epm_policy.example.id
# }

output "epm_policy_id" {
  description = "EPM policy ID (matches the policy argument)."
  value       = data.commander_epm_policy.example.id
}

output "epm_policy_name" {
  description = "Display name of the policy."
  value       = data.commander_epm_policy.example.policy_name
}

output "epm_policy_type" {
  description = "Policy type: elevation, file_access, command, or least_privilege."
  value       = data.commander_epm_policy.example.policy_type
}

output "epm_policy_status" {
  description = "Policy status: enforce, monitor, monitor_and_notify, or off."
  value       = data.commander_epm_policy.example.status
}

output "epm_policy_applications" {
  description = "Application collection IDs scoped by this policy, or \"*\" for all."
  value       = data.commander_epm_policy.example.applications
}

output "epm_policy_machine_collections" {
  description = "Machine collection IDs, or \"*\" for all machines."
  value       = data.commander_epm_policy.example.machine_collections
}

output "epm_policy_user_groups" {
  description = "User collection IDs, or \"*\" for all users."
  value       = data.commander_epm_policy.example.user_groups
}

output "epm_policy_control" {
  description = "Control actions: audit, notify, mfa, justify, approval."
  value       = data.commander_epm_policy.example.control
}

output "epm_policy_date_filter" {
  description = "Date filter ranges in YYYY-MM-DD:YYYY-MM-DD format."
  value       = data.commander_epm_policy.example.date_filter
}

output "epm_policy_day_filter" {
  description = "Days of week included in the policy schedule."
  value       = data.commander_epm_policy.example.day_filter
}

output "epm_policy_time_filter" {
  description = "Time-of-day ranges as start-end hours (0–23), e.g. 9-12."
  value       = data.commander_epm_policy.example.time_filter
}
