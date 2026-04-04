
output "container_name" {
  value = docker_container.container.name
}

output "container_ports" {
  value = jsonencode(docker_container.container.ports)
}

output "container_network" {
  value = docker_container.container.network.name
}
