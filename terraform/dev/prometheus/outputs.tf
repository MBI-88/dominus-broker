
output "prometheus_image_id" {
  value = module.prometheus_image.image_id
}

output "prometheus_container_ports" {
  value = jsonencode(var.prometheus_container_ports)
}

output "prometheus_container_name" {
  value = module.prometheus_container.container_name
}
