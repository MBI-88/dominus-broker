

resource "docker_volume" "volume" {
  name = var.name
  driver = var.driver
}
