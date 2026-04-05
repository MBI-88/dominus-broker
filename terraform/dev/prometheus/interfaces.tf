variable "prometheus_image_name" {
  type    = string
  default = "prom/prometheus:latest"
}

variable "prometheus_container_name" {
  type    = string
  default = "prometheus"
}

variable "prometheus_config_file" {
  type        = string
  default     = "prometheus.yml"
  description = "Config file relative to this module directory; bind-mounted read-only"
}

variable "prometheus_volume_name" {
  type        = string
  default     = "prometheus_data"
  description = "Docker volume name for TSDB data (Compose top-level volume name)"
}

variable "prometheus_volume_driver" {
  type    = string
  default = "local"
}

variable "prometheus_container_ports" {
  type = list(object({
    internal = number
    external = number
  }))
  default = [
    { internal = 9090, external = 9090 },
  ]
}

variable "prometheus_network_name" {
  type = string
}

variable "prometheus_container_cpu_resources" {
  type    = string
  default = "2.5"
}

variable "prometheus_container_healcheck" {
  type    = list(string)
  default = ["CMD", "wget", "-qO-", "http://prometheus:9090/-/ready"]
}
