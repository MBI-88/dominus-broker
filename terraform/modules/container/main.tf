resource "docker_container" "container" {
  image = var.container_image
  name  = var.container_name

  dynamic "ports" {
    for_each = var.container_ports
    content {
      internal = ports.value.internal
      external = ports.value.external
    }
  }
}
