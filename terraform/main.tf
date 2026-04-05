
module "network" {
  source = "./modules/network"

  network_name       = var.network_name
  network_driver     = var.network_mode
  network_attachable = true

}

module "dominus_server" {
  source = "./dev/dominus"

  dominus_file            = var.dominus_server_file
  dominus_container_ports = var.dominus_server_container_ports
  dominus_network_name    = module.network.network_name

  depends_on = [module.redis_server]
}

module "grafana_server" {
  source = "./dev/grafana"

  grafana_network_name    = module.network.network_name
  grafana_container_ports = var.grafana_server_container_ports

  depends_on = [module.dominus_server, module.prometheus_server, module.sidecar_server]
}

module "prometheus_server" {
  source = "./dev/prometheus"

  prometheus_container_ports = var.prometheus_server_container_ports
  prometheus_network_name    = module.network.network_name

  depends_on = [module.dominus_server, module.sidecar_server]
}

module "redis_server" {
  source = "./dev/redis"

  redis_network_name    = module.network.network_name
  redis_container_ports = var.redis_server_container_ports
}

module "sidecar_server" {
  source = "./dev/sidecar"

  sidecar_container_ports = var.sidecar_server_container_ports
  sidecar_network_name    = module.network.network_name

  depends_on = [module.dominus_server]
}
