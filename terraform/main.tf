
module "dominus_server" {
  source = "./dev/dominus"
  dominus_file = var.dominus_server_file
  dominus_container_name = var.dominus_server_container_name
  dominus_container_ports = var.dominus_server_container_ports
  dominus_volume_cert = var.dominus_server_volume_cert
  dominus_volume_cert_driver = var.dominus_server_volume_cert_driver
  dominus_volume_env = var.dominus_server_volume_env
  dominus_volume_env_driver = var.dominus_server_volume_env_driver
  dominus_network_mode = var.network_mode
  dominus_network_name = var.network_name
}
