locals {
  # Volúmenes derivados de variables `dominus_volume_*`: cada bloque merge corresponde
  # a un grupo nombre + driver. Para uno nuevo, declara las variables y añade un mapa aquí.
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
