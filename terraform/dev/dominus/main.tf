

module "dominus_image" {
  source = "../../modules/image"

  image = var.dominus_image_name
  path  = var.dominus_file
}

module "dominus_container" {
  source = "../../modules/container"

  container_image            = module.dominus_image.image_id
  container_name             = var.dominus_container_name
  container_ports            = var.dominus_container_ports
  container_network_name     = var.dominus_network_name
  container_cpu_resources    = var.dominus_container_cpu_resources
  container_memory_resources = var.dominus_container_memory_resources
  container_healcheck        = var.dominus_container_healcheck

  container_volume_mounts = [
    {
      volume_name    = module.dominus_volume.volume_name
      container_path = "/etc/dominus/certs"
      read_only      = false
    },
  ]
}

module "dominus_volume" {
  source = "../../modules/volume"

  volume_name   = var.dominus_volume_cert
  volume_driver = var.dominus_volume_driver
}
