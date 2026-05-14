# Look up an existing PAM configuration by UID.
data "commander_pam_configuration" "existing" {
  pam_configuration = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"
}

# ──────────────────────────────────────────────────────────────────────────────
# General attributes
# ──────────────────────────────────────────────────────────────────────────────

output "pam_config_id" {
  description = "PAM configuration UID."
  value       = data.commander_pam_configuration.existing.id
}

output "pam_config_title" {
  value = data.commander_pam_configuration.existing.title
}

output "pam_config_environment" {
  description = "Environment type (local, aws, azure, gcp, domain)."
  value       = data.commander_pam_configuration.existing.environment
}

output "pam_config_gateway" {
  value = data.commander_pam_configuration.existing.gateway
}

output "pam_config_application_folder" {
  value = data.commander_pam_configuration.existing.application_folder
}

output "pam_config_schedule" {
  value = data.commander_pam_configuration.existing.schedule
}

output "pam_config_port_mapping" {
  value = data.commander_pam_configuration.existing.port_mapping
}

# ──────────────────────────────────────────────────────────────────────────────
# Permission / recording settings
# ──────────────────────────────────────────────────────────────────────────────

output "pam_config_connections" {
  value = data.commander_pam_configuration.existing.connections
}

output "pam_config_tunneling" {
  value = data.commander_pam_configuration.existing.tunneling
}

output "pam_config_rotation" {
  value = data.commander_pam_configuration.existing.rotation
}

output "pam_config_remote_browser_isolation" {
  value = data.commander_pam_configuration.existing.remote_browser_isolation
}

output "pam_config_connections_recording" {
  value = data.commander_pam_configuration.existing.connections_recording
}

output "pam_config_typescript_recording" {
  value = data.commander_pam_configuration.existing.typescript_recording
}

output "pam_config_ai_threat_detection" {
  value = data.commander_pam_configuration.existing.ai_threat_detection
}

output "pam_config_ai_terminate_session_on_detection" {
  value = data.commander_pam_configuration.existing.ai_terminate_session_on_detection
}

# ──────────────────────────────────────────────────────────────────────────────
# Environment-specific blocks (only the matching block will be populated)
# ──────────────────────────────────────────────────────────────────────────────

output "pam_config_local_network" {
  description = "Local network configuration (populated when environment is 'local')."
  value       = data.commander_pam_configuration.existing.local_network
}

output "pam_config_aws" {
  description = "AWS configuration (populated when environment is 'aws')."
  value       = data.commander_pam_configuration.existing.aws
}

output "pam_config_azure" {
  description = "Azure configuration (populated when environment is 'azure')."
  value       = data.commander_pam_configuration.existing.azure
}

output "pam_config_domain" {
  description = "Domain/AD configuration (populated when environment is 'domain')."
  value       = data.commander_pam_configuration.existing.domain
}

output "pam_config_gcp" {
  description = "GCP configuration (populated when environment is 'gcp')."
  value       = data.commander_pam_configuration.existing.gcp
}
