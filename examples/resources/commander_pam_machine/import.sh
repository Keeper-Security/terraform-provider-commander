# Import is supported. Import ID is the vault record UID of the PAM machine record.
terraform import commander_pam_machine.example "AbCdEfGhIjKlMnOpQrStUw"

# Or use the import block in configuration:
# import {
#   to = commander_pam_machine.example
#   id = "AbCdEfGhIjKlMnOpQrStUw"
# }
#
# After import, run terraform plan and align configuration with remote state
# (title, hostname_or_ip, pam_settings, etc.).
