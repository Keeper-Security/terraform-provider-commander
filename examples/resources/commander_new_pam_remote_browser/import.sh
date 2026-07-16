# Import is supported. Import ID is the vault record UID of the pamRemoteBrowser record.
terraform import commander_new_pam_remote_browser.intranet_app "6HyuRTiqOljXl0J3DB6TAw"

# Or use the import block in configuration:
# import {
#   to = commander_new_pam_remote_browser.intranet_app
#   id = "6HyuRTiqOljXl0J3DB6TAw"
# }
#
# After import, run terraform plan and align configuration with remote state
# (title, url, pam_remote_browser_settings, share, etc.).
