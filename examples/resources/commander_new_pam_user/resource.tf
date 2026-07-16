# commander_new_pam_user
#
# Creates and manages a Keeper Nested-Shared (new) PAM User record (`pamUser`)
# in the vault. PAM User records hold privileged credentials (login, password,
# private key) that can be associated with PAM Machines or Databases for
# rotation, connections, and tunneling.
#
# Required fields ............ title, login
# Optional fields ............ password, folder_location, notes,
#                              distinguished_name, private_pem_key, public_key,
#                              private_key_passphrase, connect_database,
#                              managed, rotation_settings, share
# Read-only fields ........... id (record UID assigned by Keeper after create)
#
# Note: `folder_location` cannot be changed after create for NSF PAM records.
#
# -----------------------------------------------------------------------------
# Which fields does each rotation_profile need?
# -----------------------------------------------------------------------------
#   rotation_profile = "general"      => `configuration` + `resource`
#   rotation_profile = "iam_user"       => `configuration`  (not `resource` or `saas_config`)
#   rotation_profile = "scripts_only"   => `configuration`
#   rotation_profile = "saas"           => `configuration` + `saas_config`
#
# Schedule: when enabled is not false, set exactly ONE of on_demand, use_default_rotation_schedule, schedule_cron, schedule_json (required).
# Cron uses Keeper Quartz format (6 or 7 fields, seconds first), e.g. "0 0 4 * * ?".
# Complexity: five integers length,upper,lower,digits,symbols (length min 20).
# -----------------------------------------------------------------------------

###############################################################################
# Usage 1 - Create a PAM User (no rotation) with per-user share permissions
###############################################################################

resource "commander_new_pam_user" "mysql_app_account" {
  title                  = "MySQL - billing app service account"
  login                  = "svc_billing"
  password               = "_REPLACE_WITH_STRONG_PASSWORD_"
  folder_location        = "_REPLACE_WITH_NSF_FOLDER_UID_OR_PATH_"
  notes                  = "Service account used by the billing app to connect to MySQL prod."
  distinguished_name     = "CN=svc_billing,OU=Service Accounts,DC=corp,DC=local"
  private_pem_key        = "_REPLACE_WITH_PEM_KEY_OR_USE_file()_"
  public_key             = "_REPLACE_WITH_PUBLIC_KEY_OR_USE_file()_"
  private_key_passphrase = "_REPLACE_WITH_PRIVATE_KEY_PASSPHRASE_OR_USE_file()_"
  connect_database       = "billing_prod"
  managed                = true

  # ----------------------------------------------------------------
  # Per-user share permissions (optional).
  # Map key = user email; value = one of:
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

###############################################################################
# Usage 2 - Rotation profile: "general"  (rotates ON a PAM resource)
#
# For type "general" you MUST pass:
#   - configuration : PAM Configuration UID
#   - resource      : UID of the PAM Machine or PAM Database
###############################################################################

# resource "commander_new_pam_user" "mysql_rotation_user" {
#   title              = "PAM - MySQL Rotation User"
#   login              = "sqluser"
#   password           = "_REPLACE_WITH_STRONG_PASSWORD_"
#   distinguished_name = "CN=sqluser,OU=DB,DC=corp,DC=local"
#   connect_database   = "billing_prod"
#   folder_location    = "_REPLACE_WITH_NSF_FOLDER_UID_OR_PATH_"
#
#   rotation_settings = {
#     rotation_profile = "general"
#     configuration    = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"
#     resource         = "_REPLACE_WITH_PAM_MACHINE_OR_DATABASE_UID_"
#     complexity       = "32,5,1,1,2"
#     enabled          = true
#     on_demand        = true
#   }
#
#   share = {
#     "alice@example.com" = "full-manager"
#   }
# }

###############################################################################
# Usage 3 - Rotation profile: "iam_user"  (cloud IAM / Azure AD)
#
# For type "iam_user" you MUST pass:
#   - configuration : PAM Configuration UID for IAM / Azure AD rotation
#
# Do NOT set resource or saas_config for this profile.
###############################################################################

# resource "commander_new_pam_user" "aws_iam_deploy_user" {
#   title              = "AWS IAM - prod deploy user"
#   login              = "deploy@111122223333"
#   password           = "_REPLACE_WITH_STRONG_PASSWORD_"
#   distinguished_name = "deploy@aws-prod"
#   folder_location    = "_REPLACE_WITH_NSF_FOLDER_UID_OR_PATH_"
#
#   rotation_settings = {
#     rotation_profile = "iam_user"
#     configuration    = "_REPLACE_WITH_PAM_CONFIGURATION_UID_FOR_IAM_"
#     complexity       = "32,5,1,1,2"
#     enabled          = true
#     on_demand        = true
#   }
# }

###############################################################################
# Usage 4 - Rotation profile: "saas"  (SaaS account rotation)
#
# For type "saas" you MUST pass:
#   - configuration : PAM Configuration UID
#   - saas_config   : SaaS Configuration record UID
###############################################################################

# resource "commander_new_pam_user" "saas_app_user" {
#   title           = "SaaS - app service account"
#   login           = "svc_saas_app"
#   password        = "_REPLACE_WITH_STRONG_PASSWORD_"
#   folder_location = "_REPLACE_WITH_NSF_FOLDER_UID_OR_PATH_"
#
#   rotation_settings = {
#     rotation_profile = "saas"
#     configuration    = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"
#     saas_config      = "_REPLACE_WITH_SAAS_CONFIGURATION_UID_"
#     complexity       = "32,5,1,1,2"
#     enabled          = true
#     schedule_cron    = "0 0 4 * * ?"
#   }
# }

###############################################################################
# Usage 5 - Rotation schedule (cron / json / inherit-from-config)
#
# Pick exactly ONE of: on_demand, schedule_cron, schedule_json, use_default_rotation_schedule (required).
###############################################################################

# resource "commander_new_pam_user" "scheduled_postgres_user" {
#   title              = "Postgres - scheduled rotation user"
#   login              = "svc_scheduled"
#   distinguished_name = "svc_scheduled"
#   connect_database   = "analytics_prod"
#   folder_location    = "_REPLACE_WITH_NSF_FOLDER_UID_OR_PATH_"
#
#   rotation_settings = {
#     rotation_profile = "scripts_only"
#     configuration    = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"
#     complexity       = "32,1,1,1,1"
#     enabled          = true
#
#     # ----- pick ONE of the schedule options below -----
#     schedule_cron = "0 0 3 1 * ?" # First of every month at 3 AM UTC
#     # schedule_json   = "{\"type\": \"DAILY\", \"utcTime\": \"17:56\", \"intervalCount\": 1}"
#     # use_default_rotation_schedule = true      # Inherit schedule from the PAM Configuration
#     # on_demand       = true      # Manual rotation only
#   }
# }

###############################################################################
# Usage 6 - Import an existing PAM User into Terraform
#
# After import, run `terraform plan` and align this block with the remote state.
# SaaS profiles are detected from vault dagDebug parentAclEdge rotation settings.
###############################################################################

# resource "commander_new_pam_user" "imported_pam_user" {
#   title = "Imported PAM User"
#   login = "imported_user"
#   rotation_settings = {}
# }
