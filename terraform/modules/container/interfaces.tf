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

variable "container_env" {
  type        = list(string)
  default     = []
  description = "Environment entries as KEY=value (same idea as Compose `environment:`)"
}

variable "container_bind_mounts" {
  type = list(object({
    host_path      = string
    container_path = string
    read_only      = bool
  }))
  default     = []
  description = "Host path bind mounts (Compose `volumes` with a host path)"
}

variable "container_volume_mounts" {
  type = list(object({
    volume_name    = string
    container_path = string
    read_only      = bool
  }))
  default     = []
  description = "Named Docker volume mounts (volume must exist, e.g. from module.volume)"
}

variable "container_network_name" {
  type = string
}

variable "container_network_mode" {
  type    = string
  default = "bridge"
}

variable "container_conmands" {
  type    = list(string)
  default = null
}

variable "container_healcheck" {
  type        = list(string)
  description = "Conmand secuencies"
}

variable "container_cpu_resources" {
  type        = string
  default     = "2.5"
  description = "Specify how much of the available CPU resources a container can use. e.g a value of 1.5 means the container is guaranteed at most one and a half of the CPUs. Has precedence over cpu_period and cpu_quota."
}

variable "container_memory_resources" {
  type        = number
  default     = 3
  description = "The memory limit for the container in MBs."
}

variable "container_healcheck_interval" {
  type        = string
  default     = "1m"
  description = "Time between running the check (ms|s|m|h). Defaults to 0s."
}

variable "container_healcheck_timeout" {
  type        = string
  default     = "1s"
  description = "Maximum time to allow one check to run (ms|s|m|h). Defaults to 0s."
}

variable "container_healcheck_retries" {
  type        = number
  default     = 3
  description = "Consecutive failures needed to report unhealthy. Defaults to 0."
}
