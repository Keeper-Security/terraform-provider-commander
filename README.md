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

| Name | Description |
|------|-------------|
| `commander_enterprise_node` | Create and manage enterprise nodes (MSP or enterprise account). |
| `commander_enterprise_user` | Create and manage enterprise users. |
| `commander_enterprise_role` | Create and manage enterprise roles and policies. |
| `commander_enterprise_scim` | Create and manage enterprise SCIM configurations for automated provisioning. |
| `commander_enterprise_scim_push` | Push SCIM data to a Keeper SCIM endpoint in a single step. |
| `commander_enterprise_team` | Create and manage enterprise teams. |
| `commander_enterprise_scim` | Create and manage enterprise SCIM configurations for automated provisioning. |
| `commander_enterprise_push` | Push record data from a JSON file to users' or teams' Keeper vaults. |
| `commander_enterprise_scim_push` | Push SCIM data to a Keeper SCIM endpoint in a single step. |

#### MSP Management

| Name | Description |
|------|-------------|
| `commander_managed_company` | Create and manage managed companies (MSP only). |
| `commander_shared_folder` | Create and manage shared folders. |

#### Classic Folders

| Name | Description |
|------|-------------|
| `commander_non_shared_folder` | Create and manage non-shared vault folders. |
| `commander_shared_folder` | Create and manage classic shared folders. |

#### Nested Shared Folders (NSF)

| Name | Description |
|------|-------------|
| `commander_new_folder` | Create and manage nested shared folders. |

#### KeeperPAM

| Name | Description |
|------|-------------|
| `commander_pam_configuration` | Create and manage Keeper PAM configurations. |

#### Classic PAM Records

| Name | Description |
|------|-------------|
| `commander_classic_pam_user` | Create and manage classic PAM user records in the vault. |
| `commander_classic_pam_machine` | Create and manage classic PAM machine records in the vault. |
| `commander_classic_pam_database` | Create and manage classic PAM database records in the vault. |
| `commander_classic_pam_directory` | Create and manage classic PAM directory records in the vault. |
| `commander_classic_pam_remote_browser` | Create and manage classic PAM remote browser (RBI) records in the vault. |

#### New PAM Records (NSF)

| Name | Description |
|------|-------------|
| `commander_new_pam_user` | Create and manage NSF PAM user records in the vault. |
| `commander_new_pam_machine` | Create and manage NSF PAM machine records in the vault. |
| `commander_new_pam_database` | Create and manage NSF PAM database records in the vault. |
| `commander_new_pam_directory` | Create and manage NSF PAM directory records in the vault. |
| `commander_new_pam_remote_browser` | Create and manage NSF PAM remote browser (RBI) records in the vault. |

#### Endpoint Privilege Manager (EPM)

| Name | Description |
|------|-------------|
| `commander_epm_policy` | Create and manage EPM (Endpoint Policy Management) policies. |

#### Secrets Manager

| Name | Description |
|------|-------------|
| `commander_secrets_manager` | Create and manage Keeper Secrets Manager applications. |

### Data sources

#### Enterprise Management

| Name | Description |
|------|-------------|
| `commander_enterprise_node` | Look up an enterprise node by name or ID. |
| `commander_enterprise_user` | Look up an enterprise user by email or ID. |
| `commander_enterprise_role` | Look up an enterprise role by name or ID. |
| `commander_enterprise_scim` | Look up an enterprise SCIM configuration by ID, node, or managed company. |
| `commander_enterprise_team` | Look up an enterprise team by name or ID. |
| `commander_enterprise_scim` | Look up an enterprise SCIM configuration by ID, node, or managed company. |

#### MSP Management

| Name | Description |
|------|-------------|
| `commander_managed_company` | Look up a managed company by name or ID (MSP only). |
| `commander_shared_folder` | Look up an existing shared folder by UID or name. |

#### Classic Folders

| Name | Description |
|------|-------------|
| `commander_non_shared_folder` | Look up a non-shared folder by UID. |
| `commander_shared_folder` | Look up a classic shared folder by UID. |

#### Nested Shared Folders (NSF)

| Name | Description |
|------|-------------|
| `commander_new_folder` | Look up a nested shared folder by UID. |

#### KeeperPAM

| Name | Description |
|------|-------------|
| `commander_pam_configuration` | Look up a PAM configuration by UID. |

#### Classic PAM Records

| Name | Description |
|------|-------------|
| `commander_classic_pam_user` | Look up a classic PAM user record by record UID. |
| `commander_classic_pam_machine` | Look up a classic PAM machine record by record UID. |
| `commander_classic_pam_database` | Look up a classic PAM database record by record UID. |
| `commander_classic_pam_directory` | Look up a classic PAM directory record by record UID. |
| `commander_classic_pam_remote_browser` | Look up a classic PAM remote browser record by record UID. |

#### New PAM Records (NSF)

| Name | Description |
|------|-------------|
| `commander_new_pam_user` | Look up a NSF PAM user record by record UID. |
| `commander_new_pam_machine` | Look up a NSF PAM machine record by record UID. |
| `commander_new_pam_database` | Look up a NSF PAM database record by record UID. |
| `commander_new_pam_directory` | Look up a NSF PAM directory record by record UID. |
| `commander_new_pam_remote_browser` | Look up a NSF PAM remote browser record by record UID. |

#### Endpoint Privilege Manager (EPM)

| Name | Description |
|------|-------------|
| `commander_epm_policy` | Look up an existing EPM policy by its policy ID. |

#### Secrets Manager

| Name | Description |
|------|-------------|
| `commander_secrets_manager` | Look up a Secrets Manager application by name or UID. |

## Prerequisites

- **Keeper Commander Service Mode**: A service account running latest version of Commander Service Mode REST API. Make sure you are running **Commander version 18.0.11** or **later** before starting Service Mode
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0

## Setup and Installation

In order to communicate between the Terraform and Keeper, the customer is responsible for hosting a Keeper Commander Service Mode instance. This can be accomplished many ways depending on your IT requirements. Commander Service Mode can run as a foreground service on any machine, or it can be run in a Docker container locally or remotely on a server.

### Step 1. Commander Setup

Follow the setup steps documented in the [Commander Service Mode REST API](https://docs.keeper.io/en/keeperpam/commander-cli/service-mode-rest-api) section to install Keeper Commander and start the service.<br>
Commander Service Mode can run directly in the CLI, in the background on a local machine, on a remote server as a service, or under a Docker container. Using Docker is the recommended method.

Note the following Important Items:

1) The Request Queue System (API v2) must be enabled, e.g. `-q=y`

2) Make sure the following commands are in the list:

```
this-device,sync-down,switch-to-mc,switch-to-msp,msp-add,msp-down,msp-info,msp-remove,msp-update,enterprise-info,enterprise-node,enterprise-user,enterprise-role,enterprise-team,enterprise-down,enterprise-push,team-approve,record-add,record-update,rm,get,list,record-type-info,share-folder,rmdir,rndir,mkdir,epm,scim,mv,pam,secrets-manager,ln,share-record,nsf-mkdir,nsf-get,nsf-rmdir,nsf-record-add,nsf-record-update,nsf-rm,nsf-rmdir,nsf-rndir,nsf-share-folder,nsf-share-record,nsf-ln
```

> If you encounter a 429 Too Many Requests error due to rate limiting, you can configure rate-limit for your service mode using the `-rl` or `--ratelimit flag`.
>
>This allows you to configure the allowed number of requests per endpoint per IP address, for example:
> - `1000/minute`
> - `100000/hour`
> - `2000000/day`
> 
>Adjust these limits based on your expected traffic and system capacity.

After service creation, the API key will be displayed in the console output. Make sure to copy and store it securely. If you are using Docker, you can pull the API key from the logs with this command:

```bash
docker compose logs | grep -i "generated api key"
```

When the Commander service is up and running, you should be able to submit a curl request to the endpoint. For example:

```bash
curl -X POST 'https://localhost:8080/api/v2/executecommand-async' \
--header 'Content-Type: application/json' \
    --header 'api-key: <your-api-key>' \
    --data '{"command": "this-device"}'
```

If the tunnel is running and the API key is correct, you should get a response like this:

```json
{
    "success": true,
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "queued",
    "message": "Request queued successfully..."
}
```

Now that the service is up and running, you can use Service Mode URL and API Key in provider configuration.

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

Below example explain how you can manage your enterprise team with help of "commander_enterprise_team" resource.<br>
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

Below example explain how you can read your existing enterprise team with help of "commander_enterprise_team" data source.<br>
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

Please email commander@keepersecurity.com with any specific requirements that you have.
