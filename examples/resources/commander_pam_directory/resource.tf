# ------------------------------------------------------------------
# Example 1 – Minimal: Active Directory with only required fields
# ------------------------------------------------------------------
resource "commander_pam_directory" "minimal" {
  title = "PAM Directory - corp-dc-01"

  hostname_or_ip = {
    hostname            = "dc-01.corp.example.com"
    administrative_port = 636
  }
}

# ------------------------------------------------------------------
# Example 2 – Active Directory with all directory fields
# ------------------------------------------------------------------
resource "commander_pam_directory" "active_directory" {
  title = "PAM Directory - corp-dc-01"

  hostname_or_ip = {
    hostname            = "dc-01.corp.example.com"
    administrative_port = 636
  }

  use_ssl        = true
  domain_name    = "corp.example.com"
  directory_type = "active_directory" # active_directory | openldap

  alternative_ips = [
    "10.0.1.10",
    "10.0.1.11",
  ]

  directory_id    = "d-926724afcc"
  user_match      = "OU=ServiceAccounts,DC=corp,DC=example,DC=com"
  provider_group  = "Azure"
  provider_region = "us-east-1"
  notes           = "Primary domain controller for corp.example.com"
  folder          = "_REPLACE_WITH_SHARED_FOLDER_UID_"
}

# ------------------------------------------------------------------
# Example 3 – OpenLDAP directory
# ------------------------------------------------------------------
resource "commander_pam_directory" "openldap" {
  title = "PAM Directory - ldap-server"

  hostname_or_ip = {
    hostname            = "ldap.internal.example.com"
    administrative_port = 389
  }

  use_ssl        = false
  domain_name    = "internal.example.com"
  directory_type = "openldap"
  user_match     = "/ou=People/"
  notes          = "Internal OpenLDAP server for user discovery."
}

# ------------------------------------------------------------------
# Example 4 – Full configuration with PAM settings (SSH connection)
# ------------------------------------------------------------------
resource "commander_pam_directory" "full" {
  title = "PAM Directory - prod-dc"

  hostname_or_ip = {
    hostname            = "prod-dc.corp.example.com"
    administrative_port = 636
  }

  use_ssl        = true
  domain_name    = "corp.example.com"
  directory_type = "active_directory"

  alternative_ips = [
    "10.0.2.10",
    "10.0.2.11",
    "10.0.2.12",
  ]

  directory_id    = "d-926724afcc"
  user_match      = "OU=Users,DC=corp,DC=example,DC=com"
  provider_group  = "AWS"
  provider_region = "us-west-2"
  notes           = "Production domain controller with full PAM settings."
  folder          = "_REPLACE_WITH_SHARED_FOLDER_UID_"

  pam_settings {
    configuration              = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"
    administrative_credentials = "_REPLACE_WITH_ADMIN_CREDENTIAL_UID_"
    allow_supply_host          = false

    tunnel {
      enable                   = true
      remote_target_port       = 636
      re_use_port              = true
      use_specified_local_port = true
      local_port               = 10636
    }

    connection {
      enable            = true
      protocol          = "ssh"
      connection_port   = 22
      launch_credential = "_REPLACE_WITH_LAUNCH_CREDENTIAL_UID_"

      ssh {
        session_recording      = true
        typescript_recording   = true
        recording_include_keys = true
        allow_supply_user      = false
        read_only              = false
        color_scheme           = "black-white"
        font_name              = "monospace"
        font_size              = 12
        scrollback             = 2000
        backspace              = "127"
        terminal_type          = "xterm-256color"
        locale                 = "en_US.UTF-8"
        timezone               = "America/New_York"
        server_alive_interval  = 60
        command                = "/bin/bash"
        disable_copy           = false
        disable_paste          = false

        sftp {
          enable_sftp = true
        }
      }
    }
  }
}

# ------------------------------------------------------------------
# Example 5 – PAM settings with tunnel only (no connection)
# ------------------------------------------------------------------
resource "commander_pam_directory" "tunnel_only" {
  title = "PAM Directory - tunnel-dc"

  hostname_or_ip = {
    hostname = "tunnel-dc.corp.example.com"
  }

  domain_name    = "corp.example.com"
  directory_type = "active_directory"

  pam_settings {
    configuration = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"

    tunnel {
      enable             = true
      remote_target_port = 389
    }
  }
}
