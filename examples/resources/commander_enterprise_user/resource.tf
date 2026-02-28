# Enterprise user: email and node (required). Optional: name, job_title, roles, teams, managed_company (MSP only).
# id and status are read-only (set by the API).

resource "commander_enterprise_user" "example" {
  email     = "alice@example.com"
  name      = "Alice Smith"
  job_title = "Software Engineer"
  node      = "Engineering"

  roles = ["Developer", "1234567890"]
  teams = ["Backend Developers", "1234567890"]

  # Optional, MSP only: scope to a specific managed company
  # managed_company = "Acme Corp"
}
