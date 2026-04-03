variable "dominus_server_file" {
  type = string

}
variable "dominus_server_container_name" {
  type = string
}

variable "dominus_server_container_ports" {
  type = list(object({
    internal = number
    external = number
  }))
  default = [
    { internal = 8000, external = 8000 },
    { internal = 5000, external = 5000 },
  ]
}

variable "dominus_server_volume_cert" {
  type = string
  default = "dominus_certs:/etc/dominus/certs:ro"
}

variable "dominus_server_volume_cert_driver" {
  type = string
  default = "local"
}

variable "dominus_server_volume_env" {
  type = string
  default = "./env:/env:ro"
}

variable "dominus_server_volume_env_driver" {
  type = string
  default = "local"
}

variable "network_mode" {
  type = string
}

variable "network_name" {
  type = string
}
