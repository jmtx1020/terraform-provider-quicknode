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