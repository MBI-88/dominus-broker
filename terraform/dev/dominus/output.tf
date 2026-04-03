
output "dominus_image_name" {
  value = module.dominus_image
}

output "dominus_container_name" {
  value = module.dominus_container.container_name
}

output "docker_container_ports" {
  value = jsonencode(module.dominus_container.container_ports)
}
