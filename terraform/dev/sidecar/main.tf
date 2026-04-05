
module "sidecar_image" {
  source = "../../modules/image"

  image = var.sidecar_image_name
  path  = null
}

module "sidecar_container" {
  source = "../../modules/container"

  container_image         = module.sidecar_image.image_id
  container_name          = var.sidecar_container_name
  container_ports         = var.sidecar_container_ports
  container_network_name  = var.sidecar_network_name
  container_env           = var.sidecar_environments
  container_cpu_resources = var.sidecar_container_cpu_resources
  container_healcheck     = var.sidecar_container_healcheck

  container_bind_mounts = [
    {
      host_path      = abspath("${path.module}/${var.sidecar_path}")
      container_path = "/etc/nginx/nginx.conf"
      read_only      = true
    }

  ]
}
