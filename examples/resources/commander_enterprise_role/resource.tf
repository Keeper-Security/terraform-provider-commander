# Enterprise role: name and node (required). Optional: users, teams, managing_nodes, enforcement_policies, managed_company (MSP only).
# enforcement_policies: key = policy name (e.g. REQUIRE_TWO_FACTOR), value = string.
#   - GENERATED_PASSWORD_COMPLEXITY: use jsonencode([...]) with a non-empty array of rule objects.
#   - TWO_FACTOR_DURATION_*: one of login, 12_hours, 24_hours, 30_days, forever.
#   - KEEPER_FILL_*: one of enforce, disable, null.

resource "commander_enterprise_role" "example" {
  name = "Developer"
  node = "Engineering"

  users = ["alice@example.com", "bob@example.com"]
  teams = ["Backend Developers"]

  # Grant admin privileges on specific nodes (map key = node name/ID)
  managing_nodes = {
    "Engineering" = {
      privileges = ["manage_user", "manage_roles", "manage_teams"]
      cascade    = true
    }
  }

  # Enforcement policies: key = policy name, value = string. GPC must be non-empty JSON array.
  enforcement_policies = {
    "REQUIRE_TWO_FACTOR"             = "true"
    "MASTER_PASSWORD_MINIMUM_LENGTH" = "20"
    "TWO_FACTOR_DURATION_WEB"        = "24_hours" # login | 12_hours | 24_hours | 30_days | forever
    "KEEPER_FILL_AUTO_FILL"          = "enforce"  # enforce | disable | null
    "GENERATED_PASSWORD_COMPLEXITY" = jsonencode([
      {
        domains                 = ["_default_"]
        length                  = 40
        "lower-use"             = true
        "lower-min"             = 2
        "upper-use"             = true
        "upper-min"             = 2
        "digit-use"             = true
        "digit-min"             = 2
        "special-use"           = true
        "special-min"           = 2
        special                 = "!@#$%^?();'\",=+[]<>{}&"
        "passphrase-allow"      = true
        "passphrase-length"     = 5
        "passphrase-capitalize" = false
        "passphrase-number"     = false
        "passphrase-separator"  = "-"
        "apply-privacy-screen"  = true
      }
    ])
  }

  # Optional, MSP only
  # managed_company = "Acme Corp"
}
