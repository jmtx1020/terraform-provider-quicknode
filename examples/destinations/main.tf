terraform {
  required_providers {
    quicknode = {
      source = "jose.com/edu/quicknode"
    }
  }
}

provider "quicknode" {
  host  = "https://api.quicknode.com"
  token = "QN_800bff829847439eba41e7867c79af68"
}

data "quicknode_destinations" "edu" {}
