# Keeper Security Commander Provider Examples

This directory contains examples used by [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs) to generate provider documentation. You can also run them manually with the Terraform CLI (set provider config and credentials first).

## Layout

The documentation generator expects files in these paths:

| Path | Purpose |
|------|---------|
| `provider/provider.tf` | Provider configuration example (provider index page) |
| `resources/<resource_type>/resource.tf` | Resource configuration example for the named resource |
| `resources/<resource_type>/import.sh` | Import examples (CLI and import block) for the named resource |
| `data-sources/<data_source_type>/data-source.tf` | Data source example for the named data source |

Other `.tf` files in these directories are ignored by the docs tool but can be used for runnable or testable examples.

## Resources

- **commander_managed_company** — Managed company (MSP)
- **commander_enterprise_node** — Enterprise node
- **commander_enterprise_team** — Enterprise team
- **commander_enterprise_role** — Enterprise role
- **commander_enterprise_user** — Enterprise user
- **commander_enterprise_push** — One-time push of records to user vaults (write-only; see `resources/commander_enterprise_push/README.md` for details)
- **commander_enterprise_scim_push** — One-time push of SCIM data (Google, AD, or record) to a SCIM endpoint
- **commander_epm_policy** — EPM (Endpoint Policy Management) policy
- **commander_classic_shared_folder** — Classic shared folder (vault path, default and per-record/per-user permissions)
- **commander_classic_pam_database** — Classic PAM database vault record
- **commander_classic_pam_directory** — Classic PAM directory vault record
- **commander_classic_pam_machine** — Classic PAM machine vault record
- **commander_classic_pam_remote_browser** — Classic PAM remote browser (RBI) vault record
- **commander_classic_pam_user** — Classic PAM user vault record (login/password/PEM credentials, optional rotation settings)
- **commander_new_pam_database** — Nested-shared (new) PAM database vault record
- **commander_new_pam_directory** — Nested-shared (new) PAM directory vault record
- **commander_new_pam_machine** — Nested-shared (new) PAM machine vault record
- **commander_new_pam_remote_browser** — Nested-shared (new) PAM remote browser (RBI) vault record
- **commander_new_pam_user** — Nested-shared (new) PAM user vault record (login/password/PEM credentials, optional rotation settings)

## Data sources

- **commander_managed_company** — Look up a managed company
- **commander_enterprise_node** — Look up an enterprise node
- **commander_enterprise_team** — Look up an enterprise team
- **commander_enterprise_role** — Look up an enterprise role
- **commander_enterprise_user** — Look up an enterprise user
- **commander_epm_policy** — Look up an EPM policy by policy ID
- **commander_classic_shared_folder** — Look up a classic shared folder by UID or vault path
- **commander_pam_database** — Look up a PAM database vault record by record UID
- **commander_pam_directory** — Look up a PAM directory vault record by record UID
- **commander_pam_machine** — Look up a PAM machine vault record by record UID
- **commander_pam_remote_browser** — Look up a PAM remote browser vault record by record UID
- **commander_pam_user** — Look up a PAM user vault record by record UID

## Regenerating docs

From the repo root, run:

```bash
make generate
```

This runs tfplugindocs and updates `docs/` from the schema and these examples.
