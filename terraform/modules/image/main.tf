resource "docker_image" "image" {
  name         = var.image
  keep_locally = var.keep_localy

  dynamic "build" {
    for_each = var.path != null && var.path != "" ? [1] : []
    content {
      context    = var.path
      dockerfile = var.dockerfile
    }
  }
}
