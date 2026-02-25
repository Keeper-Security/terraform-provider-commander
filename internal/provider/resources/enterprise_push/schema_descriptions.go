// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisepush

// ResourceDescription is the plain-text description for the commander_enterprise_push resource schema.
// It is used in Schema.Description (e.g. for tfplugindocs and Terraform Registry).
const ResourceDescription = "One-time action resource that pushes JSON file content to user vaults via the enterprise-push Commander CLI. Write-only: the API does not support read or delete. email and team are updatable: adding emails/teams triggers Update and pushes only to newly added targets; removing emails/teams only updates state (no push). Changes to file_path, file content, or managed_company trigger replace (destroy + create) and a full push."

// ResourceMarkdownDescription is the markdown description for the commander_enterprise_push resource schema.
// It is used in Schema.MarkdownDescription and becomes the intro content in generated docs (docs/resources/enterprise_push.md).
// This file is not shown in Terraform Registry; only the string value emitted into the schema is included in the generated docs.
const ResourceMarkdownDescription = "One-time action resource that pushes JSON file content to user vaults via the **enterprise-push** Commander CLI.\n\n" +
	"## Write-only resource\n\n" +
	"- The API does **not** support read or delete operations.\n" +
	"- Terraform **cannot** detect drift. If the file or targets change outside Terraform, Terraform will not see it.\n" +
	"- **email** and **team** are **updatable**: adding emails/teams triggers **Update** and pushes only to **newly added** targets; **removing** emails/teams only updates state (no push).\n" +
	"- Changes to **file_path**, file content, or **managed_company** trigger **replace** (destroy + create) and a full push.\n" +
	"- **Delete** only removes the resource from Terraform state; nothing is deleted on the server.\n\n" +
	"## Behavior\n\n" +
	"- **Deterministic ID**: `id` is computed as `sha256(file_content + sorted_emails + sorted_teams + managed_company)`.\n" +
	"- **Create**: Push runs to all configured email/team. **Update** (when only email/team change): push runs **only to newly added** emails/teams; if only removals, state is updated and **no push** runs.\n" +
	"- **No-op Read**: Read does not call the API; state is left as-is.\n" +
	"- **No-op Delete**: Delete does not call the API; the resource is only removed from state.\n\n" +
	"## Example push-data.json\n\n" +
	"Create a file that will be pushed (e.g. `push-data.json`). It must be **valid JSON**. The entire file content is sent as `filedata` to the enterprise-push command.\n\n" +
	"```json\n" +
	"{\n" +
	"  \"records\": [{\n" +
	"    \"title\":\"Google\",\n" +
	"    \"folders\": [\n" +
	"      {\n" +
	"        \"folder\": \"My Websites\\\\Online\"\n" +
	"      }\n" +
	"    ],\n" +
	"    \"login\": \"testing\",\n" +
	"    \"password\": \"lk4j139sk4j\",\n" +
	"    \"login_url\": \"https://google.com\",\n" +
	"    \"notes\": \"These are some notes.\",\n" +
	"    \"custom_fields\": {\"Favorite Food\":\"Cheetos\"}\n" +
	"  },\n" +
	"  {\n" +
	"    \"title\":\"Facebook\",\n" +
	"    \"folders\": [\n" +
	"      {\n" +
	"        \"folder\": \"Social Media\"\n" +
	"      },\n" +
	"      {\n" +
	"        \"shared_folder\": \"Shared Social\",\n" +
	"        \"can_edit\": false,\n" +
	"        \"can_share\": false\n" +
	"      }\n" +
	"    ],\n" +
	"    \"login\": \"me@gmail.com\",\n" +
	"    \"password\": \"123123123123\",\n" +
	"    \"login_url\": \"https://facebook.com\",\n" +
	"    \"notes\": \"This is our corporate shared record.\",\n" +
	"    \"custom_fields\": {\n" +
	"      \"Facebook Application ID\":\"ABC12345\",\n" +
	"      \"$oneTimeCode\": \"otpauth://totp/Amazon:me@company.com?secret=JBSWY3DPEHPK3PXP&issuer=Amazon&algorithm=SHA1&digits=6&period=30\"\n" +
	"    }\n" +
	"  }]\n" +
	"}\n" +
	"```\n\n" +
	"Place the file where Terraform can read it (e.g. `\"${path.module}/push-data.json\"`).\n\n" +
	"See [Enterprise Push documentation](https://docs.keeper.io/en/keeperpam/commander-cli/command-reference/enterprise-management-commands#enterprise-push-command) for more information.\n\n" +
	"## Adding or removing email/team: Update behavior\n\n" +
	"**email** and **team** do **not** use ForceNew. When you change only those attributes, Terraform runs **Update** (not replace).\n\n" +
	"- **When you add** one or more emails or teams: Update runs and pushes the file content **only to the newly added** emails/teams.\n" +
	"- **When you remove** emails or teams: Update runs and **only updates state**; no push is performed.\n\n" +
	"### Example: adding one email\n\n" +
	"**Step 1 — Initial config and apply**\n\n" +
	"```terraform\n" +
	"resource \"commander_enterprise_push\" \"example\" {\n" +
	"  file_path = \"${path.module}/data.json\"\n" +
	"  email     = [\"alice@example.com\", \"bob@example.com\"]\n" +
	"  team      = [\"Engineering\"]\n" +
	"}\n" +
	"```\n\n" +
	"- First `terraform apply`: push runs to **alice@example.com**, **bob@example.com**, and **Engineering**. Resource is created.\n\n" +
	"**Step 2 — Add one email and apply again**\n\n" +
	"```terraform\n" +
	"resource \"commander_enterprise_push\" \"example\" {\n" +
	"  file_path = \"${path.module}/data.json\"\n" +
	"  email     = [\"alice@example.com\", \"bob@example.com\", \"carol@example.com\"]  # added carol\n" +
	"  team      = [\"Engineering\"]\n" +
	"}\n" +
	"```\n\n" +
	"- `terraform plan`: Terraform sees `email` changed → **in-place update** (no replace).\n" +
	"- `terraform apply`: **Update** runs. Push runs **only to carol@example.com** (the newly added email). alice, bob, and Engineering do **not** receive another push.\n\n" +
	"### Example: adding one team\n\n" +
	"If you add `[\"Support\"]` to `team`, Update runs and the push goes **only to the Support team** (and any other newly added emails/teams). Existing targets are not pushed again.\n\n" +
	"### Example: removing email/team\n\n" +
	"If you remove `bob@example.com` from `email` (and leave alice and Engineering), **Update** runs and **only state is updated**; no push is performed.\n\n" +
	"### Summary\n\n" +
	"| Change | What Terraform does | Who gets the push |\n" +
	"|--------|---------------------|-------------------|\n" +
	"| Add email(s) and/or team(s) | Update | **Only newly added** emails/teams |\n" +
	"| Remove email(s) and/or team(s) | Update | **No push** (state only) |\n\n" +
	"## When does a push happen?\n\n" +
	"| Situation | Result |\n" +
	"|-----------|--------|\n" +
	"| First `terraform apply` with this resource | Push runs to all configured email/team; resource created. |\n" +
	"| No change to file or attributes | No re-push. |\n" +
	"| File content changed (same path) | Replace → push runs again to **all** current email/team. |\n" +
	"| `file_path` or `managed_company` changed | Replace → push runs again to all email/team. |\n" +
	"| **Add** email(s) and/or team(s) | Update → push runs **only to newly added** targets. |\n" +
	"| **Remove** email(s) and/or team(s) | Update → **no push** (state only). |\n" +
	"| Resource removed from config and applied | Resource removed from state only; no API delete. |\n\n" +
	"## All scenarios: how the resource behaves\n\n" +
	"Below are **all** scenarios and what Terraform and the provider do. Attributes that use **ForceNew** (`file_path`, file content via `file_content_sha256`, `managed_company`) trigger **replace** (destroy + create). Only **email** and **team** are updatable in place.\n\n" +
	"| # | Scenario | Terraform operation | API / push behavior | State after apply |\n" +
	"|---|----------|---------------------|----------------------|-------------------|\n" +
	"| 1 | **First time** adding the resource and applying | **Create** | Read file → push to **all** configured email + team | Resource created; `id` and `file_content_sha256` set |\n" +
	"| 2 | **No change** to config or file, apply again | No change | No API call | Unchanged |\n" +
	"| 3 | **Add one or more emails** (same file_path, content, team, managed_company) | **Update** | Push **only to the newly added** email(s) | State updated with new email set and new `id` |\n" +
	"| 4 | **Add one or more teams** (same file_path, content, email, managed_company) | **Update** | Push **only to the newly added** team(s) | State updated with new team set and new `id` |\n" +
	"| 5 | **Add both emails and teams** in one change | **Update** | Push **only to all newly added** emails and teams (one push with all new targets) | State updated with new email/team sets and new `id` |\n" +
	"| 6 | **Remove one or more emails** (keep at least one target in config) | **Update** | **No push**; only state updated | State updated with smaller email set and new `id` |\n" +
	"| 7 | **Remove one or more teams** (keep at least one target in config) | **Update** | **No push**; only state updated | State updated with smaller team set and new `id` |\n" +
	"| 8 | **Remove some emails and some teams** (still have at least one email or team) | **Update** | **No push**; only state updated | State updated; new `id` |\n" +
	"| 9 | **Add emails and remove teams** (or add teams and remove emails) in same change | **Update** | Push **only to newly added** emails/teams; removals do not trigger push | State updated; new `id` |\n" +
	"| 10 | **File content changed** (same `file_path`, file edited on disk) | **Replace** (destroy + create) | Destroy: no-op. Create: read new file → push to **all** current email + team | New resource instance; new `id`, new `file_content_sha256` |\n" +
	"| 11 | **`file_path` changed** (different file, same or different content) | **Replace** (destroy + create) | Destroy: no-op. Create: read new path → push to **all** current email + team | New resource instance; new `id`, new `file_content_sha256` |\n" +
	"| 12 | **`managed_company` changed** (MSP) | **Replace** (destroy + create) | Destroy: no-op. Create: push to **all** current email + team in new company context | New resource instance; new `id` |\n" +
	"| 13 | **File content + add email** in same apply | **Replace** (content change wins) | Destroy: no-op. Create: push to **all** plan email + team (including the new email) with **new** file content | New resource; full push to all targets with new content |\n" +
	"| 14 | **Refresh** (e.g. `terraform refresh` or during plan/apply) | **Read** | No API call; state left as-is (no-op Read) | Unchanged |\n" +
	"| 15 | **Resource removed from config** and apply | **Delete** | No API call; resource removed from state only | Resource gone from state |\n\n" +
	"### Notes on scenarios\n\n" +
	"- **Scenarios 3–5:** Only **new** emails/teams receive the push; existing targets are not pushed again.\n" +
	"- **Scenarios 6–8:** Removals never trigger a push; the provider only updates Terraform state.\n" +
	"- **Scenarios 10–12:** Any change to `file_path`, file content (detected via `file_content_sha256`), or `managed_company` causes **replace**. Create then pushes to **all** current email and team in the plan.\n" +
	"- **Scenario 13:** When replace is forced by file content (or file_path/managed_company), Create runs with the **full** plan (all emails/teams); there is no \"delta\" push on replace.\n" +
	"- **At least one target required on Create:** You must specify at least one of `email` or `team` when creating the resource; otherwise Terraform will report an error. On **Update**, you may remove all emails and teams (state is updated, no push is performed).\n\n" +
	"## Import\n\n" +
	"**Import is not supported.** This is a write-only resource with no server-side identity. To \"re-import\" a push, add the resource again with the same `file_path`, `email`, `team`, and `managed_company`; Terraform will compute the same `id` and will not re-push if the file content is unchanged."
