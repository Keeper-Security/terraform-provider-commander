terraform {
  required_providers {
    commander = {
      source = "keeper-security/commander"
    }
  }
}

provider "commander" {
  service_mode_url     = "http://localhost:8080/api/v2/"
  service_mode_api_key = "XXXXXXXXXXXXXX"
}

# Minimal example: open guest network with no password.
resource "commander_wifi" "guest" {
  title = "Guest WiFi"
  ssid  = "Company-Guest"
  notes = "Open network in the lobby"
}

# Full example: WPA home network with hidden SSID and a custom field.
# `encryption` must be one of: wep, wpa, noEncryption.
resource "commander_wifi" "home" {
  title          = "Home WiFi"
  folder         = "Personal"
  ssid           = "MyHomeNetwork"
  password       = "WiFiPassword123"
  encryption     = "wpa"
  is_ssid_hidden = true
  notes          = "Living-room router"

  custom {
    type  = "text"
    label = "Frequency"
    value = "5GHz"
  }

  custom {
    type  = "text"
    label = "Router"
    value = "Eero Pro 6"
  }
}
