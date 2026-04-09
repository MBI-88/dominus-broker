// ====== Dominus ===== //
variable "dominus_server_file" {
  type = string
}


variable "dominus_server_container_ports" {
  type = list(object({
    internal = number
    external = number
  }))
}

variable "dominus_server_container_cpu" {
  type = string
}

variable "dominus_server_container_memory" {
  type = number
}

variable "dominus_server_token" {
  type = string
}

// ===== Grafana ==== //
variable "grafana_server_container_ports" {
  type = list(object({
    internal = number
    external = number
  }))
}

// ==== Redis ==== //
variable "redis_server_container_ports" {
  type = list(object({
    internal = number
    external = number
  }))
}


// ==== Prometheus ==== //
variable "prometheus_server_container_ports" {
  type = list(object({
    internal = number
    external = number
  }))
}


// ==== Sidecar ==== //
variable "sidecar_server_container_ports" {
  type = list(object({
    internal = number
    external = number
  }))
}


// ======== Shared ========= //
variable "network_mode" {
  type = string
}

variable "network_name" {
  type = string
}
