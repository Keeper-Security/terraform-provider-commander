# Changelog

## 1.3.0 (YYYY-MM-DD)

### Added

**Resources**

- **Classic Records:** `commander_classic_login`, `commander_classic_wifi`, `commander_classic_contact`, `commander_classic_address`, `commander_classic_payment_card`, `commander_classic_bank_account`, `commander_classic_membership`, `commander_classic_health_insurance`, `commander_classic_driver_license`, `commander_classic_passport`, `commander_classic_ssn_card`, `commander_classic_birth_certificate`, `commander_classic_ssh_keys`, `commander_classic_saas_configuration`, `commander_classic_server`, `commander_classic_database`, `commander_classic_software_license`, `commander_classic_secure_note`
- **New PAM Records (NSF):** `commander_new_pam_user`, `commander_new_pam_machine`, `commander_new_pam_database`, `commander_new_pam_directory`, `commander_new_pam_remote_browser`

**Data sources**

- **Classic Records:** `commander_classic_login`, `commander_classic_wifi`, `commander_classic_contact`, `commander_classic_address`, `commander_classic_payment_card`, `commander_classic_bank_account`, `commander_classic_membership`, `commander_classic_health_insurance`, `commander_classic_driver_license`, `commander_classic_passport`, `commander_classic_ssn_card`, `commander_classic_birth_certificate`, `commander_classic_ssh_keys`, `commander_classic_saas_configuration`, `commander_classic_server`, `commander_classic_database`, `commander_classic_software_license`, `commander_classic_secure_note`
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
