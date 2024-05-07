# Fetch One Notification's Information
resource "quicknode_destination" "destination" {
  name         = var.destination_name
  to           = var.endpoint_url
  webhook_type = "POST"
  service      = "webhook"
  payload_type = 1
}

resource "quicknode_notification" "test" {
  name            = var.notification_name
  network         = "ethereum-mainnet"
  expression      = "dHhfdG8gPT0gJzB4ZDhkYTZiZjI2OTY0YWY5ZDdlZWQ5ZTAzZTUzNDE1ZDM3YWE5NjA0Nic="
  destination_ids = [resource.quicknode_destination.destination.id]
  enabled         = true
}

# retrieves one notification by id
data "quicknode_notification" "notification" {
  id = resource.quicknode_notification.test.id
}