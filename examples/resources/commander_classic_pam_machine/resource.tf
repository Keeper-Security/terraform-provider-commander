resource "commander_classic_pam_machine" "example" {
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
  folder_location  = "_REPLACE_WITH_SHARED_FOLDER_UID_"

  # ----------------------------------------------------------------
  # Per-user share permissions (optional).
  # Map key = user email. Each value is { can_share, can_edit }.
  # Both flags default to false (view-only).
  # ----------------------------------------------------------------
  share = {
    "alice@example.com" = {
      can_share = true
      can_edit  = true
    }
    "bob@example.com" = {
      can_edit = true
    }
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
      protocol          = "ssh" # one of: kubernetes, rdp, ssh, telnet, vnc
      connection_port   = 22
      launch_credential = "_REPLACE_WITH_LAUNCH_CREDENTIAL_UID_"

      # Only the block matching "protocol" should be uncommented.

      # -----------------------------------------------------------------------
      # Kubernetes
      # -----------------------------------------------------------------------
      # kubernetes {
      #   session_recording      = true
      #   typescript_recording   = true
      #   recording_include_keys = true
      #   allow_supply_user      = true
      #   read_only              = false
      #   rotate_on_termination  = true
      #   use_ssl                = true
      #   ignore_cert            = false
      #   namespace              = "production"
      #   pod                    = "api-server-7d9f8c6b4-xk2zp"
      #   container              = "app"
      #   command                = "/bin/sh"
      #   color_scheme           = "green-black"    # black-white | gray-black | green-black | white-black | Guacamole syntax
      #   font_name              = "monospace"
      #   font_size              = 14               # 8,9,10,11,12,14,18,24,30,36,48,60,72,96
      #   scrollback             = 1000
      #   backspace              = "127"             # "127" (default) | "8"
      #   ca_cert                = <<-EOT
      #     -----BEGIN CERTIFICATE-----
      #     MIIBkTCB+wIJAKExampleFakeCA...
      #     -----END CERTIFICATE-----
      #   EOT
      #   client_cert            = <<-EOT
      #     -----BEGIN CERTIFICATE-----
      #     MIIBkTCB+wIJAKExampleFakeClient...
      #     -----END CERTIFICATE-----
      #   EOT
      #   client_key             = <<-EOT
      #     -----BEGIN PRIVATE KEY-----
      #     MIIEvQIBADANBgkqhkExampleFakeKey...
      #     -----END PRIVATE KEY-----
      #   EOT
      # }

      # -----------------------------------------------------------------------
      # MySQL (postgresql and sql_server share the same attributes)
      # -----------------------------------------------------------------------
      # mysql {
      #   session_recording      = true
      #   typescript_recording   = true
      #   recording_include_keys = true
      #   allow_supply_user      = true
      #   read_only              = false
      #   database               = "production_db"
      #   color_scheme           = "white-black"
      #   font_name              = "monospace"
      #   font_size              = 12
      #   scrollback             = 500
      #   disable_copy           = false
      #   disable_paste          = false
      #   disable_csv_export     = true
      #   disable_csv_import     = true
      # }

      # -----------------------------------------------------------------------
      # PostgreSQL
      # -----------------------------------------------------------------------
      # postgresql {
      #   session_recording      = true
      #   typescript_recording   = true
      #   recording_include_keys = true
      #   allow_supply_user      = true
      #   read_only              = false
      #   database               = "analytics"
      #   color_scheme           = "black-white"
      #   font_name              = "monospace"
      #   font_size              = 12
      #   scrollback             = 500
      #   disable_copy           = false
      #   disable_paste          = false
      #   disable_csv_export     = false
      #   disable_csv_import     = false
      # }

      # -----------------------------------------------------------------------
      # SQL Server
      # -----------------------------------------------------------------------
      # sql_server {
      #   session_recording      = true
      #   typescript_recording   = true
      #   recording_include_keys = true
      #   allow_supply_user      = true
      #   read_only              = true
      #   database               = "finance_prod"
      #   color_scheme           = "black-white"
      #   font_name              = "monospace"
      #   font_size              = 12
      #   scrollback             = 500
      #   disable_copy           = false
      #   disable_paste          = false
      #   disable_csv_export     = false
      #   disable_csv_import     = false
      # }

      # -----------------------------------------------------------------------
      # RDP
      # -----------------------------------------------------------------------
      # rdp {
      #   session_recording      = true
      #   recording_include_keys = true
      #   allow_supply_user      = true
      #   read_only              = false
      #
      #   # Display
      #   color_depth            = 24               # 8 (default) | 16 | 24 | 32
      #   server_layout          = "en-us-qwerty"   # en-us-qwerty (default) | en-gb-qwerty | de-de-qwertz | fr-fr-azerty | fr-ch-qwertz | it-it-qwerty | ja-jp-qwerty | pt-br-qwerty | es-es-qwerty | sv-se-qwerty | tr-tr-qwerty | failsafe
      #   dpi                    = 96
      #   width                  = 1920
      #   height                 = 1080
      #   resize_method          = "display-update"  # display-update (default) | reconnect
      #   force_lossless         = false
      #
      #   # Visual experience
      #   enable_wallpaper           = true
      #   enable_theming             = true
      #   enable_font_smoothing      = true
      #   enable_desktop_composition = true
      #   enable_full_window_drag    = false
      #   enable_menu_animations     = false
      #
      #   # Caching
      #   disable_bitmap_caching     = false
      #   disable_offscreen_caching  = false
      #   disable_glyph_caching      = false
      #
      #   # Security & clipboard
      #   security            = "nla"               # any (default) | nla | tls | vmconnect | rdp
      #   ignore_cert         = false
      #   disable_auth        = false
      #   normalize_clipboard = "preserve"           # preserve (default) | unix | windows
      #   disable_copy        = false
      #   disable_paste       = false
      #
      #   # Audio
      #   console_audio      = false
      #   disable_audio      = false
      #   enable_audio_input = true
      #
      #   # Printing
      #   enable_printing         = true
      #   redirected_printer_name = "FollowMe-Printer"
      #
      #   # Remote app
      #   remote_app      = "C:\\Program Files\\App\\app.exe"
      #   remote_app_dir  = "C:\\Program Files\\App"
      #   remote_app_args = "--config prod"
      #
      #   # Misc
      #   enable_touch    = false
      #   console         = false
      #   timezone        = "America/New_York"
      #   client_name     = "terraform-jumpbox"
      #   initial_program = ""
      #
      #   # Pre-connection
      #   load_balance_info  = ""
      #   preconnection_id   = ""
      #   preconnection_blob = ""
      #
      #   # SFTP file transfer (sftp_user_uid is required when enable_sftp = true)
      #   sftp {
      #     enable_sftp                = true
      #     sftp_resource_uid          = "_REPLACE_WITH_SFTP_RESOURCE_UID_"
      #     sftp_user_uid              = "_REPLACE_WITH_SFTP_USER_UID_"
      #     sftp_directory             = "/uploads"
      #     sftp_server_alive_interval = 30
      #   }
      # }

      # -----------------------------------------------------------------------
      # SSH
      # -----------------------------------------------------------------------
      ssh {
        session_recording      = true
        typescript_recording   = true
        recording_include_keys = true
        allow_supply_user      = true
        read_only              = false
        color_scheme           = "black-white" # black-white (default) | gray-black | green-black | white-black | Guacamole syntax
        font_name              = "monospace"
        font_size              = 12 # 8,9,10,11,12,14,18,24,30,36,48,60,72,96
        scrollback             = 2000
        backspace              = "127" # "127" (default) | "8"
        terminal_type          = "xterm-256color"
        locale                 = "en_US.UTF-8"      # default: "$LANG"
        timezone               = "America/New_York" # default: "$TZ"
        server_alive_interval  = 60
        command                = "/bin/bash"
        disable_copy           = false
        disable_paste          = false
        host_key               = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExampleHostKey..."

        sftp {
          enable_sftp = true
        }
      }

      # -----------------------------------------------------------------------
      # Telnet
      # -----------------------------------------------------------------------
      # telnet {
      #   session_recording      = true
      #   typescript_recording   = true
      #   recording_include_keys = true
      #   allow_supply_user      = true
      #   read_only              = false
      #   color_scheme           = "green-black"
      #   font_name              = "monospace"
      #   font_size              = 11
      #   scrollback             = 500
      #   backspace              = "8"              # "127" (default) | "8"
      #   terminal_type          = "vt100"
      #   disable_copy           = false
      #   disable_paste          = false
      #   username_regex         = "[Uu]sername:"
      #   password_regex         = "[Pp]assword:"
      #   login_success_regex    = ">"
      #   login_failure_regex    = "% Login invalid|% Bad"
      # }

      # -----------------------------------------------------------------------
      # VNC
      # -----------------------------------------------------------------------
      # vnc {
      #   session_recording      = true
      #   recording_include_keys = true
      #   allow_supply_user      = true
      #   read_only              = false
      #   color_depth            = 24               # 8 | 16 | 24 | 32
      #   cursor                 = "local"          # local | remote
      #   clipboard_encoding     = "UTF-8"          # UTF-8 (default) | UTF-16 | ISO8859-1 | CP1252
      #   disable_copy           = false
      #   disable_paste          = false
      #   swap_red_blue          = false
      #   force_lossless         = false
      #
      #   # Audio (PulseAudio)
      #   enable_audio     = true
      #   audio_servername = "pulseaudio.corp.example.com"
      #
      #   # VNC gateway / repeater
      #   dest_host = "10.0.1.50"
      #   dest_port = 5900
      #
      #   # SFTP file transfer (sftp_user_uid is required when enable_sftp = true)
      #   sftp {
      #     enable_sftp                = true
      #     sftp_resource_uid          = "_REPLACE_WITH_SFTP_RESOURCE_UID_"
      #     sftp_user_uid              = "_REPLACE_WITH_SFTP_USER_UID_"
      #     sftp_directory             = "/home/lab/uploads"
      #     sftp_server_alive_interval = 60
      #   }
      # }
    }
  }
}
