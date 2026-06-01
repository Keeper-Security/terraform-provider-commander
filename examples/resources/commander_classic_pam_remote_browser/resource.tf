resource "commander_classic_pam_remote_browser" "intranet_app" {
  title = "PAM RBI - intranet.example.com"
  url   = "https://intranet.example.com/app/login"

  notes  = "Internal HR portal; RBI + autofill + URL allow lists."
  folder = "Shared Folders/PAM/Remote Browser" # optional; use folder UID or path Commander accepts

  pam_remote_browser_settings = {
    # Required when this nested block is present — PAM configuration that owns RBI for this record.
    configuration = "_REPLACE_WITH_PAM_CONFIGURATION_UID_"

    remote_browser_isolation = true
    connections_recording    = true
    key_events               = true
    allow_url_navigation     = true
    ignore_server_cert       = true

    allowed_urls = toset([
      "localhost.com",
      "127.0.0.1:9000",
    ])

    allowed_resource_urls = toset([
      "localhost.com",
      "127.0.0.1:9000",
    ])

    auto_fill_credentials = "_REPLACE_WITH_CREDENTIALS_RECORD_UID_"
    auto_fill_targets = toset([
      "_REPLACE_WITH_AUTOFILL_TARGET_RECORD_UID_",
    ])

    allow_copy    = true
    allow_paste   = true
    disable_audio = false

    audio_channels    = 2
    audio_bit_depth   = 16
    audio_sample_rate = 48000
  }
}