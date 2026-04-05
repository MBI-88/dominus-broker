output "redis_image_id" {
  value = module.redis_image.image_id
}

output "redis_container_ports" {
  value = jsonencode(var.redis_container_ports)
}

output "redis_container_name" {
  value = module.redis_container.container_name
}
