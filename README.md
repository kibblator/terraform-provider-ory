<div align="center">
  <h1>Ory Terraform Provider</h1>

[![Release](https://img.shields.io/github/v/release/kibblator/ory-terraform-provider?logo=terraform&include_prereleases)](https://github.com/kibblator/ory-terraform-provider/releases)
[![License](https://img.shields.io/github/license/kibblator/ory-terraform-provider.svg)](https://github.com/kibblator/ory-terraform-provider/blob/main/LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/kibblator/ory-terraform-provider/ci-cd.yml?branch=main)](https://github.com/kibblator/ory-terraform-provider/actions?query=branch%3Amain)

</div>

-------------------------------------

The is a Terraform Provider for managing Ory Network configuration through [Terraform](https://www.terraform.io/).

-------------------------------------

## Documentation

TODO

## Getting Started

### Requirements

- [Terraform](https://www.terraform.io/downloads)
- An [Ory Network](https://ory.sh) account

### Installation

Terraform uses the [Terraform Registry](https://registry.terraform.io/) to download and install providers. To install
this provider, copy and paste the following code into your Terraform configuration. Then, run `terraform init`.

```terraform
terraform {
  required_providers {
    ory = {
      source = "registry.terraform.io/kibblator/ory"
    }
  }
}

provider "ory" {}
```

```shell
$ terraform init
