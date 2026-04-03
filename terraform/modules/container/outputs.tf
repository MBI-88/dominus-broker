
output "container_name" {
  value = var.container_name
}

output "container_ports" {
  value = jsonencode(var.container_ports)
}
