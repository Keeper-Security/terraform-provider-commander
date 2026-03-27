// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisepush

// ResourceDescription is the plain-text description for the commander_enterprise_push resource schema.
// It is used in Schema.Description (e.g. for tfplugindocs and Terraform Registry).
const ResourceDescription = "Push a JSON file of records into selected users' or teams' vaults. Provide the file path and at least one email or team. Adding more emails or teams later pushes only to the new recipients; removing recipients only updates Terraform state (no push). Changing the file or its path runs a full push to everyone again. Optional for MSP: managed_company."

// ResourceMarkdownDescription is the markdown description for the commander_enterprise_push resource schema.
// It is used in Schema.MarkdownDescription and becomes the intro content in generated docs (docs/resources/enterprise_push.md).
const ResourceMarkdownDescription = "Use this resource to **push record data** from a JSON file into your users' or teams' Keeper vaults. You choose the file and who receives it (by email and/or team).\n\n" +
	"## What this resource does\n\n" +
	"- On **first apply**, it pushes the file contents to everyone you list in **email** and **team**.\n" +
	"- Terraform does **not** read or delete anything on the server—it only runs the push and tracks it.\n" +
	"- If you **add** more emails or teams later, Terraform pushes **only to the new recipients** (existing ones are not pushed again).\n" +
	"- If you **remove** emails or teams, Terraform only updates its state; **no push** runs.\n" +
	"- If you **change the file** (same or different path), Terraform runs a **full push** to all current recipients again.\n\n" +
	"## What you need to provide\n\n" +
	"- **file_path** — Path to a JSON file on the machine running Terraform. The file must be valid JSON (see example below).\n" +
	"- **email** and/or **team** — At least one is required. List the users (by email) and/or teams that should receive the records.\n" +
	"- **managed_company** — (Optional, MSP only) Name or ID of the managed company. Omit to use your current account.\n\n" +
	"## Example JSON file\n\n" +
	"Create a file (e.g. `push-data.json`) with your records. It must be **valid JSON**. The whole file is sent as the push payload.\n\n" +
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
	"Point **file_path** at this file (e.g. `\"${path.module}/push-data.json\"`).\n\n" +
	"See [Enterprise Push documentation](https://docs.keeper.io/en/keeperpam/commander-cli/command-reference/enterprise-management-commands#enterprise-push-command) for more details.\n\n" +
	"## Adding or removing recipients\n\n" +
	"- **Add** one or more emails or teams → Terraform updates in place → push goes **only to the newly added** emails/teams.\n" +
	"- **Remove** one or more emails or teams → Terraform updates in place → **no push** (state only).\n" +
	"- **Change** the file or file_path → Terraform replaces the resource → **all** current email and team get the push.\n" +
	"- **Change** managed_company (MSP) → Terraform replaces the resource → **all** current email and team in the new company get the push.\n\n" +
	"**Tip:** When you add people or teams, only they receive the push. When you remove people or teams, nothing is pushed—Terraform just stops tracking them.\n\n" +
	"## When does a push run?\n\n" +
	"- **First apply** — Push runs to everyone in email and team; resource is created.\n" +
	"- **No change** to file or config — No push.\n" +
	"- **You add** email(s) or team(s) — Push runs only to the new email(s)/team(s).\n" +
	"- **You remove** email(s) or team(s) — No push; only Terraform state is updated.\n" +
	"- **You edit** the file or change file_path — Full push runs again to all current email and team.\n" +
	"- **You remove** the resource from config and apply — Resource is removed from state only; nothing is deleted on the server.\n\n" +
	"## Import\n\n" +
	"Import is not supported. This resource has no server-side identity. To run the same push again with the same settings, add the resource again with the same **file_path**, **email**, **team**, and **managed_company**; Terraform will use the same computed `id` and will not re-push if the file content is unchanged."
