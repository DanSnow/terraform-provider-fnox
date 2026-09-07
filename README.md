# terraform-provider-fnox

A Terraform provider that exposes secrets managed by [fnox](https://fnox.jdx.dev) as Terraform-native sensitive values.

[fnox](https://github.com/jdx/fnox) is a CLI secret manager that can read secrets from many backends (1Password, AWS Secrets Manager/SSM, Azure Key Vault, GCP Secret Manager, HashiCorp Vault, Bitwarden, age, and more), all described in a single `fnox.toml` config. This provider shells out to the `fnox` CLI so you can reference any fnox-managed secret directly from Terraform.

## Requirements

- The `fnox` binary must be installed and available on `$PATH` (or pointed to via `binary_path`).
- A `fnox.toml` configuration file describing your secrets and their backend(s).

## Usage

```hcl
terraform {
  required_providers {
    fnox = {
      source = "DanSnow/fnox"
    }
  }
}

provider "fnox" {
  # config_path = "${path.module}/fnox.toml"
  # profile     = "production"
}

data "fnox_secret" "db_password" {
  key = "DB_PASSWORD"
}

resource "some_resource" "example" {
  password = data.fnox_secret.db_password.value
}
```

See [docs/](docs/) for the full provider and data source reference, and [examples/](examples/) for more usage examples.

## Development

```console
go build ./...
go test ./...
TF_ACC=1 go test ./... -run TestAcc -v
```

Docs are generated with [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs):

```console
tfplugindocs generate
```
