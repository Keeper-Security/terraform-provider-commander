# Import is supported. Import ID is the vault record UID of the PAM directory record.
terraform import commander_new_pam_directory.active_directory "AbCdEfGhIjKlMnOpQrStUw"

# Or use the import block in configuration:
# import {
#   to = commander_new_pam_directory.active_directory
#   id = "AbCdEfGhIjKlMnOpQrStUw"
# }
#
# After import, run terraform plan and align configuration with remote state
# (title, hostname_or_ip, directory_type, pam_settings, share, etc.).
