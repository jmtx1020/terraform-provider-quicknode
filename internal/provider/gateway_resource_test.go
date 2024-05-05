package provider

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	charset = "0123456789"
	length  = 6
)

func TestGatewayResource(t *testing.T) {
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
				`, randomString(length)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("quicknode_gateway.test", "id"),
					resource.TestCheckResourceAttrSet("quicknode_gateway.test", "enabled"),
					resource.TestCheckResourceAttrSet("quicknode_gateway.test", "private"),
					resource.TestCheckResourceAttrSet("quicknode_gateway.test", "status"),
					resource.TestCheckResourceAttrSet("quicknode_gateway.test", "updated_at"),
					resource.TestCheckResourceAttrSet("quicknode_gateway.test", "created_at"),
				),
			},
		},
	})
}

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
