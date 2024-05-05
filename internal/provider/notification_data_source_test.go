package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestNotificationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read Testing
			{
				Config: providerConfig + `
				resource "quicknode_destination" "test" {
		  			name         = "test_destination"
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
				data "quicknode_notification" "test" {
					id = resource.quicknode_notification.test.id
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.quicknode_notification.test", "id"),
					resource.TestCheckResourceAttrSet("data.quicknode_notification.test", "enabled"),
					resource.TestCheckResourceAttrSet("data.quicknode_notification.test", "expression"),
					resource.TestCheckResourceAttrSet("data.quicknode_notification.test", "name"),
					resource.TestCheckResourceAttrSet("data.quicknode_notification.test", "network"),
					resource.TestCheckResourceAttrSet("data.quicknode_notification.test", "updated_at"),
					resource.TestCheckResourceAttrSet("data.quicknode_notification.test", "created_at"),
					resource.TestCheckResourceAttrSet("data.quicknode_notification.test", "destinations.0.id"),
					resource.TestCheckResourceAttrSet("data.quicknode_notification.test", "destinations.0.name"),
					resource.TestCheckResourceAttrSet("data.quicknode_notification.test", "destinations.0.payload_type"),
					resource.TestCheckResourceAttrSet("data.quicknode_notification.test", "destinations.0.service"),
					resource.TestCheckResourceAttrSet("data.quicknode_notification.test", "destinations.0.to"),
					resource.TestCheckResourceAttrSet("data.quicknode_notification.test", "destinations.0.webhook_type"),
				),
			},
		},
	})
}
