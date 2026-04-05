// ==== Dominus outputs ==== //
output "dominus_container_name" {
  value = module.dominus_server.dominus_container_name
}

output "dominus_image_id" {
  value = module.dominus_server.dominus_image_id
}

output "dominus_ports" {
  value = jsonencode(module.dominus_server.dominus_container_ports)
}



// ==== Grafana ==== //
output "grafana_container_name" {
  value = module.grafana_server.grafana_container_name
}

output "grafana_image_id" {
  value = module.grafana_server.grafana_image_id
}

output "grafana_container_ports" {
  value = jsondecode(module.grafana_server.grafana_container_ports)
}



// ==== Redis ==== //
output "redis_container_name" {
  value = module.redis_server.redis_container_name
}

output "redis_image_id" {
  value = module.redis_server.redis_image_id
}

output "redis_container_ports" {
  value = jsondecode(module.redis_server.redis_container_ports)
}



// ==== Prometheus ==== //
output "prometheus_container_ports" {
  value = jsondecode(module.prometheus_server.prometheus_container_ports)
}

output "prometheus_image_id" {
  value = module.prometheus_server.prometheus_image_id
}

output "prometheus_container_name" {
  value = module.prometheus_server.prometheus_container_name
}



// ==== Sidecar ==== //
output "sidecar_container_ports" {
  value = jsondecode(module.sidecar_server.sidecar_container_ports)
}

output "sidecar_image_id" {
  value = module.sidecar_server.sidecar_image_id
}

output "sidecar_container_name" {
  value = module.sidecar_server.sidecar_container_name
}
