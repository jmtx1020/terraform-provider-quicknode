# Fetch One Destination's Information
resource "quicknode_destination" "dest" {
  name         = var.name
  to           = var.endpoint_url
  webhook_type = "POST"
  service      = "webhook"
  payload_type = 1
}

data "quicknode_destination" "dest" {
  id = resource.quicknode_destination.dest.id
}
