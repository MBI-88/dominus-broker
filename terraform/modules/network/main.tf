

resource "docker_network" "network"{
  name = var.name
  driver = var.driver
  attachable = var.attachable
}
