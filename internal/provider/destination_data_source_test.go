package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDestinationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
				resource "quicknode_destination" "one" {
		  			name         = "au-test-api"
					to           = "https://us-central1-serious-truck-412423.cloudfunctions.net/function-1"
					webhook_type = "POST"
					service      = "webhook"
					payload_type = 1
				}
				data "quicknode_destination" "one" {
					id = resource.quicknode_destination.one.id
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the first coffee to ensure all attributes are set
					resource.TestCheckResourceAttrSet("data.quicknode_destination.one", "id"),
					resource.TestCheckResourceAttrSet("data.quicknode_destination.one", "name"),
					resource.TestCheckResourceAttrSet("data.quicknode_destination.one", "payload_type"),
					resource.TestCheckResourceAttrSet("data.quicknode_destination.one", "service"),
					resource.TestCheckResourceAttrSet("data.quicknode_destination.one", "to"),
					resource.TestCheckResourceAttrSet("data.quicknode_destination.one", "token"),
					resource.TestCheckResourceAttrSet("data.quicknode_destination.one", "updated_at"),
					resource.TestCheckResourceAttrSet("data.quicknode_destination.one", "created_at"),
				),
			},
		},
	})
}
