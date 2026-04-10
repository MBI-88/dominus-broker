network_driver = "bridge"
network_name   = "dominus"

dominus_server_file = "./../"
dominus_server_container_ports = [
  { internal = 8000, external = 8000 },
  { internal = 5000, external = 5000 },
]

dominus_server_container_cpu    = "4"
dominus_server_container_memory = 1000


grafana_server_container_ports = [
  { internal = 3000,
    external = 3000
  },
]

redis_server_container_ports = [{
  internal = 6379,
  external = 6379
}]


prometheus_server_container_ports = [{
  internal = 9090,
  external = 9090
}]

sidecar_server_container_ports = [{
  internal = 80,
  external = 80
}]
