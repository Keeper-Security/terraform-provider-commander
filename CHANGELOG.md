# Changelog

## 1.3.0 (YYYY-MM-DD)

### Added

**Resources**

- **New PAM Records (NSF):** `commander_new_pam_user`, `commander_new_pam_machine`, `commander_new_pam_database`, `commander_new_pam_directory`, `commander_new_pam_remote_browser`

**Data sources**

- **New PAM Records (NSF):** `commander_new_pam_user`, `commander_new_pam_machine`, `commander_new_pam_database`, `commander_new_pam_directory`, `commander_new_pam_remote_browser`

## 1.2.0 (YYYY-MM-DD)

### Added

**Resources**

- **Classic Folders:** `commander_non_shared_folder`
- **Nested Shared Folders (NSF):** `commander_new_folder`
- **KeeperPAM:** `commander_pam_configuration`
- **Classic PAM Records:** `commander_classic_pam_user`, `commander_classic_pam_machine`, `commander_classic_pam_database`, `commander_classic_pam_directory`, `commander_classic_pam_remote_browser`
- **Secrets Manager:** `commander_secrets_manager`

**Data sources**

- **Classic Folders:** `commander_non_shared_folder`
- **Nested Shared Folders (NSF):** `commander_new_folder`
- **KeeperPAM:** `commander_pam_configuration`
- **Classic PAM Records:** `commander_classic_pam_user`, `commander_classic_pam_machine`, `commander_classic_pam_database`, `commander_classic_pam_directory`, `commander_classic_pam_remote_browser`
- **Secrets Manager:** `commander_secrets_manager`

## 1.1.0 (YYYY-MM-DD)

Added new resources and data sources

- **Resources:** `commander_epm_policy`, `commander_shared_folder`, `commander_enterprise_push`, `commander_enterprise_scim`, `commander_enterprise_scim_push`
- **Data sources:** `commander_epm_policy`, `commander_shared_folder`, `commander_enterprise_scim`

## 1.0.0 (YYYY-MM-DD)

FEATURES:

- **Resources:** `commander_managed_company`, `commander_enterprise_node`, `commander_enterprise_team`, `commander_enterprise_role`, `commander_enterprise_user`
- **Data sources:** `commander_managed_company`, `commander_enterprise_node`, `commander_enterprise_role`, `commander_enterprise_team`, `commander_enterprise_user`

Initial release of the Terraform provider for Keeper Security. Manages enterprise nodes, roles, teams, users, and companies via the Keeper Commander Service Mode REST API.
