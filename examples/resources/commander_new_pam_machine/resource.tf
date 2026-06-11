# ------------------------------------------------------------------
# Example 1 - Minimal: only required fields (title + hostname)
# ------------------------------------------------------------------
resource "commander_new_pam_machine" "minimal" {
  title = "PAM Machine - bastion-01"

  hostname_or_ip = {
    hostname            = "bastion-01.corp.example.com"
    administrative_port = 22
  }
}

# ------------------------------------------------------------------
# Example 2 - SSH machine with full PAM settings + share permissions
# ------------------------------------------------------------------
resource "commander_new_pam_machine" "example" {
  title = "PAM Machine - prod-server-01"

  hostname_or_ip = {
    hostname            = "prod-server-01.corp.example.com"
    administrative_port = 22
  }

  operating_system = "Ubuntu 22.04 LTS"
  instance_name    = "prod-server-01"
  instance_id      = "i-0abc123def456"
  provider_group   = "AWS"
  provider_region  = "us-east-1"
  notes            = "Primary production server."
  folder_location  = "_REPLACE_WITH_NSF_FOLDER_UID_"

  # ----------------------------------------------------------------
  # Per-user share permissions (optional).
  # Map key = user email; value = one of:
  #   viewer, share-manager, content-manager,
  #   content-share-manager, full-manager.
  # ----------------------------------------------------------------
  share = {
    "alice@example.com" = "full-manager"
    "bob@example.com"   = "content-manager"
    "carol@example.com" = "viewer"
  }

  pam_settings {
    configuration              = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"
    administrative_credentials = "_REPLACE_WITH_ADMIN_CREDENTIAL_UID_"
    allow_supply_host          = false

    tunnel {
      enable                   = true
      remote_target_port       = 22
      re_use_port              = true
      use_specified_local_port = true
      local_port               = 10022
    }

    connection {
      enable            = true
      protocol          = "ssh" # one of: kubernetes, mysql, postgresql, rdp, sql-server, ssh, telnet, vnc
      connection_port   = 22
      launch_credential = "_REPLACE_WITH_LAUNCH_CREDENTIAL_UID_"

      # Only the block matching "protocol" should be set; see the
      # commander_classic_pam_machine example for connection variants
      # (kubernetes, mysql, postgresql, sql_server, rdp, telnet, vnc).
      ssh {
        session_recording      = true
        typescript_recording   = true
        recording_include_keys = true
        allow_supply_user      = true
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
        host_key               = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExampleHostKey..."

        sftp {
          enable_sftp = true
        }
      }
    }
  }
}
