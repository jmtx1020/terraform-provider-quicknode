---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "quicknode Provider"
subcategory: ""
description: |-
  
---

# quicknode Provider



## Example Usage

```terraform
# Configuration-based authentication
provider "quicknode" {
  host  = "https://api.quicknode.com"
  token = "TOKEN_VALUE"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `host` (String) API Hostname
- `token` (String, Sensitive) API Token to use to authenticate.
