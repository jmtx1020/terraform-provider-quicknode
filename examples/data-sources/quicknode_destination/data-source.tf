# Fetch One Destination's Information
resource "quicknode_destination" "one" {
  name         = "au-test-api"
  to           = "https://us-central1-serious-truck-412423.cloudfunctions.net/function-1"
  webhook_type = "POST"
  service      = "webhook"
  payload_type = 1
}

data "quicknode_destination" "one" {
  id = resource.quicknode_destination.one.id
}
