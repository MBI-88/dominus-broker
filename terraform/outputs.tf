// ==== Dominus outputs ==== //
output "dominus_container_name" {
  value = module.dominus_server.dominus_container_name
}

output "dominus_image_id" {
  value = module.dominus_server.dominus_image_name
}

output "dominus_ports" {
  value = jsonencode(module.dominus_server.docker_container_ports)

}
