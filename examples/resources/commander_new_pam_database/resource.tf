# ------------------------------------------------------------------
# Example 1 - Minimal: only required fields (title + hostname)
# ------------------------------------------------------------------
resource "commander_new_pam_database" "minimal" {
  title = "PAM Database - mongo-prod"

  hostname_or_ip = {
    hostname = "mongo.internal.example.com"
  }
}

# ------------------------------------------------------------------
# Example 2 - All database fields + share permissions
# ------------------------------------------------------------------
# Supported database_type values:
#   postgresql, postgresql-flexible, mysql, mysql-flexible,
#   mariadb, mariadb-flexible, mssql, oracle, mongodb
resource "commander_new_pam_database" "full_fields" {
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
  folder_location = "_REPLACE_WITH_NSF_FOLDER_UID_"

  # ----------------------------------------------------------------
  # Per-user share permissions (optional).
  # Map key = user email. Each value is one of:
  #   viewer, share-manager, content-manager,
  #   content-share-manager, full-manager.
  # Omit the block entirely to skip share reconciliation.
  # ----------------------------------------------------------------
  share = {
    "alice@example.com" = "full-manager"
    "bob@example.com"   = "content-manager"
    "carol@example.com" = "viewer"
  }
}

# ------------------------------------------------------------------
# Example 3 - Full PAM settings (tunnel + connection) + share
# ------------------------------------------------------------------
# Connection protocol on PAM Database records must be one of:
#   mysql, postgresql, sql-server
resource "commander_new_pam_database" "with_pam_settings" {
  title = "PAM Database - mysql-prod"

  hostname_or_ip = {
    hostname            = "mysql-prod.db.example.com"
    administrative_port = 3306
  }

  use_ssl         = true
  database_type   = "mysql"
  database_id     = "db-mysql-prod-01"
  provider_group  = "AWS"
  provider_region = "us-west-2"
  notes           = "Production MySQL cluster with Commander-managed sessions."
  folder_location = "_REPLACE_WITH_NSF_FOLDER_UID_"

  pam_settings {
    configuration              = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"
    administrative_credentials = "_REPLACE_WITH_ADMIN_CREDENTIAL_UID_"
    allow_supply_host          = false

    tunnel {
      enable                   = true
      remote_target_port       = 3306
      re_use_port              = true
      use_specified_local_port = true
      local_port               = 33306
    }

    connection {
      enable            = true
      protocol          = "mysql"
      connection_port   = 3306
      launch_credential = "_REPLACE_WITH_LAUNCH_CREDENTIAL_UID_"

      mysql {
        session_recording      = true
        typescript_recording   = true
        recording_include_keys = true
        allow_supply_user      = false
        read_only              = false
        database               = "app_production"
        disable_csv_export     = true
        disable_csv_import     = true
        color_scheme           = "black-white"
        font_size              = 12
      }
    }
  }

  share = {
    "alice@example.com" = "full-manager"
    "bob@example.com"   = "viewer"
  }
}
