
module "redis_image" {
  source = "../../modules/image"

  image = var.redis_image_name
  path  = abspath("${path.module}/")

}

module "redis_container" {
  source = "../../modules/container"

  container_image        = module.redis_image.image_id
  container_name         = var.redis_container_name
  container_network_name = var.redis_network_name
  container_ports        = var.redis_container_ports
  container_env          = null
  container_conmands     = var.redis_container_conmand
  container_healcheck    = var.redis_container_healcheck

  container_bind_mounts = [
    {
      host_path      = abspath("${path.module}/${var.redis_path}")
      container_path = "/usr/local/etc/redis/redis.conf"
      read_only      = false
    }
  ]

  container_volume_mounts = [
    {
      volume_name    = module.redis_volume.volume_name
      container_path = "/data"
      read_only      = false
    }
  ]

}

module "redis_volume" {
  source = "../../modules/volume"

  volume_name   = var.redis_volume_name
  volume_driver = var.redis_volume_driver
}
