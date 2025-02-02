/*
Provides a resource to create a as execute_scaling_policy

Example Usage

```hcl
resource "tencentcloud_as_execute_scaling_policy" "execute_scaling_policy" {
  auto_scaling_policy_id = "asp-519acdug"
  honor_cooldown = false
  trigger_source = "API"
}
```

Import

as execute_scaling_policy can be imported using the id, e.g.

```
terraform import tencentcloud_as_execute_scaling_policy.execute_scaling_policy execute_scaling_policy_id
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAsExecuteScalingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsExecuteScalingPolicyCreate,
		Read:   resourceTencentCloudAsExecuteScalingPolicyRead,
		Delete: resourceTencentCloudAsExecuteScalingPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_scaling_policy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Auto-scaling policy ID. This parameter is not available to a target tracking policy.",
			},

			"honor_cooldown": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to check if the auto scaling group is in the cooldown period. Default value: false.",
			},

			"trigger_source": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Source that triggers the scaling policy. Valid values: API and CLOUD_MONITOR. Default value: API. The value CLOUD_MONITOR is specific to the Cloud Monitor service.",
			},
		},
	}
}

func resourceTencentCloudAsExecuteScalingPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_execute_scaling_policy.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = as.NewExecuteScalingPolicyRequest()
		response   = as.NewExecuteScalingPolicyResponse()
		activityId string
	)
	if v, ok := d.GetOk("auto_scaling_policy_id"); ok {
		request.AutoScalingPolicyId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("honor_cooldown"); ok {
		request.HonorCooldown = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("trigger_source"); ok {
		request.TriggerSource = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().ExecuteScalingPolicy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate as executeScalingPolicy failed, reason:%+v", logId, err)
		return err
	}

	activityId = *response.Response.ActivityId
	d.SetId(activityId)

	return resourceTencentCloudAsExecuteScalingPolicyRead(d, meta)
}

func resourceTencentCloudAsExecuteScalingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_execute_scaling_policy.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudAsExecuteScalingPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_execute_scaling_policy.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
