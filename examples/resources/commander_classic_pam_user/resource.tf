# commander_classic_pam_user
#
# Creates and manages a Keeper PAM User record (`pamUser`) in the vault. PAM User
# records hold privileged credentials (login / password / private key) that can
# be associated with PAM Machines or Databases for rotation, connections, and
# tunneling.
#
# Required fields ............ title
# Optional fields ............ login, password, folder, notes,
#                              distinguished_name, private_pem_key,
#                              connect_database, managed, rotation_settings
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
# Usage 1 - Create a PAM User (no rotation)
#
# A simple credential vault entry. Use this when you just need to store a
# privileged credential without configuring automated rotation.
###############################################################################

resource "commander_classic_pam_user" "mysql_app_account" {
  title              = "MySQL - billing app service account"
  login              = "svc_billing"
  password           = "_REPLACE_WITH_STRONG_PASSWORD_"
  folder             = "Shared Folders/PAM/Database Users"
  notes              = "Service account used by the billing app to connect to MySQL prod."
  distinguished_name = "CN=svc_billing,OU=Service Accounts,DC=corp,DC=local"
  private_pem_key    = "_REPLACE_WITH_PEM_KEY_OR_USE_file()_"
  connect_database   = "billing_prod"
  managed            = true
}

###############################################################################
# Usage 1b - Folder: optional at create; changing `folder` moves the record
#
# If you omit `folder`, the pamUser is created in your vault root. To place it
# in a Shared Folder from the start, set `folder` to a path (e.g.
# `Shared Folders/PAM/Service Accounts`) or a folder UID.
#
# If you add or change `folder` later, the next `terraform apply` moves the
# existing record (Commander `mv`), same as other PAM record resources.
###############################################################################

# resource "commander_classic_pam_user" "folder_after_root" {
#   title = "PAM User - start in root, move later"
#   # Omit `folder` on first apply, then add the line below and apply again:
#   # folder = "Shared Folders/PAM/Service Accounts"
# }


###############################################################################
# Usage 2.1 - Rotation profile: "general"  (rotates ON a PAM resource)
#
# For type "general" you MUST pass:
#   - resource    : UID of the PAM Machine / PAM Database record where the
#                   credential lives and where Keeper will perform the rotation.
#   - admin_user  : UID of the PAM User Keeper uses to log in and rotate the
#                   target account (usually a privileged admin account).
#
# Use this for Linux SSH users, Windows local accounts, MySQL / Postgres / SQL
# Server users that live on a specific machine or database.
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
#     resource         = "_REPLACE_WITH_PAM_MACHINE_OR_DATABASE_UID_" # <-- required for "general"
#     admin_user       = "_REPLACE_WITH_ADMIN_PAM_USER_UID_"          # <-- privileged rotator UID
#     complexity       = "32,5,5,5,5"                                 # length, upper, lower, digits, symbols
#     enabled          = true
#     on_demand        = true
#   }
# }

###############################################################################
# Usage 2.2 - Rotation profile: "iam_user"  (cloud IAM / Azure AD)
#
# For type "iam_user" you MUST pass:
#   - configuration : PAM Configuration UID that talks to your cloud IAM provider
#                     (AWS IAM, Azure AD, GCP IAM, etc.). Do NOT set `resource`
#                     for this profile.
#
# Use this for cloud IAM accounts that don't live on a single machine.
###############################################################################

# resource "commander_classic_pam_user" "aws_iam_deploy_user" {
#   title              = "AWS IAM - prod deploy user"
#   login              = "deploy@111122223333"
#   password           = "_REPLACE_WITH_STRONG_PASSWORD_"
#   distinguished_name = "deploy@aws-prod"
#
#   rotation_settings = {
#     rotation_profile = "iam_user"
#     configuration    = "_REPLACE_WITH_PAM_CONFIGURATION_UID_FOR_IAM_" # <-- required for "iam_user"
#     complexity       = "32,5,5,5,5"
#     enabled          = true
#     on_demand        = true
#   }
# }

###############################################################################
# Usage 2.3 - Rotation profile: "scripts_only"  (custom rotation via script)
#
# For type "scripts_only" you MUST pass:
#   - configuration : PAM Configuration UID that hosts the rotation script.
#                     `login` / `password` are optional - omit them if the
#                     script doesn't need a starting credential.
#
# Use this for custom rotation logic or non-standard targets.
###############################################################################

# resource "commander_classic_pam_user" "custom_script_runner" {
#   title              = "PAM - custom rotation runner"
#   distinguished_name = "svc_custom"
#   connect_database   = "internal_app_db"
#   # login            = "svc_custom"                    # optional for scripts_only
#   # password         = "_REPLACE_WITH_STRONG_PASSWORD_" # optional for scripts_only
#
#   rotation_settings = {
#     rotation_profile = "scripts_only"
#     configuration    = "_REPLACE_WITH_PAM_CONFIGURATION_UID_" # <-- required for "scripts_only"
#     complexity       = "20,1,1,1,1"
#     enabled          = true
#   }
# }

###############################################################################
# Usage 2.4 - Rotation schedule (cron / json / inherit-from-config)
#
# Pick exactly ONE of: schedule_cron, schedule_json, schedule_config.
# `schedule_config = true` inherits the schedule from the PAM Configuration
# (the easiest option when you manage many records uniformly).
###############################################################################

# resource "commander_classic_pam_user" "scheduled_postgres_user" {
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
#     # schedule_cron   = "0 0 2 1 1,4,7,10 ?"                                                                            # First of every quarter (Jan, Apr, Jul, Oct) at 2 AM UTC
#     # schedule_json   = "{\"type\": \"DAILY\", \"utcTime\": \"17:56\", \"intervalCount\": 1}"                            # Daily at 5:56 PM UTC
#     # schedule_json   = "{\"type\": \"WEEKLY\", \"utcTime\": \"00:00\", \"weekday\": \"SATURDAY\", \"intervalCount\": 1}" # Weekly on Saturday at 12:00 AM UTC
#     # schedule_config = true                                                                                            # Inherit schedule from the PAM Configuration
#   }
# }

###############################################################################
# Usage 3 - Import an existing PAM User into Terraform
#
# Define an empty (or near-empty) resource block and run the `terraform import`
# command in import.sh. After import, run `terraform plan` and align this block
# with the remote state.
###############################################################################

# resource "commander_classic_pam_user" "imported_pam_user" {
#   rotation_settings = {}
# }
