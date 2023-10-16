package imagecheck

import (
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"
)

type Checker struct {
	Driver
}

func NewChecker() Checker {
	driver := facades.Config().GetString("imagecheck.driver", "aliyun")
	config := cast.ToStringMapString(facades.Config().Get("imagecheck." + driver))
	var newDriver Driver

	switch driver {
	case "aliyun":
		newDriver = &Aliyun{
			AccessKeyId:     config["access_key_id"],
			AccessKeySecret: config["access_key_secret"],
		}
	case "cos":
		newDriver = &COS{
			AccessKey: config["access_key"],
			SecretKey: config["secret_key"],
			Bucket:    config["bucket"],
		}
	}

	return Checker{Driver: newDriver}
}

// Check 检查图片是否违规 true: 违规 false: 未违规
func (c *Checker) Check(url string) (bool, error) {
	return c.Driver.Check(url)
}
