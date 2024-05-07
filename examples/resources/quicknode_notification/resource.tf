resource "quicknode_notification" "notification" {
  name            = var.notification_name
  network         = "ethereum-mainnet"
  expression      = "dHhfdG8gPT0gJzB4ZDhkYTZiZjI2OTY0YWY5ZDdlZWQ5ZTAzZTUzNDE1ZDM3YWE5NjA0Nic="
  destination_ids = [resource.quicknode_destination.destination.id]
  enabled         = true
}