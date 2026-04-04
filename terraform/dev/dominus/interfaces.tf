

variable "dominus_image_name" {
  type = string
  default = "dominus-broker:latest"
}

variable "dominus_file" {
  type        = string
  description = "Path to the Docker build context (repo root with Dockerfile), relative to Terraform cwd or absolute"
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

variable "dominus_network_mode" {
  type = string
  default = "bridge"
}

variable "dominus_network_name" {
  type = string

}

variable "dominus_network_driver" {
  type = string
  default = "bridge"
}

variable "dominus_network_attachable" {
  type = bool
  default = true
}
