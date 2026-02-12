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

- **commander_manage_company** — Managed company (MSP)
- **commander_enterprise_node** — Enterprise node
- **commander_enterprise_team** — Enterprise team
- **commander_enterprise_role** — Enterprise role
- **commander_enterprise_user** — Enterprise user

## Data sources

- **commander_manage_company** — Look up a managed company
- **commander_enterprise_node** — Look up an enterprise node
- **commander_enterprise_team** — Look up an enterprise team
- **commander_enterprise_role** — Look up an enterprise role
- **commander_enterprise_user** — Look up an enterprise user

## Regenerating docs

From the repo root, run:

```bash
make generate
```

This runs tfplugindocs and updates `docs/` from the schema and these examples.
