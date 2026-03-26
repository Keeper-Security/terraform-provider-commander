# Enterprise team: name and node (required). Optional: users, roles, record options, managed_company (MSP only).

resource "commander_enterprise_team" "example" {
  name = "Backend Developers"
  node = "Engineering"

  users = ["alice@example.com", "1234567890"]
  roles = ["1234567890", "Developer"]

  restrict_record_edit     = true  # false by default
  restrict_record_re_share = true  # false by default
  enable_privacy_screen    = false # false by default

  # Optional, MSP only
  # managed_company = "Acme Corp"
}
