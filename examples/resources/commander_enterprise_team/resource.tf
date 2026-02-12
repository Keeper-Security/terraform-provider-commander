# Enterprise team: name and node (required). Optional: users, roles, record options, managed_company (MSP only).

resource "commander_enterprise_team" "example" {
  name = "Backend Developers"
  node = "Engineering"

  users = ["alice@example.com", "bob@example.com"]
  roles = ["Developer"]

  restrict_record_edit     = true
  restrict_record_re_share = true
  enable_privacy_screen    = false

  # Optional, MSP only
  # managed_company = "Acme Corp"
}
