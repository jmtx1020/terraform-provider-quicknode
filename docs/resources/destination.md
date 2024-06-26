---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "quicknode_destination Resource - quicknode"
subcategory: ""
description: |-
  
---

# quicknode_destination (Resource)



## Example Usage

```terraform
# Manage example destination.
resource "quicknode_destination" "destination" {
  name         = var.name
  to           = var.endpoint_url
  webhook_type = "POST"
  service      = "webhook"
  payload_type = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) User supplied name given to the destination.
- `payload_type` (Number) The type of payload to send. ENUM: 1,2,3,4,5,6,7
- `service` (String) The destination service. Currently only "webhook" is supported.
- `to` (String) The webhook URL to which QuickAlerts will send alert payloads.
- `webhook_type` (String) The type of destination. ENUM: 'POST', 'GET'

### Read-Only

- `created_at` (String) The date and time the destination was created.
- `id` (String) ID given by API for the destination.
- `token` (String) The token for this destination. This is used to optionally verify a QuickAlerts payload.
- `updated_at` (String) The date and time the destination was last updated.

## Import

Import is supported using the following syntax:

```shell
# Destination can be imported by specifying the ID in API.
terraform import quicknode_destination.destination $DESTINATION_ID
```
