
variable "grafana_image_name" {
  type    = string
  default = "grafana/grafana:latest"
}

variable "grafana_container_name" {
  type    = string
  default = "grafana"
}

variable "grafana_volume_name" {
  type    = string
  default = "grafana_data"
}

variable "grafana_volume_driver" {
  type    = string
  default = "local"
}

variable "grafana_container_ports" {
  type = list(object({
    internal = number
    external = number
  }))
}

variable "grafana_container_env" {
  type        = list(string)
  default     = [
      "GF_SECURITY_ADMIN_PASSWORD=admin",
      "GF_SECURITY_ADMIN_USER=admin",
      "GF_AUTH_ANONYMOUS_ENABLED=false",
      "GF_INSTALL_PLUGINS=grafana-clock-panel,grafana-simple-json-datasource"
  ]
  sensitive   = true
  description = "KEY=value strings for the container (Grafana: GF_*). Use tfvars; marked sensitive for passwords."
}

variable "grafana_network_name" {
  type = string
}

variable "grafana_container_cpu_resources" {
  type = string
  default = "2.5"
}

variable "grafana_container_healcheck" {
  type = list(string)
  default = ["CMD", "wget", "-qO-", "http://grafana:3000/login"]
}
