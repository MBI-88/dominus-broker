
output "grafana_image_id" {
  value = module.grafana_image.image_id
}

output "grafana_container_ports" {
  value =  jsonencode(var.grafana_container_ports)
}

output "grafana_container_name" {
  value = module.grafana_container.container_name
}
