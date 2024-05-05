package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGatewayDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "quicknode_gateway" "test" {
				  name    = "test-gateway-%s"
				  private = true
				  enabled = false
				}
				data "quicknode_gateway" "test" {
					name = resource.quicknode_gateway.test.name
				}
				`, randomString(length)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.quicknode_gateway.test", "id"),
					resource.TestCheckResourceAttrSet("data.quicknode_gateway.test", "name"),
					resource.TestCheckResourceAttrSet("data.quicknode_gateway.test", "enabled"),
					resource.TestCheckResourceAttrSet("data.quicknode_gateway.test", "private"),
					resource.TestCheckResourceAttrSet("data.quicknode_gateway.test", "status"),
					resource.TestCheckResourceAttrSet("data.quicknode_gateway.test", "updated_at"),
					resource.TestCheckResourceAttrSet("data.quicknode_gateway.test", "created_at"),
				),
			},
		},
	})
}
