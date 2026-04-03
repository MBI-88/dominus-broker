resource "docker_image" "image" {
  name         = var.image
  keep_locally = var.keep_localy

  build {
    context    = var.path
    dockerfile = var.dockerfile
  }
}
