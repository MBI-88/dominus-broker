terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 4.1.0"
    }
  }
  backend "local" {
    path = "terraform.tfstate"
  }
}

provider "docker" {
  host = "npipe:////.//pipe//docker_engine"
}
