locals {
  docker_volumes = merge(
    {
      cert = {
        name   = var.dominus_volume_cert
        driver = var.dominus_volume_cert_driver
      }
    },
    {
      env = {
        name   = var.dominus_volume_env
        driver = var.dominus_volume_env_driver
      }
    },
  )
}
