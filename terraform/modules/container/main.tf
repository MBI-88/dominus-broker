
resource "docker_container" "container" {
  image        = var.container_image
  name         = var.container_name
  env          = var.container_env
  network_mode = var.container_network_mode

  command = var.container_conmands
  memory  = var.container_memory_resources
  cpus    = var.container_cpu_resources

  healthcheck {
    interval = var.container_healcheck_interval
    timeout  = var.container_healcheck_timeout
    retries  = var.container_healcheck_retries
    test     = var.container_healcheck
  }

  networks_advanced {
    name = var.container_network_name
  }

  dynamic "ports" {
    for_each = var.container_ports
    content {
      internal = ports.value.internal
      external = ports.value.external
    }
  }

  dynamic "volumes" {
    for_each = var.container_bind_mounts
    content {
      host_path      = volumes.value.host_path
      container_path = volumes.value.container_path
      read_only      = volumes.value.read_only
    }
  }

  dynamic "volumes" {
    for_each = var.container_volume_mounts
    content {
      volume_name    = volumes.value.volume_name
      container_path = volumes.value.container_path
      read_only      = volumes.value.read_only
    }
  }
}
