# Import an existing EPM policy by Keeper policy ID.
# After import, Terraform runs Read to populate the rest of the attributes.
#
# Policy ID only (current account / MSP context):
terraform import commander_epm_policy.example "your-epm-policy-id"
#
# MSP: managed company name or ID, then policy ID (comma-separated):
terraform import commander_epm_policy.example "Acme Corp,your-epm-policy-id"
#
# Or use an import block in configuration:
# import {
#   to = commander_epm_policy.example
#   id = "your-epm-policy-id"
# }
# import {
#   to = commander_epm_policy.example
#   id = "Acme Corp,your-epm-policy-id"
# }
