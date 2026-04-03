


resource "docker_image" "image" {
  name = var.image
  keep_locally = var.keep_localy
  platform = var.platform
  depends_on = var.dependencies

  provisioner "file" {
    destination = var.path
  }
}
