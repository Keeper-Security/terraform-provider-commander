# commander_classic_pam_user
#
# Creates and manages a Keeper PAM User record (`pamUser`) in the vault. PAM User
# records hold privileged credentials (login / password / private key) that can
# be associated with PAM Machines or Databases for rotation, connections, and
# tunneling.
#
# Required fields ............ title, login
# Optional fields ............ password, folder_location, notes,
#                              distinguished_name, private_pem_key, public_key,
#                              private_key_passphrase, connect_database,
#                              managed, rotation_settings, share
# Read-only fields ........... id (record UID assigned by Keeper after create)
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
# Usage 1 - Create a PAM User (no rotation)
#
# A simple credential vault entry. Use this when you just need to store a
# privileged credential without configuring automated rotation.
###############################################################################

resource "commander_classic_pam_user" "mysql_app_account" {
  title                  = "MySQL - billing app service account"
  login                  = "svc_billing"
  password               = "_REPLACE_WITH_STRONG_PASSWORD_"
  folder_location        = "Shared Folders/PAM/Database Users"
  notes                  = "Service account used by the billing app to connect to MySQL prod."
  distinguished_name     = "CN=svc_billing,OU=Service Accounts,DC=corp,DC=local"
  private_pem_key        = "_REPLACE_WITH_PEM_KEY_OR_USE_file()_"
  public_key             = "_REPLACE_WITH_PUBLIC_KEY_OR_USE_file()_"
  private_key_passphrase = "_REPLACE_WITH_PRIVATE_KEY_PASSPHRASE_OR_USE_file()_"
  connect_database       = "billing_prod"
  managed                = true

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

###############################################################################
# Usage 1b - Folder: optional at create; changing `folder_location` moves the record
#
# If you omit `folder_location`, the pamUser is created in your vault root. To place it
# in a Shared Folder from the start, set `folder_location` to a path (e.g.
# `Shared Folders/PAM/Service Accounts`) or a folder UID.
#
# If you add or change `folder_location` later, the next `terraform apply` moves the
# existing record (Commander `mv`), same as other classic PAM record resources.
###############################################################################

# resource "commander_classic_pam_user" "folder_after_root" {
#   title = "PAM User - start in root, move later"
#   login = "svc_example"
#   # Omit `folder_location` on first apply, then add the line below and apply again:
#   # folder_location = "Shared Folders/PAM/Service Accounts"
# }


###############################################################################
# Usage 2.1 - Rotation profile: "general"  (rotates ON a PAM resource)
#
# For type "general" you MUST pass:
#   - configuration : PAM Configuration UID for the rotation gateway / policy
#   - resource      : UID of the PAM Machine or PAM Database where rotation runs
#
# Do NOT set saas_config for this profile.
###############################################################################

# resource "commander_classic_pam_user" "mysql_rotation_user" {
#   title              = "PAM - MySQL Rotation User"
#   login              = "sqluser"
#   password           = "_REPLACE_WITH_STRONG_PASSWORD_"
#   distinguished_name = "CN=sqluser,OU=DB,DC=corp,DC=local"
#   connect_database   = "billing_prod"
#
#   rotation_settings = {
#     rotation_profile = "general"
#     configuration    = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"
#     resource         = "_REPLACE_WITH_PAM_MACHINE_OR_DATABASE_UID_"
#     complexity       = "32,5,1,1,2" # length, upper, lower, digits, symbols
#     enabled          = true
#     on_demand        = true
#   }
# }

###############################################################################
# Usage 2.2 - Rotation profile: "iam_user"  (cloud IAM / Azure AD)
#
# For type "iam_user" you MUST pass:
#   - configuration : PAM Configuration UID for IAM / Azure AD rotation
#
# Do NOT set resource or saas_config for this profile.
###############################################################################

# resource "commander_classic_pam_user" "aws_iam_deploy_user" {
#   title              = "AWS IAM - prod deploy user"
#   login              = "deploy@111122223333"
#   password           = "_REPLACE_WITH_STRONG_PASSWORD_"
#   distinguished_name = "deploy@aws-prod"
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
# Usage 2.3 - Rotation profile: "scripts_only"  (custom rotation via script)
#
# For type "scripts_only" you MUST pass:
#   - configuration : PAM Configuration UID that hosts the rotation script
#
# `login` / `password` are optional if the script does not need a starting credential.
# Do NOT set resource or saas_config for this profile.
###############################################################################

# resource "commander_classic_pam_user" "custom_script_runner" {
#   title              = "PAM - custom rotation runner"
#   login              = "svc_custom"
#   distinguished_name = "svc_custom"
#   connect_database   = "internal_app_db"
#   # password         = "_REPLACE_WITH_STRONG_PASSWORD_" # optional for scripts_only
#
#   rotation_settings = {
#     rotation_profile = "scripts_only"
#     configuration    = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"
#     complexity       = "32,1,1,1,1"
#     enabled          = true
#   }
# }

###############################################################################
# Usage 2.4 - Rotation profile: "saas"  (SaaS account rotation)
#
# For type "saas" you MUST pass:
#   - configuration : PAM Configuration UID
#   - saas_config   : SaaS Configuration record UID linked to that PAM Configuration
#
# Do NOT set resource for this profile.
###############################################################################

# resource "commander_classic_pam_user" "saas_app_user" {
#   title    = "SaaS - app service account"
#   login    = "svc_saas_app"
#   password = "_REPLACE_WITH_STRONG_PASSWORD_"
#
#   rotation_settings = {
#     rotation_profile = "saas"
#     configuration    = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"
#     saas_config      = "_REPLACE_WITH_SAAS_CONFIGURATION_UID_"
#     complexity       = "32,5,1,1,2"
#     enabled          = true
#     schedule_cron    = "0 0 4 * * ?" # Daily at 4 AM (Keeper Quartz, 6 fields)
#   }
# }

###############################################################################
# Usage 2.5 - Rotation schedule (cron / json / inherit-from-config)
#
# Pick exactly ONE of: on_demand, schedule_cron, schedule_json, use_default_rotation_schedule (required).
# `use_default_rotation_schedule = true` inherits the schedule from the PAM Configuration.
###############################################################################

# resource "commander_classic_pam_user" "scheduled_postgres_user" {
#   title              = "Postgres - scheduled rotation user"
#   login              = "svc_scheduled"
#   distinguished_name = "svc_scheduled"
#   connect_database   = "analytics_prod"
#
#   rotation_settings = {
#     rotation_profile = "scripts_only"
#     configuration    = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"
#     complexity       = "32,1,1,1,1"
#     enabled          = true
#
#     # ----- pick ONE of the schedule options below -----
#     schedule_cron = "0 0 3 1 * ?" # First of every month at 3 AM UTC
#     # schedule_cron   = "0 0 2 1 1,4,7,10 ?"                                                                            # First of every quarter at 2 AM UTC
#     # schedule_json   = "{\"type\": \"DAILY\", \"utcTime\": \"17:56\", \"intervalCount\": 1}"                            # Daily at 5:56 PM UTC
#     # schedule_json   = "{\"type\": \"WEEKLY\", \"utcTime\": \"00:00\", \"weekday\": \"SATURDAY\", \"intervalCount\": 1}" # Weekly on Saturday at 12:00 AM UTC
#     # use_default_rotation_schedule = true                                                                                            # Inherit schedule from the PAM Configuration
#     # on_demand       = true                                                                                            # Manual rotation only
#   }
# }

###############################################################################
# Usage 3 - Import an existing PAM User into Terraform
#
# Define an empty (or near-empty) resource block and run `terraform import`.
# After import, run `terraform plan` and align this block with the remote state.
# SaaS profiles are detected from vault dagDebug parentAclEdge rotation settings.
###############################################################################

# resource "commander_classic_pam_user" "imported_pam_user" {
#   title = "Imported PAM User"
#   login = "imported_user"
#   rotation_settings = {}
# }
