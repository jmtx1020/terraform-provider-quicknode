# Manage example destination.
resource "quicknode_destination" "destination" {
  name         = var.name
  to           = var.endpoint_url
  webhook_type = "POST"
  service      = "webhook"
  payload_type = 1
}
