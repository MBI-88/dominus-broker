output "sidecar_image_id" {
  value = module.sidecar_image.image_id
}

output "sidecar_container_ports" {
  value = jsonencode(var.sidecar_container_ports)
}

output "sidecar_container_name" {
  value = module.sidecar_container.container_name
}
