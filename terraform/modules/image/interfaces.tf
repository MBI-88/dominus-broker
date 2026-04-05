variable "image" {
  type        = string
  description = "Image name:tag to assign to the built image (e.g. dominus-broker:latest)"
  default     = null
  nullable    = true
}

variable "keep_localy" {
  type    = bool
  default = false
}

variable "path" {
  type        = string
  default     = null
  nullable    = true
  description = "Build context path; null or empty = only pull var.image from the registry"
}

variable "dockerfile" {
  type        = string
  default     = "Dockerfile"
  description = "Dockerfile path relative to the build context"
}
