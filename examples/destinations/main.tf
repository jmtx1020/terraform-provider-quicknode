terraform {
  required_providers {
    quicknode = {
      source = "jose.com/edu/quicknode"
    }
  }
}

provider "quicknode" {}

data "quicknode_destinations" "edu" {}

output "edu_destinations" {
  value = data.quicknode_destinations.edu
}

resource "quicknode_destination" "edx" {
  name         = "tz-testing-go-api"
  to           = "https://us-central1-serious-truck-412423.cloudfunctions.net/function-1"
  webhook_type = "POST"
  service      = "webhook"
  payload_type = 1
}

resource "quicknode_destination" "edy" {
  name         = "tq-testing-go-api"
  to           = "https://us-central1-serious-truck-412423.cloudfunctions.net/function-1"
  webhook_type = "POST"
  service      = "webhook"
  payload_type = 1
}

resource "quicknode_destination" "edz" {
  name         = "tq-testing-go-api-v4"
  to           = "https://us-central1-serious-truck-412423.cloudfunctions.net/function-1"
  webhook_type = "POST"
  service      = "webhook"
  payload_type = 1
}

output "edz_id" {
  value = resource.quicknode_destination.edz.id
}

data "quicknode_destination" "single" {
  id = "1dd594ec-31b1-44cf-83ba-e6b891b130a1"
}

output "single_node" {
  value = data.quicknode_destination.single
}

data "quicknode_notifications" "nots" {}

output "notifications" {
  value = data.quicknode_notifications.nots
}

resource "quicknode_notification" "test" {
  name            = "test_notification"
  network         = "ethereum-mainnet"
  expression      = "dHhfdG8gPT0gJzB4ZDhkYTZiZjI2OTY0YWY5ZDdlZWQ5ZTAzZTUzNDE1ZDM3YWE5NjA0Nic="
  destination_ids = [resource.quicknode_destination.edz.id]
  enabled         = true
}

output "notification" {
  value = resource.quicknode_notification.test
}

data "quicknode_notification" "test" {
  id = resource.quicknode_notification.test.id
}

output "notification_test" {
  value = data.quicknode_notification.test
}

resource "quicknode_gateway" "test" {
  name    = "test-gateway-1008"
  private = true
  enabled = false
}

# output "gateway_test" {
#   value = resource.quicknode_gateway.test
# }

# data "quicknode_gateway" "test" {
#   name = resource.quicknode_gateway.test.name
# }

# data "quicknode_gateways" "test" {}
