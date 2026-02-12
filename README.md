# Terraform Provider Commander

This provider manages your keeper resources via the **Keeper Commander Service Mode** REST API. Use it to manage enterprise nodes, roles, teams, users, and managed companies in Terraform.

For more on Service Mode, see [Keeper Commander Service Mode](https://docs.keeper.io/en/keeperpam/commander-cli/service-mode-rest-api#keeper-commander-service-mode).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24 (only if building the provider from source)

## Usage

Add the provider to your Terraform configuration and configure it with your Commander Service Mode URL and API key:

```hcl
terraform {
  required_providers {
    commander = {
      source  = "registry.terraform.io/Keeper-Security/commander"
    }
  }
}

provider "commander" {
  service_mode_url    = "https://your-commander-service-mode.example.com"
  service_mode_api_key = "your-api-key"
}

# Example: manage an enterprise user
resource "commander_enterprise_user" "example" {
  email = "user@example.com"
  name  = "Example User"
  node  = "Root"
}
```

See the [provider documentation](https://registry.terraform.io/providers/Keeper-Security/commander/latest/docs) and the [examples](examples/) directory for more.

## Building the Provider

1. Clone the repository.
2. Enter the repository directory.
3. Build and install the provider:

```shell
go install
```

The provider binary will be installed to `$GOPATH/bin` (or `$HOME/go/bin` by default).

## Developing the Provider

You need [Go](https://golang.org/doc/install) installed (see [Requirements](#requirements)).

### Setting up for local development (unpublished provider)

Because the provider is not yet on the Terraform Registry, use **development overrides** so Terraform uses your local binary instead of downloading it.

1. **Build and install the provider** (from the repo root):

   ```shell
   go install
   ```

   The binary is placed in `$GOPATH/bin` (or `$HOME/go/bin` by default), typically as `terraform-provider-commander_dev`.

2. **Configure Terraform to use the local binary** by adding a `dev_overrides` block to your CLI config.

   Create or edit `~/.terraformrc` (Windows: `%APPDATA%\terraform.rc`):

   ```hcl
   provider_installation {
     dev_overrides {
       "registry.terraform.io/Keeper-Security/commander" = "/path/to/dir/containing/binary"
     }
     direct {}
   }
   ```

   Replace `/path/to/dir/containing/binary` with the directory that contains the provider binary, for example:

   - **macOS/Linux:** `$HOME/go/bin` (i.e. Output of `echo $HOME/go/bin`)
   - Or: `$GOPATH/bin` if you use a custom GOPATH

3. **Use the provider in Terraform** as you would once published: in your config use `required_providers` with `source = "Keeper-Security/commander"` and a version (e.g. `version = "0.1.0"`). Run `terraform init` and Terraform will use your local binary instead of the registry.

4. Now provider is setup for local development. Follow the [Usage](#usage) steps to use it.
