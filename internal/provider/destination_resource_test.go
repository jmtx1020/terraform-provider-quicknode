package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDestinationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
					resource "quicknode_destination" "test" {
						name         = "ds-tf-testing-api"
						to           = "https://us-central1-serious-truck-412423.cloudfunctions.net/function-1"
						webhook_type = "POST"
						service      = "webhook"
						payload_type = 1
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first coffee item has Computed attributes filled.
					resource.TestCheckResourceAttr("quicknode_destination.test", "name", "ds-tf-testing-api"),
					resource.TestCheckResourceAttr("quicknode_destination.test", "to", "https://us-central1-serious-truck-412423.cloudfunctions.net/function-1"),
					resource.TestCheckResourceAttr("quicknode_destination.test", "webhook_type", "POST"),
					resource.TestCheckResourceAttr("quicknode_destination.test", "service", "webhook"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("quicknode_destination.test", "id"),
					resource.TestCheckResourceAttrSet("quicknode_destination.test", "token"),
					resource.TestCheckResourceAttrSet("quicknode_destination.test", "updated_at"),
					resource.TestCheckResourceAttrSet("quicknode_destination.test", "created_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "quicknode_destination.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					resource "quicknode_destination" "test" {
						name         = "ds-tf-testing-update"
						to           = "https://us-central1-serious-truck-412423.cloudfunctions.net/function-1"
						webhook_type = "POST"
						service      = "webhook"
						payload_type = 1
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first destination has attributes updated.
					resource.TestCheckResourceAttr("quicknode_destination.test", "name", "ds-tf-testing-update"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
