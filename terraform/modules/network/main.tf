

resource "docker_network" "network"{
  name = var.network
  driver = var.driver
  attachable = var.attachable
}
