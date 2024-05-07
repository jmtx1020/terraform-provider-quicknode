# creates a gateway
resource "quicknode_gateway" "gateway" {
  name    = var.gateway_name
  private = true
  enabled = false
}