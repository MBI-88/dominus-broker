
variable "sidecar_image_name" {
  type    = string
  default = "nginx:latest"
}

variable "sidecar_container_name" {
  type    = string
  default = "sidecar"
}

variable "sidecar_container_ports" {
  type = list(object({
    internal = number
    external = number
  }))
}

variable "sidecar_network_name" {
  type = string
}

variable "sidecar_environments" {
  type    = list(string)
  default = ["NGINX_HOST=0.0.0.0", "NGINX_PORT=80", "NGINX_ENVSUBST_TEMPLATE_SUFFIX=.conf"]
}

variable "sidecar_path" {
  type    = string
  default = "./nginx.conf"
}

variable "sidecar_container_healcheck" {
  type    = list(string)
  default = ["CMD", "curl", "-f", "http://127.0.0.1:80/"]
}
