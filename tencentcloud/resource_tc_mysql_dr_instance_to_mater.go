/*
Provides a resource to create a mysql dr_instance_to_mater

Example Usage

```hcl
resource "tencentcloud_mysql_dr_instance_to_mater" "dr_instance_to_mater" {
  instance_id = "cdb-c1nl9rpv"
}
```

Import

mysql dr_instance_to_mater can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_dr_instance_to_mater.dr_instance_to_mater dr_instance_to_mater_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func resourceTencentCloudMysqlDrInstanceToMater() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlDrInstanceToMaterCreate,
		Read:   resourceTencentCloudMysqlDrInstanceToMaterRead,
		Update: resourceTencentCloudMysqlDrInstanceToMaterUpdate,
		Delete: resourceTencentCloudMysqlDrInstanceToMaterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Disaster recovery instance ID in the format of cdb-c1nl9rpv. It is the same as the instance ID displayed in the TencentDB console.",
			},
		},
	}
}

func resourceTencentCloudMysqlDrInstanceToMaterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_dr_instance_to_mater.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMysqlDrInstanceToMaterUpdate(d, meta)
}

func resourceTencentCloudMysqlDrInstanceToMaterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_dr_instance_to_mater.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	drInstanceToMater, err := service.DescribeDBInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if drInstanceToMater == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlDrInstanceToMater` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if drInstanceToMater.InstanceId != nil {
		_ = d.Set("instance_id", drInstanceToMater.InstanceId)
	}

	return nil
}

func resourceTencentCloudMysqlDrInstanceToMaterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_dr_instance_to_mater.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := mysql.NewSwitchDrInstanceToMasterRequest()
	response := mysql.NewSwitchDrInstanceToMasterResponse()
	instanceId := d.Id()

	request.InstanceId = &instanceId

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().SwitchDrInstanceToMaster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mysql drInstanceToMater failed, reason:%+v", logId, err)
		return err
	}

	asyncRequestId := *response.Response.AsyncRequestId
	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s DrInstanceToMaster status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s update mysql drInstanceToMater status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s update mysql drInstanceToMater fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudMysqlDrInstanceToMaterRead(d, meta)
}

func resourceTencentCloudMysqlDrInstanceToMaterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_dr_instance_to_mater.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
