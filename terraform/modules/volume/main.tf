

resource "docker_volume" "volume" {
  name = var.volume_name
  driver = var.volume_driver
}
