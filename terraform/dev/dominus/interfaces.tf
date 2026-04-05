
variable "dominus_image_name" {
  type    = string
  default = "dominus-broker:local"
}

variable "dominus_file" {
  type        = string
  description = "Path to the Docker build context (repo root with Dockerfile), relative to Terraform cwd or absolute"
}

variable "dominus_container_name" {
  type    = string
  default = "dominus_broker"
}

variable "dominus_container_ports" {
  type = list(object({
    internal = number
    external = number
  }))
}

variable "dominus_volume_cert" {
  type        = string
  description = "Name cert"
  default     = "dominus_certs"
}

variable "dominus_volume_driver" {
  type    = string
  default = "local"
}

variable "dominus_network_name" {
  type = string
}

variable "dominus_container_cpu_resources" {
  type    = string
  default = "4.5"
}

variable "dominus_container_memory_resources" {
  type    = number
  default = 10
}

variable "dominus_container_healcheck" {
  type    = list(string)
  default = ["CMD", "wget", "-qO-", "--header", "x-api-key: dominus_example_@10102024KeyServerToken", "http://127.0.0.1:8000/health"]
}
