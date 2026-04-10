

module "grafana_image" {
  source = "../../modules/image"

  image = var.grafana_image_name
  path  = null
}


module "grafana_container" {
  source = "../../modules/container"

  container_env          = var.grafana_container_env
  container_image        = module.grafana_image.image_id
  container_name         = var.grafana_container_name
  container_ports        = var.grafana_container_ports
  container_network_name = var.grafana_network_name
  container_healcheck    = var.grafana_container_healcheck

  container_volume_mounts = [
    {
      volume_name    = module.grafana_volume.volume_name
      container_path = "/var/lib/grafana"
      read_only      = false
    },
  ]
}


module "grafana_volume" {
  source = "../../modules/volume"

  volume_name   = var.grafana_volume_name
  volume_driver = var.grafana_volume_driver
}
