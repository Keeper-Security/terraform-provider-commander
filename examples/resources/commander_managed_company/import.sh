# Import is supported using the managed company name or ID.
terraform import commander_managed_company.example "Test Company"
# Or using the company ID:
terraform import commander_managed_company.example 1169425105420462

# Or use the import block in configuration:
# import {
#   to = commander_managed_company.example
#   id = "Test Company"
# }
# import {
#   to = commander_managed_company.example
#   id = "1169425105420462"
# }
