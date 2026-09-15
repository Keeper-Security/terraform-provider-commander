###############################################################################
# Example 1 - Minimal: open guest network with no password.
###############################################################################

resource "commander_classic_wifi" "guest" {
  title = "Guest WiFi"
  ssid  = "Company-Guest"
  notes = "Open network in the lobby"
}

###############################################################################
# Example 2 - Full: WPA home network with hidden SSID, custom fields, and
# per-user sharing.
#
# `encryption` accepts exactly: wep | wpa | noEncryption.
#
# `share` is a map keyed by email address. Each value sets:
#   can_share — let the user re-share the record with others
#   can_edit  — let the user edit the record
# Removing an email from `share` on a subsequent apply revokes that user's
# access automatically.
###############################################################################

resource "commander_classic_wifi" "home" {
  title           = "Home WiFi"
  folder_location = "Personal"
  ssid            = "MyHomeNetwork"
  password        = "WiFiPassword123"
  encryption      = "wpa"
  is_ssid_hidden  = true
  notes           = "Living-room router"


  custom = [
    {
      type  = "text"
      label = "Frequency"
      value = "5GHz"
    },
    {
      type  = "text"
      label = "Router"
      value = "Eero Pro 6"
    },
  ]

  share = {
    "alice@example.com" = {
      can_share = true
      can_edit  = true
    }
    "bob@example.com" = {
      can_share = false
      can_edit  = false
    }
  }
}

###############################################################################
# Example 3 - Explicit open network (no encryption).
###############################################################################

resource "commander_classic_wifi" "lab" {
  title           = "Lab Open WiFi"
  folder_location = "Lab"
  ssid            = "lab-open"
  encryption      = "noEncryption"
  is_ssid_hidden  = false
}

###############################################################################
# Outputs - record UIDs are useful for chaining into data sources or other
# resources (e.g. shared folder_locations).
###############################################################################

output "wifi_home_id" {
  description = "UID of the managed Home WiFi record."
  value       = commander_classic_wifi.home.id
}
