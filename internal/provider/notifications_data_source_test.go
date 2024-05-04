package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// read testing
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
				data "quicknode_notifications" "test" {}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.quicknode_notifications.test", "notifications.0.id"),
					resource.TestCheckResourceAttrSet("data.quicknode_notifications.test", "notifications.0.name"),
					resource.TestCheckResourceAttrSet("data.quicknode_notifications.test", "notifications.0.network"),
					resource.TestCheckResourceAttrSet("data.quicknode_notifications.test", "notifications.0.enabled"),
					resource.TestCheckResourceAttrSet("data.quicknode_notifications.test", "notifications.0.expression"),
					resource.TestCheckResourceAttrSet("data.quicknode_notifications.test", "notifications.0.updated_at"),
					resource.TestCheckResourceAttrSet("data.quicknode_notifications.test", "notifications.0.created_at"),
				),
			},
		},
	})
}
