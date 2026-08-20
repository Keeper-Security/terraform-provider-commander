# Terraform Provider For Commander

## About

**Terraform Provider For Commander** lets you manage your Keeper Security enterprise or MSP configuration as infrastructure-as-code. The provider uses the **Keeper Commander Service Mode REST API** to manage your Keeper resources from Terraform, so you get declarative config, version control, and a clear audit trail while staying on Keeper’s zero-knowledge infrastructure. See [Available resources and data sources](#available-resources-and-data-sources) for the full list.

## Features

- **Resources:** Create and manage your Keeper resources from Terraform.
- **Import**: Most resources support import state so you can bring existing keeper resources under Terraform management with terraform import .
- **Data sources**: Read the existing resource data via data sources.
- **MSP support:** Use the optional `managed_company` attribute on enterprise resources and data sources to scope operations to a specific managed company.

## Available resources and data sources

> Full resource and data source documentation on the [Terraform Registry](https://registry.terraform.io/providers/keeper-security/commander/latest/docs).

### Resources

#### Enterprise Management

| Name                             | Description                                                                  |
| -------------------------------- | ---------------------------------------------------------------------------- |
| `commander_enterprise_node`      | Create and manage enterprise nodes (MSP or enterprise account).              |
| `commander_enterprise_user`      | Create and manage enterprise users.                                          |
| `commander_enterprise_role`      | Create and manage enterprise roles and policies.                             |
| `commander_enterprise_team`      | Create and manage enterprise teams.                                          |
| `commander_enterprise_scim`      | Create and manage enterprise SCIM configurations for automated provisioning. |
| `commander_enterprise_push`      | Push record data from a JSON file to users' or teams' Keeper vaults.         |
| `commander_enterprise_scim_push` | Push SCIM data to a Keeper SCIM endpoint in a single step.                   |

#### MSP Management

| Name                        | Description                                     |
| --------------------------- | ----------------------------------------------- |
| `commander_managed_company` | Create and manage managed companies (MSP only). |
| `commander_shared_folder`   | Create and manage shared folders.               |

#### Classic Folders

| Name                          | Description                                 |
| ----------------------------- | ------------------------------------------- |
| `commander_non_shared_folder` | Create and manage non-shared vault folders. |
| `commander_shared_folder`     | Create and manage classic shared folders.   |

#### Nested Shared Folders (NSF)

| Name                   | Description                              |
| ---------------------- | ---------------------------------------- |
| `commander_new_folder` | Create and manage nested shared folders. |

#### KeeperPAM

| Name                          | Description                                  |
| ----------------------------- | -------------------------------------------- |
| `commander_pam_configuration` | Create and manage Keeper PAM configurations. |

#### Endpoint Privilege Manager (EPM)

| Name                   | Description                                                  |
| ---------------------- | ------------------------------------------------------------ |
| `commander_epm_policy` | Create and manage EPM (Endpoint Policy Management) policies. |

#### Secrets Manager

| Name                        | Description                                            |
| --------------------------- | ------------------------------------------------------ |
| `commander_secrets_manager` | Create and manage Keeper Secrets Manager applications. |

#### Classic PAM Records

| Name                                   | Description                                                              |
| -------------------------------------- | ------------------------------------------------------------------------ |
| `commander_classic_pam_user`           | Create and manage classic PAM user records in the vault.                 |
| `commander_classic_pam_machine`        | Create and manage classic PAM machine records in the vault.              |
| `commander_classic_pam_database`       | Create and manage classic PAM database records in the vault.             |
| `commander_classic_pam_directory`      | Create and manage classic PAM directory records in the vault.            |
| `commander_classic_pam_remote_browser` | Create and manage classic PAM remote browser (RBI) records in the vault. |

#### New PAM Records (NSF)

| Name                               | Description                                                          |
| ---------------------------------- | -------------------------------------------------------------------- |
| `commander_new_pam_user`           | Create and manage NSF PAM user records in the vault.                 |
| `commander_new_pam_machine`        | Create and manage NSF PAM machine records in the vault.              |
| `commander_new_pam_database`       | Create and manage NSF PAM database records in the vault.             |
| `commander_new_pam_directory`      | Create and manage NSF PAM directory records in the vault.            |
| `commander_new_pam_remote_browser` | Create and manage NSF PAM remote browser (RBI) records in the vault. |

#### Classic Records

| Name                                   | Description                                                          |
| -------------------------------------- | -------------------------------------------------------------------- |
| `commander_classic_login`              | Create and manage classic login records in the vault.                |
| `commander_classic_wifi`               | Create and manage classic WiFi credentials records in the vault.     |
| `commander_classic_contact`            | Create and manage classic contact records in the vault.              |
| `commander_classic_address`            | Create and manage classic address records in the vault.              |
| `commander_classic_payment_card`       | Create and manage classic payment card records in the vault.         |
| `commander_classic_bank_account`       | Create and manage classic bank account records in the vault.         |
| `commander_classic_membership`         | Create and manage classic membership records in the vault.           |
| `commander_classic_health_insurance`   | Create and manage classic health insurance records in the vault.     |
| `commander_classic_driver_license`     | Create and manage classic driver's license records in the vault.     |
| `commander_classic_passport`           | Create and manage classic passport records in the vault.             |
| `commander_classic_ssn_card`           | Create and manage classic identity (SSN) card records in the vault.  |
| `commander_classic_birth_certificate`  | Create and manage classic birth certificate records in the vault.    |
| `commander_classic_ssh_keys`           | Create and manage classic SSH keys records in the vault.             |
| `commander_classic_saas_configuration` | Create and manage classic SaaS configuration records in the vault.   |
| `commander_classic_server`             | Create and manage classic server credentials records in the vault.   |
| `commander_classic_database`           | Create and manage classic database credentials records in the vault. |
| `commander_classic_software_license`   | Create and manage classic software license records in the vault.     |
| `commander_classic_secure_note`        | Create and manage classic secure note records in the vault.          |

#### Classic Folders

| Name                          | Description                                 |
| ----------------------------- | ------------------------------------------- |
| `commander_non_shared_folder` | Create and manage non-shared vault folders. |
| `commander_shared_folder`     | Create and manage classic shared folders.   |

#### Nested Shared Folders (NSF)

| Name                   | Description                              |
| ---------------------- | ---------------------------------------- |
| `commander_new_folder` | Create and manage nested shared folders. |

#### KeeperPAM

| Name                          | Description                                  |
| ----------------------------- | -------------------------------------------- |
| `commander_pam_configuration` | Create and manage Keeper PAM configurations. |

#### Classic PAM Records

| Name                                   | Description                                                              |
| -------------------------------------- | ------------------------------------------------------------------------ |
| `commander_classic_pam_user`           | Create and manage classic PAM user records in the vault.                 |
| `commander_classic_pam_machine`        | Create and manage classic PAM machine records in the vault.              |
| `commander_classic_pam_database`       | Create and manage classic PAM database records in the vault.             |
| `commander_classic_pam_directory`      | Create and manage classic PAM directory records in the vault.            |
| `commander_classic_pam_remote_browser` | Create and manage classic PAM remote browser (RBI) records in the vault. |

#### New PAM Records (NSF)

| Name                               | Description                                                          |
| ---------------------------------- | -------------------------------------------------------------------- |
| `commander_new_pam_user`           | Create and manage NSF PAM user records in the vault.                 |
| `commander_new_pam_machine`        | Create and manage NSF PAM machine records in the vault.              |
| `commander_new_pam_database`       | Create and manage NSF PAM database records in the vault.             |
| `commander_new_pam_directory`      | Create and manage NSF PAM directory records in the vault.            |
| `commander_new_pam_remote_browser` | Create and manage NSF PAM remote browser (RBI) records in the vault. |

#### Endpoint Privilege Manager (EPM)

| Name                   | Description                                                  |
| ---------------------- | ------------------------------------------------------------ |
| `commander_epm_policy` | Create and manage EPM (Endpoint Policy Management) policies. |

#### Secrets Manager

| Name                        | Description                                            |
| --------------------------- | ------------------------------------------------------ |
| `commander_secrets_manager` | Create and manage Keeper Secrets Manager applications. |

### Data sources

#### Enterprise Management

| Name                        | Description                                                               |
| --------------------------- | ------------------------------------------------------------------------- |
| `commander_enterprise_node` | Look up an enterprise node by name or ID.                                 |
| `commander_enterprise_user` | Look up an enterprise user by email or ID.                                |
| `commander_enterprise_role` | Look up an enterprise role by name or ID.                                 |
| `commander_enterprise_team` | Look up an enterprise team by name or ID.                                 |
| `commander_enterprise_scim` | Look up an enterprise SCIM configuration by ID, node, or managed company. |

#### MSP Management

| Name                        | Description                                         |
| --------------------------- | --------------------------------------------------- |
| `commander_managed_company` | Look up a managed company by name or ID (MSP only). |
| `commander_shared_folder`   | Look up an existing shared folder by UID or name.   |

#### Classic Folders

| Name                          | Description                             |
| ----------------------------- | --------------------------------------- |
| `commander_non_shared_folder` | Look up a non-shared folder by UID.     |
| `commander_shared_folder`     | Look up a classic shared folder by UID. |

#### Nested Shared Folders (NSF)

| Name                   | Description                            |
| ---------------------- | -------------------------------------- |
| `commander_new_folder` | Look up a nested shared folder by UID. |

#### KeeperPAM

| Name                          | Description                         |
| ----------------------------- | ----------------------------------- |
| `commander_pam_configuration` | Look up a PAM configuration by UID. |

#### Endpoint Privilege Manager (EPM)

| Name                   | Description                                      |
| ---------------------- | ------------------------------------------------ |
| `commander_epm_policy` | Look up an existing EPM policy by its policy ID. |

#### Secrets Manager

| Name                        | Description                                           |
| --------------------------- | ----------------------------------------------------- |
| `commander_secrets_manager` | Look up a Secrets Manager application by name or UID. |

#### Classic PAM Records

| Name                                   | Description                                                |
| -------------------------------------- | ---------------------------------------------------------- |
| `commander_classic_pam_user`           | Look up a classic PAM user record by record UID.           |
| `commander_classic_pam_machine`        | Look up a classic PAM machine record by record UID.        |
| `commander_classic_pam_database`       | Look up a classic PAM database record by record UID.       |
| `commander_classic_pam_directory`      | Look up a classic PAM directory record by record UID.      |
| `commander_classic_pam_remote_browser` | Look up a classic PAM remote browser record by record UID. |

#### New PAM Records (NSF)

| Name                               | Description                                            |
| ---------------------------------- | ------------------------------------------------------ |
| `commander_new_pam_user`           | Look up a NSF PAM user record by record UID.           |
| `commander_new_pam_machine`        | Look up a NSF PAM machine record by record UID.        |
| `commander_new_pam_database`       | Look up a NSF PAM database record by record UID.       |
| `commander_new_pam_directory`      | Look up a NSF PAM directory record by record UID.      |
| `commander_new_pam_remote_browser` | Look up a NSF PAM remote browser record by record UID. |

#### Classic Records

| Name                                   | Description                                                  |
| -------------------------------------- | ------------------------------------------------------------ |
| `commander_classic_login`              | Look up a classic login record by record UID.                |
| `commander_classic_wifi`               | Look up a classic WiFi credentials record by record UID.     |
| `commander_classic_contact`            | Look up a classic contact record by record UID.              |
| `commander_classic_address`            | Look up a classic address record by record UID.              |
| `commander_classic_payment_card`       | Look up a classic payment card record by record UID.         |
| `commander_classic_bank_account`       | Look up a classic bank account record by record UID.         |
| `commander_classic_membership`         | Look up a classic membership record by record UID.           |
| `commander_classic_health_insurance`   | Look up a classic health insurance record by record UID.     |
| `commander_classic_driver_license`     | Look up a classic driver's license record by record UID.     |
| `commander_classic_passport`           | Look up a classic passport record by record UID.             |
| `commander_classic_ssn_card`           | Look up a classic identity (SSN) card record by record UID.  |
| `commander_classic_birth_certificate`  | Look up a classic birth certificate record by record UID.    |
| `commander_classic_ssh_keys`           | Look up a classic SSH keys record by record UID.             |
| `commander_classic_saas_configuration` | Look up a classic SaaS configuration record by record UID.   |
| `commander_classic_server`             | Look up a classic server credentials record by record UID.   |
| `commander_classic_database`           | Look up a classic database credentials record by record UID. |
| `commander_classic_software_license`   | Look up a classic software license record by record UID.     |
| `commander_classic_secure_note`        | Look up a classic secure note record by record UID.          |

#### Classic Folders

| Name                          | Description                             |
| ----------------------------- | --------------------------------------- |
| `commander_non_shared_folder` | Look up a non-shared folder by UID.     |
| `commander_shared_folder`     | Look up a classic shared folder by UID. |

#### Nested Shared Folders (NSF)

| Name                   | Description                            |
| ---------------------- | -------------------------------------- |
| `commander_new_folder` | Look up a nested shared folder by UID. |

#### KeeperPAM

| Name                          | Description                         |
| ----------------------------- | ----------------------------------- |
| `commander_pam_configuration` | Look up a PAM configuration by UID. |

#### Classic PAM Records

| Name                                   | Description                                                |
| -------------------------------------- | ---------------------------------------------------------- |
| `commander_classic_pam_user`           | Look up a classic PAM user record by record UID.           |
| `commander_classic_pam_machine`        | Look up a classic PAM machine record by record UID.        |
| `commander_classic_pam_database`       | Look up a classic PAM database record by record UID.       |
| `commander_classic_pam_directory`      | Look up a classic PAM directory record by record UID.      |
| `commander_classic_pam_remote_browser` | Look up a classic PAM remote browser record by record UID. |

#### New PAM Records (NSF)

| Name                               | Description                                            |
| ---------------------------------- | ------------------------------------------------------ |
| `commander_new_pam_user`           | Look up a NSF PAM user record by record UID.           |
| `commander_new_pam_machine`        | Look up a NSF PAM machine record by record UID.        |
| `commander_new_pam_database`       | Look up a NSF PAM database record by record UID.       |
| `commander_new_pam_directory`      | Look up a NSF PAM directory record by record UID.      |
| `commander_new_pam_remote_browser` | Look up a NSF PAM remote browser record by record UID. |

#### Endpoint Privilege Manager (EPM)

| Name                   | Description                                      |
| ---------------------- | ------------------------------------------------ |
| `commander_epm_policy` | Look up an existing EPM policy by its policy ID. |

#### Secrets Manager

| Name                        | Description                                           |
| --------------------------- | ----------------------------------------------------- |
| `commander_secrets_manager` | Look up a Secrets Manager application by name or UID. |

## Prerequisites

- **Keeper Commander Service Mode**: A service account running Commander Service Mode REST API using `terraform-app-setup` command.
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0

## Setup and Installation

In order to communicate between the Terraform and Keeper, the customer is responsible for hosting a Keeper Commander Service Mode instance using `terraform-app-setup`. This can be accomplished many ways depending on your IT requirements.

### Step 0. Commander Installation

[Install Keeper Commander](https://docs.keeper.io/keeperpam/commander-cli/commander-installation-setup) on a workstation.

### Step 1. Commander Setup

Log in to Commander with that account:

```
keeper shell
login serviceuser@company.com
```

Ensure Docker is available on the host where Service Mode will run.

#### Run Terraform setup

```
My Vault> terraform-app-setup
```

The command writes a `docker-compose.yml` with a **Commander-only** service configured for the Terraform provider.

##### Service Mode / Docker

Creates the shared folder, Docker config record, KSM application, and client config, then prompts for:

| Prompt                        | Description                                                                                          |
| ----------------------------- | ---------------------------------------------------------------------------------------------------- |
| **Port**                      | Local port for Commander Service Mode. Default: `8900`.                                              |
| **Enable ngrok?**             | Optional public URL via ngrok. Default: No.                                                          |
| **Ngrok Auth Token**          | Required if ngrok is enabled.                                                                        |
| **Ngrok Custom Domain**       | Optional (for example `myapp.ngrok.io`). Press Enter to skip.                                        |
| **Enable Cloudflare?**        | Asked only if ngrok is disabled. Default: No.                                                        |
| **Cloudflare Tunnel Token**   | Required if Cloudflare is enabled.                                                                   |
| **Cloudflare Custom Domain**  | Required if Cloudflare is enabled (for example `commander.company.com`).                             |
| **Enable advanced security?** | Optional IP allow/deny lists, rate limiting, response encryption, and token expiration. Default: No. |

> **Ngrok and Cloudflare are mutually exclusive.** The Service Mode URL must be reachable from where Terraform runs. If Commander is on a private network, enable **ngrok** or **Cloudflare Tunnel** and use that public HTTPS URL as `service_mode_url` (or `COMMANDER_SERVICE_MODE_URL`) in the provider.

Queue mode (API v2) is enabled automatically. The command allowlist is fixed for Terraform-safe operations (enterprise/MSP, records, sharing, PAM, Secrets Manager, and related NSF commands). Commands are not prompted interactively.

Resources created (defaults):

| Resource                   | Default name                                       |
| -------------------------- | -------------------------------------------------- |
| Shared folder              | `Commander Service Mode - Terraform`               |
| KSM application            | `Commander Service Mode - Terraform KSM App`       |
| Docker config record       | `Commander Service Mode Terraform Config`          |
| Docker service / container | `commander-terraform` / `keeper-service-terraform` |

> Re-running setup rewrites `docker-compose.yml` (manual edits are lost).

#### Deploy

```
My Vault> quit
rm ~/.keeper/config.json
docker compose up -d
docker ps
docker logs keeper-service-terraform
curl http://localhost:<port>/health
```

Delete the local `config.json` before starting Docker so the container does not conflict with the same device token. Docker loads its own config through KSM.

Now that the service is up and running, you can use Service Mode URL (async - */api/v2/*) and API Key in provider configuration which is present in your `keeper-service-terraform` docker container logs.

> If you encounter a 429 Too Many Requests error due to rate limiting, you can configure rate-limit for your service mode in **terraform-app-setup command > Enable advanced security?** .
>
> This allows you to configure the allowed number of requests per endpoint per IP address, for example:
>
> - `1000/minute`
> - `100000/hour`
> - `2000000/day`
>
> Adjust these limits based on your expected traffic and system capacity.

> Keep the Commander Service Mode running in order to stay connected

### Step 2. Provider Installation

#### Registry install

To install this provider, add the following code to your Terraform configuration and run `terraform init`

```hcl
terraform {
  required_providers {
    commander = {
      source = "keeper-security/commander"
    }
  }
}

provider "commander" {
  # Configuration options
}
```

## Usage

### Configure the Provider

The provider needs to be configured with commander service mode url and api key before it can be used.

```hcl
terraform {
  required_providers {
    commander = {
      source = "keeper-security/commander"
    }
 }
}

provider "commander" {
  service_mode_url     = "http://localhost:8080/api/v2/"
  service_mode_api_key = "XXXXXXXXXXXXXX"
  timeout              = 60  # optional; defaults to 60 seconds (if not provided or is set to 0 or less) for HTTP and async command polling
}
```

You can omit `service_mode_url` and `service_mode_api_key` in the configuration and set them via environment variables instead: `COMMANDER_SERVICE_MODE_URL` and `COMMANDER_SERVICE_MODE_API_KEY`. Config values take precedence over environment variables.

> **Note: Using managed companies (MSP accounts)**  
> Many resources and data sources support an optional `managed_company` attribute. When your account is an MSP, set `managed_company` to a managed company name or ID to manage that resource inside that company. Omit it to work in the logged-in account context (MSP or single enterprise).

> **Note: MSP — Using both a managed company and your main account in the same config**  
> If you use some resources or data sources with `managed_company` (operations run inside that company) and others without it (operations run in the logged-in account context), Terraform may run them in parallel. Commander processes requests one at a time in a queue, so an action can run in the wrong context and fail (e.g. "resource not found").
>
> **Fix:** Add dependencies between those resources or data sources (e.g. `depends_on` or referencing one from the other) so they are not executed in parallel.
>
> **Example:** Force ordering so the main-account resource runs after the managed-company one:
>
> ```hcl
> # Runs in managed company "Acme"
> resource "commander_enterprise_team" "mc_team" {
>   name = "MC Team"
>   node = "Root"
>   managed_company = "Acme"
> }
>
> # Runs in logged-in account; depends on mc_team so it doesn't run in parallel
> resource "commander_enterprise_team" "main_team" {
>   name    = "Main Team"
>   node    = "Root"
>   # no managed_company = main account
>   depends_on = [commander_enterprise_team.mc_team]
> }
> ```

### Examples

#### Manage Enterprise Team

Below example explain how you can manage your enterprise team with help of "commander_enterprise_team" resource.

Use this resource to create and manage teams in the MSP or Enterprise account

```hcl
terraform {
  required_providers {
    commander = {
      source = "keeper-security/commander"
    }
  }
}

provider "commander" {
  service_mode_url     = "http://localhost:8080/api/v2/"
  service_mode_api_key = "XXXXXXXXXXXXXX"
}

resource "commander_enterprise_team" "example" {
  name                     = "Backend Developers"
  node                     = "Engineering"
  users                    = ["alice@example.com", "bob@example.com"]
  roles                    = ["Developer"]
  restrict_record_edit     = true
  restrict_record_re_share = true
  enable_privacy_screen    = false
  # Optional, MSP Account only
  # managed_company = "Acme Corp"
}
```

#### Read Enterprise Team

Below example explain how you can read your existing enterprise team with help of "commander_enterprise_team" data source.

Use this data source to look up an enterprise team by name or ID. Returns the team's ID, name, users, and roles so you can reference them in other resources.

```hcl
terraform {
  required_providers {
    commander = {
      source = "keeper-security/commander"
    }
 }
}

provider "commander" {
  service_mode_url     = "http://localhost:8080/api/v2/"
  service_mode_api_key = "XXXXXXXXXXXXXX"
}

data "commander_enterprise_team" "example" {
  team = "Backend Developers"
  # Optional, MSP only
  # managed_company = "Acme Corp"
}

output "team_id" {
  value = data.commander_enterprise_team.example.id
}

output "team_name" {
  value = data.commander_enterprise_team.example.name
}

output "team_users" {
  value = data.commander_enterprise_team.example.users
}

output "team_roles" {
  value = data.commander_enterprise_team.example.roles
}
```

For more examples on different resources and data sources, check out the detailed provider documentation [here](https://registry.terraform.io/providers/Keeper-Security/commander/latest/docs) .

Please email [commander@keepersecurity.com](mailto:commander@keepersecurity.com) with any specific requirements that you have.
