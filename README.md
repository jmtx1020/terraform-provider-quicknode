# Terraform Provider QuickNode

## Using the provider

### Configuring the Provider

```hcl
terraform {
  required_providers {
    quicknode = {
      source = "jmtx1020/quicknode"
      version = "0.0.1"
    }
  }
}

provider "quicknode" {
  host  = "https://api.quicknode.com"
  token = "TOKEN_VALUE"
}
```

The provider can also read environment variables for those values, if you configure it like this:

```hcl
provider "quicknode" {}
```

Then run the terraform as follows:

```shell
QUICKNODE_API_HOST='https://api.quicknode.com' \
QUICKNODE_API_TOKEN='<TOKEN_VALUE>' \
terraform apply
```

Associated documentation for all resources and datasources can be found [here](https://registry.terraform.io/providers/jmtx1020/quicknode/latest/docs) on the terraform registry.

### Creating Resources

```hcl
resource "quicknode_destination" "destination" {
  name         = "dest"
  to           = var.endpoint_url
  webhook_type = "POST"
  service      = "webhook"
  payload_type = 1
}

# retrieves one destination by id
data "quicknode_destination" "dest" {
  id = resource.quicknode_destination.destination.id
}

# retrieves all destinations
data "quicknode_notifications" "dests" {}


resource "quicknode_notification" "test" {
  name            = var.notification_name
  network         = "ethereum-mainnet"
  expression      = "dHhfdG8gPT0gJzB4ZDhkYTZiZjI2OTY0YWY5ZDdlZWQ5ZTAzZTUzNDE1ZDM3YWE5NjA0Nic="
  destination_ids = [resource.quicknode_destination.destination.id]
  enabled         = true
}

# retrieves one destination by id
data "quicknode_notification" "notification" {
  id = resource.quicknode_notification.test.id
}

# gets all notifications
data "quicknode_notifications" "notifications" {}

# creates a gateway
resource "quicknode_gateway" "gateway" {
  name    = var.gateway_name
  private = true
  enabled = false
}

# gets one one gateway by name
data "quicknode_gateway" "gateway" {
  name = resource.quicknode_gateway.gateway.name
}

# gets all gateways
data "quicknode_gateways" "gateways" {}
```

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```
