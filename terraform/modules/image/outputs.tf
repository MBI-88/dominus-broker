output "image_id" {
  value       = docker_image.image.image_id
  description = "Image ID for docker_container.image (or use image_name if the engine prefers a tag)"
}

output "image_name" {
  value       = docker_image.image.name
  description = "Configured image name:tag"
}

output "image_tag" {
  value       = docker_image.image.name
  description = "Same as image_name; kept for backward compatibility"
}
