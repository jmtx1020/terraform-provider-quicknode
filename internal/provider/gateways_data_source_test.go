package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGatewaysDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// this test is split into two parts because the resource wasn't showing in the data
			{
				Config: providerConfig + fmt.Sprintf(`
					resource "quicknode_gateway" "test" {
					  name    = "test-gateway-%s"
					  private = true
					  enabled = false
					}
				`, randomString(length)),
			},
			{
				Config: providerConfig + `
				data "quicknode_gateways" "test" {}
				`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.quicknode_gateways.test", "gateways.0.id"),
					resource.TestCheckResourceAttrSet("data.quicknode_gateways.test", "gateways.0.name"),
					resource.TestCheckResourceAttrSet("data.quicknode_gateways.test", "gateways.0.status"),
					resource.TestCheckResourceAttrSet("data.quicknode_gateways.test", "gateways.0.enabled"),
					resource.TestCheckResourceAttrSet("data.quicknode_gateways.test", "gateways.0.private"),
					resource.TestCheckResourceAttrSet("data.quicknode_gateways.test", "gateways.0.updated_at"),
					resource.TestCheckResourceAttrSet("data.quicknode_gateways.test", "gateways.0.created_at"),
				),
			},
		},
	})
}
