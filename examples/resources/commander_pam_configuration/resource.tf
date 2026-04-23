# PAM Configuration: environment, title, gateway, application_folder (required).
# Optional: schedule, port_mapping, feature toggles (connections, tunneling, rotation, etc.),
# and exactly one environment-specific block matching the chosen environment.
# id is read-only (assigned by Keeper after creation).

# ──────────────────────────────────────────────────────────────────────────────
# Example 1 — Local Network (environment = "local")
# ──────────────────────────────────────────────────────────────────────────────

resource "commander_pam_configuration" "local_example" {
  environment        = "local"
  title              = "On-Prem Data Center"
  gateway            = "my-gateway-uid"
  application_folder = "PAM Application Folder"

  schedule     = "0 2 * * *"
  port_mapping = ["22:2222", "3389:33389"]

  connections              = true
  tunneling                = true
  rotation                 = true
  remote_browser_isolation = false
  connections_recording    = true
  typescript_recording     = true
  ai_threat_detection      = true

  ai_terminate_session_on_detection = false

  local_network = {
    network_id   = "DC-East-1"
    network_cidr = "10.0.0.0/16"
  }
}

# ──────────────────────────────────────────────────────────────────────────────
# Example 2 — AWS (environment = "aws")
# ──────────────────────────────────────────────────────────────────────────────

resource "commander_pam_configuration" "aws_example" {
  environment        = "aws"
  title              = "AWS Production"
  gateway            = "aws-gateway-uid"
  application_folder = "PAM Application Folder"

  connections           = true
  tunneling             = true
  rotation              = true
  connections_recording = true
  typescript_recording  = true

  aws = {
    aws_id            = "aws-prod-account"
    access_key_id     = "AKIAIOSFODNN7EXAMPLE"
    access_secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    region_names      = ["us-east-1", "us-west-2", "eu-west-1"]
  }
}

# ──────────────────────────────────────────────────────────────────────────────
# Example 3 — Azure (environment = "azure")
# ──────────────────────────────────────────────────────────────────────────────

resource "commander_pam_configuration" "azure_example" {
  environment        = "azure"
  title              = "Azure Corp"
  gateway            = "azure-gateway-uid"
  application_folder = "PAM Application Folder"

  connections              = true
  tunneling                = true
  rotation                 = true
  remote_browser_isolation = true
  connections_recording    = true
  typescript_recording     = true
  ai_threat_detection      = true

  ai_terminate_session_on_detection = true

  azure = {
    azure_id        = "azure-corp-001"
    client_id       = "00000000-1111-2222-3333-444444444444"
    client_secret   = "your-azure-client-secret"
    subscription_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    tenant_id       = "ffffffff-0000-1111-2222-333333333333"
    resource_groups = ["rg-production", "rg-staging"]
  }
}

# ──────────────────────────────────────────────────────────────────────────────
# Example 4 — Domain / Active Directory (environment = "domain")
# ──────────────────────────────────────────────────────────────────────────────

resource "commander_pam_configuration" "domain_example" {
  environment        = "domain"
  title              = "Corporate AD"
  gateway            = "ad-gateway-uid"
  application_folder = "PAM Application Folder"

  schedule = "0 3 * * 6"

  connections           = true
  tunneling             = true
  rotation              = true
  connections_recording = true

  domain = {
    domain_id           = "CORP.EXAMPLE.COM"
    domain_hostname     = "dc01.corp.example.com"
    domain_port         = "636"
    domain_use_ssl      = true
    domain_scan_dc_cidr = true
    domain_network_cidr = "172.16.0.0/12"
    domain_admin        = "admin-record-uid"
  }
}

# ──────────────────────────────────────────────────────────────────────────────
# Example 5 — GCP (environment = "gcp")
# ──────────────────────────────────────────────────────────────────────────────

resource "commander_pam_configuration" "gcp_example" {
  environment        = "gcp"
  title              = "GCP US Central"
  gateway            = "gcp-gateway-uid"
  application_folder = "PAM Application Folder"

  connections           = true
  tunneling             = true
  rotation              = true
  connections_recording = true
  typescript_recording  = true

  gcp = {
    gcp_id              = "GCP-US-CENTRAL1"
    service_account_key = file("${path.module}/gcp-service-account.json")
    google_admin_email  = "admin@example.com"
    gcp_region          = "us-central1"
  }
}

# ──────────────────────────────────────────────────────────────────────────────
# Example 6 — Minimal local configuration (only required fields)
# ──────────────────────────────────────────────────────────────────────────────

# resource "commander_pam_configuration" "minimal" {
#   environment        = "local"
#   title              = "Minimal Config"
#   gateway            = "gateway-uid"
#   application_folder = "PAM Application Folder"
# }

# ──────────────────────────────────────────────────────────────────────────────
# Outputs
# ──────────────────────────────────────────────────────────────────────────────

output "local_pam_config_id" {
  description = "UID of the local network PAM configuration"
  value       = commander_pam_configuration.local_example.id
}

output "aws_pam_config_id" {
  description = "UID of the AWS PAM configuration"
  value       = commander_pam_configuration.aws_example.id
}

output "azure_pam_config_id" {
  description = "UID of the Azure PAM configuration"
  value       = commander_pam_configuration.azure_example.id
}

output "domain_pam_config_id" {
  description = "UID of the Domain PAM configuration"
  value       = commander_pam_configuration.domain_example.id
}

output "gcp_pam_config_id" {
  description = "UID of the GCP PAM configuration"
  value       = commander_pam_configuration.gcp_example.id
}
