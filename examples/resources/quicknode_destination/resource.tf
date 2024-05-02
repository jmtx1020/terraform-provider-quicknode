# Manage example destination.
resource "quicknode_destination" "edx" {
  name         = "tz-testing-go-api"
  to           = "https://us-central1-serious-truck-412423.cloudfunctions.net/function-1"
  webhook_type = "POST"
  service      = "webhook"
  payload_type = 1
}
