

module "dominus_image" {
  source = "../../modules/image"

  image       = var.dominus_image_name
  keep_localy = false
  path        = var.dominus_file
}

module "dominus_container" {
  source = "../../modules/container"

  container_image = module.dominus_image.docker_image
  container_name  = var.dominus_container_name
  container_ports  = var.dominus_container_ports
  network_mode = var.dominus_network_mode
  network_name = var.dominus_network_name

}

module "dominus_volume" {
  source   = "../../modules/volume"
  for_each = local.docker_volumes

  name   = each.value.name
  driver = each.value.driver
}
