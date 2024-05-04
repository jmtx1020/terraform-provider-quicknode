package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDestinationsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
				resource "quicknode_destination" "test" {
		  			name         = "au-test-api"
					to           = "https://us-central1-serious-truck-412423.cloudfunctions.net/function-1"
					webhook_type = "POST"
					service      = "webhook"
					payload_type = 1
				}
				data "quicknode_destinations" "test" {}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttrSet("data.quicknode_destinations.test", "destinations.0.id"),
					resource.TestCheckResourceAttrSet("data.quicknode_destinations.test", "destinations.0.name"),
					resource.TestCheckResourceAttrSet("data.quicknode_destinations.test", "destinations.0.payload_type"),
					resource.TestCheckResourceAttrSet("data.quicknode_destinations.test", "destinations.0.service"),
					resource.TestCheckResourceAttrSet("data.quicknode_destinations.test", "destinations.0.to"),
					resource.TestCheckResourceAttrSet("data.quicknode_destinations.test", "destinations.0.token"),
					resource.TestCheckResourceAttrSet("data.quicknode_destinations.test", "destinations.0.updated_at"),
					resource.TestCheckResourceAttrSet("data.quicknode_destinations.test", "destinations.0.created_at"),
				),
			},
		},
	})
}
