package imagecheck

import (
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"
)

type Checker struct {
	Driver
}

var internal = &Checker{}

func NewChecker() Checker {
	if internal.Driver != nil {
		return *internal
	}

	driver := facades.Config().GetString("imagecheck.driver", "aliyun")
	config := cast.ToStringMapString(facades.Config().Get("imagecheck." + driver))
	switch driver {
	case "aliyun":
		internal.Driver = &Aliyun{
			AccessKeyId:     config["access_key_id"],
			AccessKeySecret: config["access_key_secret"],
		}
	case "cos":
		internal.Driver = &COS{
			AccessKey: config["access_key"],
			SecretKey: config["secret_key"],
			Bucket:    config["bucket"],
		}
	}

	return *internal
}

// Check 检查图片是否违规
func (c *Checker) Check(url string) (bool, error) {
	return c.Driver.Check(url)
}
