resource "docker_image" "image" {
  name         = var.image
  keep_locally = var.keep_localy
  force_remove = true

  dynamic "build" {
    for_each = var.path != null && var.path != "" ? [1] : []
    content {
      context    = var.path
      dockerfile = var.dockerfile
      no_cache   = true
      build_args = {
        GITHUB_TOKEN = var.github_token != null && var.github_token != "" ? var.github_token : ""
      }
    }
  }
}
