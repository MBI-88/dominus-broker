
variable "redis_image_name" {
  type = string
  default = "redis-dev:local"
}

variable "redis_container_name" {
  type    = string
  default = "redis"
}

variable "redis_volume_name" {
  type    = string
  default = "redis_data"
}

variable "redis_volume_driver" {
  type    = string
  default = "local"
}

variable "redis_container_ports" {
  type = list(object({
    internal = number
    external = number
  }))
}

variable "redis_network_name" {
  type = string
}

variable "redis_container_conmand" {
  type = list(string)
  default = [ "redis-server", "/usr/local/etc/redis/redis.conf" ]
}

variable "redis_path" {
  type = string
  default = "./redis.conf"
}

variable "redis_container_cpu_resources" {
  type = string
  default = "2.5"
}

variable "redis_container_healcheck" {
  type = list(string)
  default = ["CMD", "redis-cli", "-p", "6379", "ping"]
}
