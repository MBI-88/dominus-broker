output "network_driver" {
  value = docker_network.network.driver
}

output "network_name" {
  value = docker_network.network.name
}

output "network_attachable" {
  value = docker_network.network.attachable
}
