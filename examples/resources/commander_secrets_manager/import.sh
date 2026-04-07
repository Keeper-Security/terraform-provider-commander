# Import using the application UID (from `secrets-manager app list` or the Keeper admin console).
terraform import commander_secrets_manager.example "IFWGfyQDSFEJErUU4wnZAA"

# Or use the import block in configuration:
# import {
#   to = commander_secrets_manager.example
#   id = "IFWGfyQDSFEJErUU4wnZAA"
# }
