variable "image" {
  type        = string
  description = "Image name:tag to assign to the built image (e.g. dominus-broker:latest)"
}

variable "keep_localy" {
  type    = bool
  default = false
}

variable "path" {
  type        = string
  description = "Build context directory path (relative to the Terraform working directory or absolute)"
}

variable "dockerfile" {
  type        = string
  default     = "Dockerfile"
  description = "Dockerfile path relative to the build context"
}
