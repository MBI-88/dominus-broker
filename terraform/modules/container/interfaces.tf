variable "container_image" {
  type = string
}

variable "container_name" {
  type = string
}

variable "container_ports" {
  type = list(object({
    internal = number
    external = number
  }))
}

variable "network_mode" {
  type = string
}

variable "network_name" {
  type = string
}
