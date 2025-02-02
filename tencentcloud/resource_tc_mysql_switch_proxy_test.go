package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixMysqlSwitchProxyResource_basic -v
func TestAccTencentCloudNeedFixMysqlSwitchProxyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlSwitchProxy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_switch_proxy.switch_proxy", "id"),
				),
			},
		},
	})
}

const testAccMysqlSwitchProxy = `

resource "tencentcloud_mysql_switch_proxy" "switch_proxy" {
	instance_id = "cdb-fitq5t9h"
	proxy_group_id = "proxy-h1ub486b"
  }

`
