
variable "image" {
  type = string
  description = "Docker image name"
}

variable "keep_localy" {
  type = bool
  default = false
}

variable "platform" {
  type = string
}

variable "dependencies" {
  type = list(string)
}

variable "path" {
  type = string
}
