

variable "dominus_image_name" {
  type = string
  default = "dominus-broker:latest"
}

variable "dominus_platform" {
  type = string
  default = "golang:1.26.1-alpine3.23.3"
}

variable "dominus_file" {
  type = string

}

variable "dominus_dependencies" {
  type = list(string)
}

variable "dominus_container_name" {
  type = string
  default = "dominus_broker"
}

variable "dominus_container_ports" {
  type = list(object({
    internal = number
    external = number
  }))
  default = [
    { internal = 8000, external = 8000 },
    { internal = 5000, external = 5000 },
  ]
}

variable "dominus_volume_cert" {
  type = string
  description = "Name cert"
}

variable "dominus_volume_cert_driver" {
  type = string
  default = "local"
}

variable "dominus_volume_env" {
  type = string
  description = "Name env"
}

variable "dominus_volume_env_driver" {
  type        = string
  default     = "local"
  description = "Docker volume driver for the env volume"
}
