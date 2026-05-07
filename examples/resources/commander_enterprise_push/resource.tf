# commander_enterprise_push: one-time push of JSON file content to user vaults.
#
# Required: file_path — path to a valid JSON file (on the machine running Terraform).
# Optional: email, team (each must have at least one value if set), managed_company (MSP only).
# id is computed: sha256(content + sorted emails + sorted teams + managed_company).
#
# Write-only: no read/update/delete from API. Any config change causes replace → push again.
# Delete removes the resource from state only.

resource "commander_enterprise_push" "example" {
  file_path = "${path.module}/push-data.json"

  email = ["alice@example.com", "bob@example.com"]
  team  = ["Engineering", "Support"]

  # Optional, MSP only: scope to a specific managed company
  # managed_company = "Acme Corp"
}
