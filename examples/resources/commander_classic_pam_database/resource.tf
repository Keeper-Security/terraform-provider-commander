# ------------------------------------------------------------------
# Example 1 – Minimal: only required fields (title + hostname)
# ------------------------------------------------------------------
resource "commander_classic_pam_database" "minimal" {
  title = "PAM Database - mongo-prod"

  hostname_or_ip = {
    hostname = "mongo.internal.example.com"
  }
}

# ------------------------------------------------------------------
# Example 2 – All database fields
# ------------------------------------------------------------------
# Supported database_type values:
#   postgresql, postgresql-flexible, mysql, mysql-flexible,
#   mariadb, mariadb-flexible, mssql, oracle, mongodb
resource "commander_classic_pam_database" "full_fields" {
  title = "PAM Database - pg-prod"

  hostname_or_ip = {
    hostname            = "pg-prod.db.example.com"
    administrative_port = 5432
  }

  use_ssl         = true
  database_type   = "postgresql"
  database_id     = "db-926724afcc"
  provider_group  = "AWS"
  provider_region = "us-east-1"
  notes           = "Primary PostgreSQL database for production workloads."
  folder_location = "_REPLACE_WITH_SHARED_FOLDER_UID_"

  # ----------------------------------------------------------------
  # Per-user share permissions (optional).
  # Map key = user email. Each value is { can_share, can_edit }.
  # Both flags default to false (view-only).
  # Omit the block entirely to skip share reconciliation.
  # ----------------------------------------------------------------
  share = {
    "alice@example.com" = {
      can_share = true
      can_edit  = true
    }
    "bob@example.com" = {
      can_edit = true
    }
  }
}

# ------------------------------------------------------------------
# Example 3 – Full PAM settings (tunnel + connection)
# ------------------------------------------------------------------
resource "commander_classic_pam_database" "with_pam_settings" {
  title = "PAM Database - mongo-k8s"

  hostname_or_ip = {
    hostname            = "mongo-k8s.internal.example.com"
    administrative_port = 27017
  }

  use_ssl         = true
  database_type   = "mongodb"
  database_id     = "mongo-cluster-01"
  provider_group  = "AWS"
  provider_region = "us-west-2"
  notes           = "MongoDB cluster accessed via Kubernetes gateway."
  folder_location = "_REPLACE_WITH_SHARED_FOLDER_UID_"

  pam_settings {
    configuration              = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"
    administrative_credentials = "_REPLACE_WITH_ADMIN_CREDENTIAL_UID_"
    allow_supply_host          = false

    tunnel {
      enable                   = true
      remote_target_port       = 27017
      re_use_port              = true
      use_specified_local_port = true
      local_port               = 37017
    }

    connection {
      enable            = true
      protocol          = "kubernetes"
      connection_port   = 8443
      launch_credential = "_REPLACE_WITH_LAUNCH_CREDENTIAL_UID_"

      kubernetes {
        session_recording      = true
        typescript_recording   = true
        recording_include_keys = true
        allow_supply_user      = false
        use_ssl                = true
        ignore_cert            = false
        read_only              = false
        namespace              = "mongo-ns"
        pod                    = "mongos-0"
        container              = "mongos"
        color_scheme           = "black-white"
        font_size              = 12
        scrollback             = 2000
      }
    }
  }
}