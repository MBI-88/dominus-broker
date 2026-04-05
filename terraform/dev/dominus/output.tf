
output "dominus_image_id" {
  value = module.dominus_image.image_id
}

output "dominus_container_name" {
  value = module.dominus_container.container_name
}

output "dominus_container_ports" {
  value = jsonencode(var.dominus_container_ports)
}
