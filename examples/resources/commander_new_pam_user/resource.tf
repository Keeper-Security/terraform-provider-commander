# commander_new_pam_user
#
# Creates and manages a Keeper Nested-Shared (new) PAM User record (`pamUser`)
# in the vault. PAM User records hold privileged credentials (login, password,
# private key) that can be associated with PAM Machines or Databases for
# rotation, connections, and tunneling.
#
# Required fields ............ title
# Optional fields ............ login, password, folder, notes,
#                              distinguished_name, private_pem_key,
#                              connect_database, managed, rotation_settings,
#                              share
# Read-only fields ........... id (record UID assigned by Keeper after create)
#
# -----------------------------------------------------------------------------
# Which UID does each rotation_profile need?
# -----------------------------------------------------------------------------
#   rotation_profile = "general"        => set `resource`        (PAM Machine
#                                                                 or PAM Database UID)
#                                          and usually `admin_user`
#                                          (UID of the rotating account)
#   rotation_profile = "iam_user"       => set `configuration`   (PAM Configuration UID
#                                                                 for IAM / Azure AD)
#   rotation_profile = "scripts_only"   => set `configuration`   (PAM Configuration UID
#                                                                 that runs the script)
# -----------------------------------------------------------------------------

###############################################################################
# Usage 1 - Create a PAM User (no rotation) with per-user share permissions
###############################################################################

resource "commander_new_pam_user" "mysql_app_account" {
  title              = "MySQL - billing app service account"
  login              = "svc_billing"
  password           = "_REPLACE_WITH_STRONG_PASSWORD_"
  folder_location    = "_REPLACE_WITH_NSF_FOLDER_UID_OR_PATH_"
  notes              = "Service account used by the billing app to connect to MySQL prod."
  distinguished_name = "CN=svc_billing,OU=Service Accounts,DC=corp,DC=local"
  private_pem_key    = "_REPLACE_WITH_PEM_KEY_OR_USE_file()_"
  connect_database   = "billing_prod"
  managed            = true

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
#   - resource    : UID of the PAM Machine / PAM Database record where the
#                   credential lives and where Keeper will perform the rotation.
#   - admin_user  : UID of the PAM User Keeper uses to log in and rotate the
#                   target account (usually a privileged admin account).
###############################################################################

# resource "commander_new_pam_user" "mysql_rotation_user" {
#   title              = "PAM - MySQL Rotation User"
#   login              = "sqluser"
#   password           = "_REPLACE_WITH_STRONG_PASSWORD_"
#   distinguished_name = "CN=sqluser,OU=DB,DC=corp,DC=local"
#   connect_database   = "billing_prod"
#
#   rotation_settings = {
#     rotation_profile = "general"
#     resource         = "_REPLACE_WITH_PAM_MACHINE_OR_DATABASE_UID_" # required for "general"
#     admin_user       = "_REPLACE_WITH_ADMIN_PAM_USER_UID_"          # privileged rotator UID
#     complexity       = "32,5,5,5,5"                                 # length, upper, lower, digits, symbols
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
###############################################################################

# resource "commander_new_pam_user" "aws_iam_deploy_user" {
#   title              = "AWS IAM - prod deploy user"
#   login              = "deploy@111122223333"
#   password           = "_REPLACE_WITH_STRONG_PASSWORD_"
#   distinguished_name = "deploy@aws-prod"
#
#   rotation_settings = {
#     rotation_profile = "iam_user"
#     configuration    = "_REPLACE_WITH_PAM_CONFIGURATION_UID_FOR_IAM_" # required for "iam_user"
#     complexity       = "32,5,5,5,5"
#     enabled          = true
#     on_demand        = true
#   }
# }

###############################################################################
# Usage 4 - Rotation schedule (cron / json / inherit-from-config)
#
# Pick exactly ONE of: schedule_cron, schedule_json, schedule_config.
###############################################################################

# resource "commander_new_pam_user" "scheduled_postgres_user" {
#   title              = "Postgres - scheduled rotation user"
#   distinguished_name = "svc_scheduled"
#   connect_database   = "analytics_prod"
#
#   rotation_settings = {
#     rotation_profile = "scripts_only"
#     configuration    = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"
#     complexity       = "20,1,1,1,1"
#     enabled          = true
#
#     # ----- pick ONE of the schedule options below -----
#     schedule_cron = "0 0 3 1 * ?" # First of every month at 3 AM UTC
#     # schedule_json   = "{\"type\": \"DAILY\", \"utcTime\": \"17:56\", \"intervalCount\": 1}"
#     # schedule_config = true                                                                                            # Inherit schedule from the PAM Configuration
#   }
# }

###############################################################################
# Usage 5 - Import an existing PAM User into Terraform
###############################################################################

# resource "commander_new_pam_user" "imported_pam_user" {
#   rotation_settings = {}
# }
