package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestNotificationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "quicknode_destination" "test" {
		  			name         = "au-test-api"
					to           = "https://us-central1-serious-truck-412423.cloudfunctions.net/function-1"
					webhook_type = "POST"
					service      = "webhook"
					payload_type = 1
				}
				resource "quicknode_notification" "test" {
					name            = "test_notification"
					network         = "ethereum-mainnet"
					expression      = "dHhfdG8gPT0gJzB4ZDhkYTZiZjI2OTY0YWY5ZDdlZWQ5ZTAzZTUzNDE1ZDM3YWE5NjA0Nic="
					destination_ids = [resource.quicknode_destination.test.id]
					enabled         = true
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("quicknode_notification.test", "id"),
					resource.TestCheckResourceAttrSet("quicknode_notification.test", "destination_ids.0"),
					resource.TestCheckResourceAttr("quicknode_notification.test", "name", "test_notification"),
					resource.TestCheckResourceAttr("quicknode_notification.test", "network", "ethereum-mainnet"),
					resource.TestCheckResourceAttr("quicknode_notification.test", "enabled", "true"),
				),
			},
		},
	})
}
