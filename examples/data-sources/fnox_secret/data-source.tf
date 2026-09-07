data "fnox_secret" "db_password" {
  key = "DB_PASSWORD"
}

data "fnox_secret" "prod_api_key" {
  key     = "API_KEY"
  profile = "production"
}

resource "some_resource" "example" {
  password = data.fnox_secret.db_password.value
}
