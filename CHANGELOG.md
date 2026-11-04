# Changelog

## 1.5.0 (2026-11-04)

### Added

- **`commander_pam_configuration`:** `universal_secrets_sync` block for the AWS, Azure, and GCP environments (`enabled`, `folder_location`, `dry_run`, `sync_identity`, plus `vault_name` for Azure).
- **PAM connection settings:** `clipboard_data_size_limit` (256KB, 1MB, 10MB, or 50MB) for SSH, MySQL, and other PAM connection configurations.
- **PAM remote browser:** `allow_supply_user` attribute.
- **PAM database:** `keeper_db` support.
- **SaaS configuration (classic and new):** required `configuration` and `gateway` attributes, computed on the data sources. The SaaS type and fields are validated against the plugins available on the gateway, and a required `Active` custom field (`TRUE` or `FALSE`) is enforced.

### Changed

- SaaS configuration documentation and examples now explain that the SaaS Type must be set, and that the gateway and configuration must be set manually after importing records.
- Updated MinCommanderVersion to `18.1.9`

## 1.4.0 (2026-09-24)

### Added

**Resources**

- **New Records (NSF):** `commander_new_login`, `commander_new_wifi`, `commander_new_contact`, `commander_new_address`, `commander_new_payment_card`, `commander_new_bank_account`, `commander_new_membership`, `commander_new_health_insurance`, `commander_new_driver_license`, `commander_new_passport`, `commander_new_ssn_card`, `commander_new_birth_certificate`, `commander_new_ssh_keys`, `commander_new_saas_configuration`, `commander_new_server`, `commander_new_database`, `commander_new_software_license`, `commander_new_secure_note`

**Data sources**

- **New Records (NSF):** `commander_new_login`, `commander_new_wifi`, `commander_new_contact`, `commander_new_address`, `commander_new_payment_card`, `commander_new_bank_account`, `commander_new_membership`, `commander_new_health_insurance`, `commander_new_driver_license`, `commander_new_passport`, `commander_new_ssn_card`, `commander_new_birth_certificate`, `commander_new_ssh_keys`, `commander_new_saas_configuration`, `commander_new_server`, `commander_new_database`, `commander_new_software_license`, `commander_new_secure_note`

### Changed

- New (NSF) records and folders can be moved between folder locations.
- Minimum required Keeper Commander version is now `18.1.5` (previously `18.1.2`).
- The CI test matrix now runs the full Go test suite.

### Fixed

- Omitted `commander_new_folder.records` no longer causes perpetual record-membership drift.
- `commander_new_wifi` stores encryption and hidden-SSID values in their typed fields instead of duplicate custom fields; classic WiFi creation uses consistent labels.
- Deleted folders reported as "Cannot find" are removed from Terraform state; transient Service Mode 429/503 responses are retried.
- Enterprise node, role, and team name conflicts now provide Terraform import guidance.
- Empty phone and custom-field lists are rejected explicitly while omitted optional fields remain valid.

## 1.3.0 (2026-09-08)

### Added

**Resources**

- **Classic Records:** `commander_classic_login`, `commander_classic_wifi`, `commander_classic_contact`, `commander_classic_address`, `commander_classic_payment_card`, `commander_classic_bank_account`, `commander_classic_membership`, `commander_classic_health_insurance`, `commander_classic_driver_license`, `commander_classic_passport`, `commander_classic_ssn_card`, `commander_classic_birth_certificate`, `commander_classic_ssh_keys`, `commander_classic_saas_configuration`, `commander_classic_server`, `commander_classic_database`, `commander_classic_software_license`, `commander_classic_secure_note`
- **New PAM Records (NSF):** `commander_new_pam_user`, `commander_new_pam_machine`, `commander_new_pam_database`, `commander_new_pam_directory`, `commander_new_pam_remote_browser`

**Data sources**

- **Classic Records:** `commander_classic_login`, `commander_classic_wifi`, `commander_classic_contact`, `commander_classic_address`, `commander_classic_payment_card`, `commander_classic_bank_account`, `commander_classic_membership`, `commander_classic_health_insurance`, `commander_classic_driver_license`, `commander_classic_passport`, `commander_classic_ssn_card`, `commander_classic_birth_certificate`, `commander_classic_ssh_keys`, `commander_classic_saas_configuration`, `commander_classic_server`, `commander_classic_database`, `commander_classic_software_license`, `commander_classic_secure_note`
- **New PAM Records (NSF):** `commander_new_pam_user`, `commander_new_pam_machine`, `commander_new_pam_database`, `commander_new_pam_directory`, `commander_new_pam_remote_browser`

## 1.2.0 (2026-07-16)

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

## 1.1.0 (2026-05-07)

Added new resources and data sources

- **Resources:** `commander_epm_policy`, `commander_shared_folder`, `commander_enterprise_push`, `commander_enterprise_scim`, `commander_enterprise_scim_push`
- **Data sources:** `commander_epm_policy`, `commander_shared_folder`, `commander_enterprise_scim`

## 1.0.0 (2026-03-27)

FEATURES:

- **Resources:** `commander_managed_company`, `commander_enterprise_node`, `commander_enterprise_team`, `commander_enterprise_role`, `commander_enterprise_user`
- **Data sources:** `commander_managed_company`, `commander_enterprise_node`, `commander_enterprise_role`, `commander_enterprise_team`, `commander_enterprise_user`

Initial release of the Terraform provider for Keeper Security. Manages enterprise nodes, roles, teams, users, and companies via the Keeper Commander Service Mode REST API.
