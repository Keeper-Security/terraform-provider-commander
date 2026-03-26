# Look up a managed company by name or ID. Returns id, name, node, node_name, plan, file_plan, seats, add_ons.

data "commander_managed_company" "example" {
  managed_company = "Acme Corp"
}

output "managed_company_id" {
  value = data.commander_managed_company.example.id
}

output "managed_company_name" {
  value = data.commander_managed_company.example.name
}

output "managed_company_node" {
  value = data.commander_managed_company.example.node
}

output "managed_company_node_name" {
  value = data.commander_managed_company.example.node_name
}

output "managed_company_plan" {
  value = data.commander_managed_company.example.plan
}

output "managed_company_file_plan" {
  value = data.commander_managed_company.example.file_plan
}

output "managed_company_seats" {
  value = data.commander_managed_company.example.seats
}

output "managed_company_add_ons" {
  value = data.commander_managed_company.example.add_ons
}


