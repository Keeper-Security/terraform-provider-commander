# Import is supported. Use either the node name/ID only (logged-in account) or comma-separated managed_company,node.
# Node only (logged-in account):
terraform import commander_enterprise_node.example "Root"
terraform import commander_enterprise_node.example 1169425105420462
# With managed company (managed_company_name_or_id,node_name_or_id):
terraform import commander_enterprise_node.example "Test Company,Root"
terraform import commander_enterprise_node.example "1169425105420462,1169425105420462"

# Or use the import block in configuration:
# import {
#   to = commander_enterprise_node.example
#   id = "Root"
# }
# import {
#   to = commander_enterprise_node.example
#   id = "Test Company,Root"
# }
