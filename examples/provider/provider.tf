terraform {
  required_providers {
    fnox = {
      source = "DanSnow/fnox"
    }
  }
}

provider "fnox" {
  # Optional: override the fnox binary, config file, working directory, or default profile.
  # binary_path = "/usr/local/bin/fnox"
  # config_path = "${path.module}/fnox.toml"
  # profile     = "production"
}
