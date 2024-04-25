terraform {
  required_providers {
    quicknode = {
      source = "jose.com/edu/quicknode"
    }
  }
}

provider "quicknode" {}

data "quicknode_destinations" "example" {}
