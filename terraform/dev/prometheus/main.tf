module "prometheus_image" {
  source = "../../modules/image"

  image = var.prometheus_image_name
  path  = null
}

module "prometheus_volume" {
  source = "../../modules/volume"

  volume_name   = var.prometheus_volume_name
  volume_driver = var.prometheus_volume_driver
}

module "prometheus_container" {
  source = "../../modules/container"

  container_image        = module.prometheus_image.image_id
  container_name         = var.prometheus_container_name
  container_ports        = var.prometheus_container_ports
  container_network_name = var.prometheus_network_name
  container_healcheck    = var.prometheus_container_healcheck

  container_bind_mounts = [
    {
      host_path      = abspath("${path.module}/${var.prometheus_config_file}")
      container_path = "/etc/prometheus/prometheus.yml"
      read_only      = true
    },
  ]

  container_volume_mounts = [
    {
      volume_name    = module.prometheus_volume.volume_name
      container_path = "/prometheus"
      read_only      = false
    },
  ]
}
