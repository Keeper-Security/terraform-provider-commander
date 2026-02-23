# commander_enterprise_push Example

This example shows how to use the **commander_enterprise_push** resource to push JSON file content to user vaults via the enterprise-push Commander CLI.

## What this resource does

- **One-time action**: When Terraform creates or replaces this resource, it runs the enterprise-push command with your JSON file and target emails/teams.
- **Deterministic ID**: The resource `id` is a hash of file content + sorted emails + sorted teams + managed_company. Same inputs → same ID → no re-push on later applies.
- **Write-only**: The API does not support read, update, or delete. Terraform cannot detect drift. Delete only removes the resource from state.

## Detailed instructions

### 1. Create a JSON file

Create a file that will be pushed (e.g. `push-data.json` in this directory). It must be **valid JSON**. The entire file content is sent as `filedata` to the enterprise-push command.

Example `push-data.json`:

```json
{
  "records": [{
    "title":"Google",
    "folders": [
      {
        "folder": "My Websites\\Online"
      }
    ],
    "login": "testing",
    "password": "lk4j139sk4j",
    "login_url": "https://google.com",
    "notes": "These are some notes.",
    "custom_fields": {"Favorite Food":"Cheetos"}
  },
  {
    "title":"Facebook",
    "folders": [
      {
        "folder": "Social Media"
      },
      {
        "shared_folder": "Shared Social",
        "can_edit": false,
        "can_share": false
      }
    ],
    "login": "me@gmail.com",
    "password": "123123123123",
    "login_url": "https://facebook.com",
    "notes": "This is our corporate shared record.",
    "custom_fields": {
      "Facebook Application ID":"ABC12345",
      "$oneTimeCode": "otpauth://totp/Amazon:me@company.com?secret=JBSWY3DPEHPK3PXP&issuer=Amazon&algorithm=SHA1&digits=6&period=30"}
  }]
}
```

Place the file where Terraform can read it. In the example we use `"${path.module}/push-data.json"` so the file lives next to the Terraform config.

### 2. Set targets (email and/or team)

- **email** (optional): Set of user email addresses to push to. If set, must have at least one value.
- **team** (optional): Set of team names or IDs to push to. If set, must have at least one value.
- You can set one or both. Omit both only if your workflow allows it.

### 3. MSP: managed_company (optional)

For MSP accounts, set `managed_company` to the name or ID of the managed company to run the push in. Omit to use the logged-in account context.

### 4. Apply

```bash
terraform init
terraform plan   # Shows: create commander_enterprise_push.example
terraform apply  # Runs enterprise-push; resource is created with computed id
```

### 5. When does it push again?

| Action | Result |
|--------|--------|
| Run `terraform apply` again with **no config or file change** | No re-push (same id). |
| **Edit the JSON file** (same path, different content) | Provider detects content change at plan time → replace → **push runs again to all email/team**. |
| **Change** `file_path` or `managed_company` | Replace → **push runs again** to all email/team. |
| **Add** email(s) and/or team(s) | **Update** → push runs **only to newly added** targets. |
| **Remove** email(s) and/or team(s) | **Update** → **no push** (state only). |
| **Remove** the resource from config and apply | Resource removed from state only; no API delete. |

**Adding or removing email/team:** Adding emails/teams triggers **Update** and the push goes **only to the newly added** emails/teams. Removing emails/teams triggers Update and **only updates state** (no push).

### 6. Delete

Removing the resource block and running `terraform apply` only removes it from Terraform state. Nothing is deleted on the server.

## File layout for this example

```
examples/resources/commander_enterprise_push/
  README.md        # This file
  resource.tf      # Example Terraform config (used by tfplugindocs)
  push-data.json   # Create this — valid JSON to push
```

Create `push-data.json` with your desired content before running `terraform apply`.

## Provider configuration

Ensure the Commander provider is configured (e.g. in `provider.tf` or in the same directory):

```terraform
provider "commander" {
  service_mode_url     = "https://your-commander-service.example.com"
  service_mode_api_key = "your-api-key"
}
```

See [Keeper Commander Service Mode](https://docs.keeper.io/en/keeperpam/commander-cli/service-mode-rest-api) for details.
